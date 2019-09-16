// Package http implements functionality common to HTTP uploading pipelines.
package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"html/template"
	"io"
	h "net/http"
	"os"
	"runtime"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/internal/semerrgroup"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
)

const (
	// ModeBinary uploads only compiled binaries
	ModeBinary = "binary"
	// ModeArchive uploads release archives
	ModeArchive = "archive"
)

type asset struct {
	ReadCloser io.ReadCloser
	Size       int64
}

type assetOpenFunc func(string, *artifact.Artifact) (*asset, error)

// nolint: gochecknoglobals
var assetOpen assetOpenFunc

// TODO: fix this.
// nolint: gochecknoinits
func init() {
	assetOpenReset()
}

func assetOpenReset() {
	assetOpen = assetOpenDefault
}

func assetOpenDefault(kind string, a *artifact.Artifact) (*asset, error) {
	f, err := os.Open(a.Path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, errors.Errorf("%s: upload failed: the asset to upload can't be a directory", kind)
	}
	return &asset{
		ReadCloser: f,
		Size:       s.Size(),
	}, nil
}

// Defaults sets default configuration options on Put structs
func Defaults(puts []config.Put) error {
	for i := range puts {
		defaults(&puts[i])
	}
	return nil
}

func defaults(put *config.Put) {
	if put.Mode == "" {
		put.Mode = ModeArchive
	}
}

// CheckConfig validates a Put configuration returning a descriptive error when appropriate
func CheckConfig(ctx *context.Context, put *config.Put, kind string) error {
	if put.Target == "" {
		return misconfigured(kind, put, "missing target")
	}

	if put.Name == "" {
		return misconfigured(kind, put, "missing name")
	}

	if put.Mode != ModeArchive && put.Mode != ModeBinary {
		return misconfigured(kind, put, "mode must be 'binary' or 'archive'")
	}

	envName := fmt.Sprintf("%s_%s_SECRET", strings.ToUpper(kind), strings.ToUpper(put.Name))
	if _, ok := ctx.Env[envName]; !ok {
		return misconfigured(kind, put, fmt.Sprintf("missing %s environment variable", envName))
	}

	if put.TrustedCerts != "" && !x509.NewCertPool().AppendCertsFromPEM([]byte(put.TrustedCerts)) {
		return misconfigured(kind, put, "no certificate could be added from the specified trusted_certificates configuration")
	}

	return nil

}

func misconfigured(kind string, upload *config.Put, reason string) error {
	return pipe.Skip(fmt.Sprintf("%s section '%s' is not configured properly (%s)", kind, upload.Name, reason))
}

// ResponseChecker is a function capable of validating an http server response.
// It must return and error when the response must be considered a failure.
type ResponseChecker func(*h.Response) error

// Upload does the actual uploading work
func Upload(ctx *context.Context, puts []config.Put, kind string, check ResponseChecker) error {
	if ctx.SkipPublish {
		return pipe.ErrSkipPublishEnabled
	}

	// Handle every configured put
	for _, put := range puts {
		put := put
		filters := []artifact.Filter{}
		if put.Checksum {
			filters = append(filters, artifact.ByType(artifact.Checksum))
		}
		if put.Signature {
			filters = append(filters, artifact.ByType(artifact.Signature))
		}
		// We support two different modes
		//	- "archive": Upload all artifacts
		//	- "binary": Upload only the raw binaries
		switch v := strings.ToLower(put.Mode); v {
		case ModeArchive:
			filters = append(filters,
				artifact.ByType(artifact.UploadableArchive),
				artifact.ByType(artifact.LinuxPackage),
			)
		case ModeBinary:
			filters = append(filters, artifact.ByType(artifact.UploadableBinary))
		default:
			err := fmt.Errorf("%s: mode \"%s\" not supported", kind, v)
			log.WithFields(log.Fields{
				kind:   put.Name,
				"mode": v,
			}).Error(err.Error())
			return err
		}

		var filter = artifact.Or(filters...)
		if len(put.IDs) > 0 {
			filter = artifact.And(filter, artifact.ByIDs(put.IDs...))
		}
		if err := uploadWithFilter(ctx, &put, filter, kind, check); err != nil {
			return err
		}
	}

	return nil
}

func uploadWithFilter(ctx *context.Context, put *config.Put, filter artifact.Filter, kind string, check ResponseChecker) error {
	var artifacts = ctx.Artifacts.Filter(filter).List()
	log.Debugf("will upload %d artifacts", len(artifacts))
	var g = semerrgroup.New(ctx.Parallelism)
	for _, artifact := range artifacts {
		artifact := artifact
		g.Go(func() error {
			return uploadAsset(ctx, put, artifact, kind, check)
		})
	}
	return g.Wait()
}

