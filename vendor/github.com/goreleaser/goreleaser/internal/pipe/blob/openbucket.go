package blob

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/semerrgroup"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/pkg/errors"
	"gocloud.dev/blob"
	"gocloud.dev/secrets"

	// Import the blob packages we want to be able to open.
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"

	// import the secrets packages we want to be able to open:
	_ "gocloud.dev/secrets/awskms"
	_ "gocloud.dev/secrets/azurekeyvault"
	_ "gocloud.dev/secrets/gcpkms"
)

// OpenBucket is the interface that wraps the BucketConnect and UploadBucket method
type OpenBucket interface {
	Connect(ctx *context.Context, bucketURL string) (*blob.Bucket, error)
	Upload(ctx *context.Context, conf config.Blob, folder string) error
}

// Bucket is object which holds connection for Go Bucker Provider
type Bucket struct {
	BucketConn *blob.Bucket
}

// returns openbucket connection for list of providers
func newOpenBucket() OpenBucket {
	return Bucket{}
}

// Connect makes connection with provider
func (b Bucket) Connect(ctx *context.Context, bucketURL string) (*blob.Bucket, error) {
	conn, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Upload takes connection initilized from newOpenBucket to upload goreleaser artifacts
// Takes goreleaser context(which includes artificats) and bucketURL for upload destination (gs://gorelease-bucket)
func (b Bucket) Upload(ctx *context.Context, conf config.Blob, folder string) error {
	var bucketURL = fmt.Sprintf("%s://%s", conf.Provider, conf.Bucket)

	// Get the openbucket connection for specific provider
	conn, err := b.Connect(ctx, bucketURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	var filter = artifact.Or(
		artifact.ByType(artifact.UploadableArchive),
		artifact.ByType(artifact.UploadableBinary),
		artifact.ByType(artifact.Checksum),
		artifact.ByType(artifact.Signature),
		artifact.ByType(artifact.LinuxPackage),
	)
	if len(conf.IDs) > 0 {
		filter = artifact.And(filter, artifact.ByIDs(conf.IDs...))
	}

	var g = semerrgroup.New(ctx.Parallelism)
	for _, artifact := range ctx.Artifacts.Filter(filter).List() {
		artifact := artifact
		g.Go(func() error {
			log.WithFields(log.Fields{
				"provider": bucketURL,
				"folder":   folder,
				"artifact": artifact.Name,
			}).Info("uploading")

			w, err := conn.NewWriter(ctx, filepath.Join(folder, artifact.Name), nil)
			if err != nil {
				return errors.Wrap(err, "failed to obtain writer")
			}
			data, err := getData(ctx, conf, artifact.Path)
			if err != nil {
				return err
			}
			_, err = w.Write(data)
			if err != nil {
				switch {
				case errorContains(err, "NoSuchBucket", "ContainerNotFound", "notFound"):
					return errors.Wrapf(err, "provided bucket does not exist: %s", bucketURL)
				case errorContains(err, "NoCredentialProviders"):
					return errors.Wrapf(err, "check credentials and access to bucket: %s", bucketURL)
				default:
					return errors.Wrapf(err, "failed to write to bucket")
				}
			}
			if err = w.Close(); err != nil {
				switch {
				case errorContains(err, "InvalidAccessKeyId"):
					return errors.Wrap(err, "aws access key id you provided does not exist in our records")
				case errorContains(err, "AuthenticationFailed"):
					return errors.Wrap(err, "azure storage key you provided is not valid")
				case errorContains(err, "invalid_grant"):
					return errors.Wrap(err, "google app credentials you provided is not valid")
				case errorContains(err, "no such host"):
					return errors.Wrap(err, "azure storage account you provided is not valid")
				case errorContains(err, "NoSuchBucket", "ContainerNotFound", "notFound"):
					return errors.Wrapf(err, "provided bucket does not exist: %s", bucketURL)
				case errorContains(err, "NoCredentialProviders"):
					return errors.Wrapf(err, "check credentials and access to bucket %s", bucketURL)
				case errorContains(err, "ServiceCode=ResourceNotFound"):
					return errors.Wrapf(err, "missing azure storage key for provided bucket %s", bucketURL)
				default:
					return errors.Wrap(err, "failed to close Bucket writer")
				}
			}
			return err
		})
	}
	return g.Wait()
}

func getData(ctx *context.Context, conf config.Blob, path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return data, errors.Wrapf(err, "failed to open file %s", path)
	}
	if conf.KMSKey == "" {
		return data, nil
	}
	keeper, err := secrets.OpenKeeper(ctx, conf.KMSKey)
	if err != nil {
		return data, errors.Wrapf(err, "failed to open kms %s", conf.KMSKey)
	}
	defer keeper.Close()
	data, err = keeper.Encrypt(ctx, data)
	if err != nil {
		return data, errors.Wrap(err, "failed to encrypt with kms")
	}
	return data, err
}