// uploadAsset uploads file to target and logs all actions
func uploadAsset(ctx *context.Context, put *config.Put, artifact *artifact.Artifact, kind string, check ResponseChecker) error {
	envBase := fmt.Sprintf("%s_%s_", strings.ToUpper(kind), strings.ToUpper(put.Name))
	username := put.Username
	if username == "" {
		// username not configured: using env
		username = ctx.Env[envBase+"USERNAME"]
	}
	secret := ctx.Env[envBase+"SECRET"]

	// Generate the target url
	targetURL, err := resolveTargetTemplate(ctx, put, artifact)
	if err != nil {
		msg := fmt.Sprintf("%s: error while building the target url", kind)
		log.WithField("instance", put.Name).WithError(err).Error(msg)
		return errors.Wrap(err, msg)
	}

	// Handle the artifact
	asset, err := assetOpen(kind, artifact)
	if err != nil {
		return err
	}
	defer asset.ReadCloser.Close() // nolint: errcheck

	// The target url needs to contain the artifact name
	if !strings.HasSuffix(targetURL, "/") {
		targetURL += "/"
	}
	targetURL += artifact.Name

	var headers = map[string]string{}
	if put.ChecksumHeader != "" {
		sum, err := artifact.Checksum("sha256")
		if err != nil {
			return err
		}
		headers[put.ChecksumHeader] = sum
	}

	res, err := uploadAssetToServer(ctx, put, targetURL, username, secret, headers, asset, check)
	if err != nil {
		msg := fmt.Sprintf("%s: upload failed", kind)
		log.WithError(err).WithFields(log.Fields{
			"instance": put.Name,
			"username": username,
		}).Error(msg)
		return errors.Wrap(err, msg)
	}
	if err := res.Body.Close(); err != nil {
		log.WithError(err).Warn("failed to close response body")
	}

	log.WithFields(log.Fields{
		"instance": put.Name,
		"mode":     put.Mode,
	}).Info("uploaded successful")

	return nil
}

// uploadAssetToServer uploads the asset file to target
func uploadAssetToServer(ctx *context.Context, put *config.Put, target, username, secret string, headers map[string]string, a *asset, check ResponseChecker) (*h.Response, error) {
	req, err := newUploadRequest(target, username, secret, headers, a)
	if err != nil {
		return nil, err
	}

	return executeHTTPRequest(ctx, put, req, check)
}

// newUploadRequest creates a new h.Request for uploading
func newUploadRequest(target, username, secret string, headers map[string]string, a *asset) (*h.Request, error) {
	req, err := h.NewRequest(h.MethodPut, target, a.ReadCloser)
	if err != nil {
		return nil, err
	}
	req.ContentLength = a.Size
	req.SetBasicAuth(username, secret)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	return req, err
}

func getHTTPClient(put *config.Put) (*h.Client, error) {
	if put.TrustedCerts == "" {
		return h.DefaultClient, nil
	}
	pool, err := x509.SystemCertPool()
	if err != nil {
		if runtime.GOOS == "windows" {
			// on windows ignore errors until golang issues #16736 & #18609 get fixed
			pool = x509.NewCertPool()
		} else {
			return nil, err
		}
	}
	pool.AppendCertsFromPEM([]byte(put.TrustedCerts)) // already validated certs checked by CheckConfig
	return &h.Client{
		Transport: &h.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: pool,
			},
		},
	}, nil
}

// executeHTTPRequest processes the http call with respect of context ctx
func executeHTTPRequest(ctx *context.Context, put *config.Put, req *h.Request, check ResponseChecker) (*h.Response, error) {
	client, err := getHTTPClient(put)
	if err != nil {
		return nil, err
	}
	log.Debugf("executing request: %s %s (headers: %v)", req.Method, req.URL, req.Header)
	resp, err := client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	defer resp.Body.Close() // nolint: errcheck

	err = check(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}

	return resp, err
}

// targetData is used as a template struct for
// Artifactory.Target
type targetData struct {
	Version     string
	Tag         string
	ProjectName string

	// Only supported in mode binary
	Os   string
	Arch string
	Arm  string
}

// resolveTargetTemplate returns the resolved target template with replaced variables
// Those variables can be replaced by the given context, goos, goarch, goarm and more
// TODO: replace this with our internal template pkg
func resolveTargetTemplate(ctx *context.Context, put *config.Put, artifact *artifact.Artifact) (string, error) {
	data := targetData{
		Version:     ctx.Version,
		Tag:         ctx.Git.CurrentTag,
		ProjectName: ctx.Config.ProjectName,
	}

	if put.Mode == ModeBinary {
		// TODO: multiple archives here
		data.Os = replace(ctx.Config.Archive.Replacements, artifact.Goos)
		data.Arch = replace(ctx.Config.Archive.Replacements, artifact.Goarch)
		data.Arm = replace(ctx.Config.Archive.Replacements, artifact.Goarm)
	}

	var out bytes.Buffer
	t, err := template.New(ctx.Config.ProjectName).Parse(put.Target)
	if err != nil {
		return "", err
	}
	err = t.Execute(&out, data)
	return out.String(), err
}

func replace(replacements map[string]string, original string) string {
	result := replacements[original]
	if result == "" {
		return original
	}
	return result
}
