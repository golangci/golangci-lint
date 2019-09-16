// Package deb implements nfpm.Packager providing .deb bindings.
package deb

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5" // nolint:gas
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/blakesmith/ar"
	"github.com/goreleaser/nfpm"
	"github.com/goreleaser/nfpm/glob"
	"github.com/pkg/errors"
)

// nolint: gochecknoinits
func init() {
	nfpm.Register("deb", Default)
}

// nolint: gochecknoglobals
var archToDebian = map[string]string{
	"386":     "i386",
	"arm":     "armhf",
	"arm6":    "armel",
	"arm7":    "armhf",
	"mipsle":  "mipsel",
	"ppc64le": "ppc64el",
}

// Default deb packager
// nolint: gochecknoglobals
var Default = &Deb{}

// Deb is a deb packager implementation
type Deb struct{}

// Package writes a new deb package to the given writer using the given info
func (*Deb) Package(info nfpm.Info, deb io.Writer) (err error) {
	arch, ok := archToDebian[info.Arch]
	if ok {
		info.Arch = arch
	}
	dataTarGz, md5sums, instSize, err := createDataTarGz(info)
	if err != nil {
		return err
	}
	controlTarGz, err := createControl(instSize, md5sums, info)
	if err != nil {
		return err
	}
	var w = ar.NewWriter(deb)
	if err := w.WriteGlobalHeader(); err != nil {
		return errors.Wrap(err, "cannot write ar header to deb file")
	}
	if err := addArFile(w, "debian-binary", []byte("2.0\n")); err != nil {
		return errors.Wrap(err, "cannot pack debian-binary")
	}
	if err := addArFile(w, "control.tar.gz", controlTarGz); err != nil {
		return errors.Wrap(err, "cannot add control.tar.gz to deb")
	}
	if err := addArFile(w, "data.tar.gz", dataTarGz); err != nil {
		return errors.Wrap(err, "cannot add data.tar.gz to deb")
	}
	return nil
}

func addArFile(w *ar.Writer, name string, body []byte) error {
	var header = ar.Header{
		Name:    filepath.ToSlash(name),
		Size:    int64(len(body)),
		Mode:    0644,
		ModTime: time.Now(),
	}
	if err := w.WriteHeader(&header); err != nil {
		return errors.Wrap(err, "cannot write file header")
	}
	_, err := w.Write(body)
	return err
}

func createDataTarGz(info nfpm.Info) (dataTarGz, md5sums []byte, instSize int64, err error) {
	var buf bytes.Buffer
	var compress = gzip.NewWriter(&buf)
	var out = tar.NewWriter(compress)

	// the writers are properly closed later, this is just in case that we have
	// an error in another part of the code.
	defer out.Close()      // nolint: errcheck
	defer compress.Close() // nolint: errcheck

	var created = map[string]bool{}
	if err = createEmptyFoldersInsideTarGz(info, out, created); err != nil {
		return nil, nil, 0, err
	}

	md5buf, instSize, err := createFilesInsideTarGz(info, out, created)
	if err != nil {
		return nil, nil, 0, err
	}

	if err := out.Close(); err != nil {
		return nil, nil, 0, errors.Wrap(err, "closing data.tar.gz")
	}
	if err := compress.Close(); err != nil {
		return nil, nil, 0, errors.Wrap(err, "closing data.tar.gz")
	}

	return buf.Bytes(), md5buf.Bytes(), instSize, nil
}

func createFilesInsideTarGz(info nfpm.Info, out *tar.Writer, created map[string]bool) (bytes.Buffer, int64, error) {
	var md5buf bytes.Buffer
	var instSize int64
	for _, files := range []map[string]string{
		info.Files,
		info.ConfigFiles,
	} {
		for srcglob, dstroot := range files {
			globbed, err := glob.Glob(srcglob, dstroot)
			if err != nil {
				return md5buf, 0, err
			}
			for src, dst := range globbed {
				if err := createTree(out, dst, created); err != nil {
					return md5buf, 0, err
				}
				size, err := copyToTarAndDigest(out, &md5buf, src, dst)
				if err != nil {
					return md5buf, 0, err
				}
				instSize += size
			}
		}
	}
	return md5buf, instSize, nil
}

func createEmptyFoldersInsideTarGz(info nfpm.Info, out *tar.Writer, created map[string]bool) error {
	for _, folder := range info.EmptyFolders {
		// this .nope is actually not created, because createTree ignore the
		// last part of the path, assuming it is a file.
		// TODO: should probably refactor this
		if err := createTree(out, filepath.Join(folder, ".nope"), created); err != nil {
			return err
		}
	}
	return nil
}

func copyToTarAndDigest(tarw *tar.Writer, md5w io.Writer, src, dst string) (int64, error) {
	file, err := os.OpenFile(src, os.O_RDONLY, 0600) //nolint:gosec
	if err != nil {
		return 0, errors.Wrap(err, "could not add file to the archive")
	}
	// don't care if it errs while closing...
	defer file.Close() // nolint: errcheck
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	if info.IsDir() {
		// TODO: this should probably return an error
		return 0, nil
	}
	var header = tar.Header{
		Name:    filepath.ToSlash(dst[1:]),
		Size:    info.Size(),
		Mode:    int64(info.Mode()),
		ModTime: time.Now(),
		Format:  tar.FormatGNU,
	}
	if err := tarw.WriteHeader(&header); err != nil {
		return 0, errors.Wrapf(err, "cannot write header of %s to data.tar.gz", src)
	}
	var digest = md5.New() // nolint:gas
	if _, err := io.Copy(tarw, io.TeeReader(file, digest)); err != nil {
		return 0, errors.Wrap(err, "failed to copy")
	}
	if _, err := fmt.Fprintf(md5w, "%x  %s\n", digest.Sum(nil), header.Name); err != nil {
		return 0, errors.Wrap(err, "failed to write md5")
	}
	return info.Size(), nil
}

func createControl(instSize int64, md5sums []byte, info nfpm.Info) (controlTarGz []byte, err error) {
	var buf bytes.Buffer
	var compress = gzip.NewWriter(&buf)
	var out = tar.NewWriter(compress)
	// the writers are properly closed later, this is just in case that we have
	// an error in another part of the code.
	defer out.Close()      // nolint: errcheck
	defer compress.Close() // nolint: errcheck

	var body bytes.Buffer
	if err := writeControl(&body, controlData{
		Info:          info,
		InstalledSize: instSize / 1024,
	}); err != nil {
		return nil, err
	}

	for name, content := range map[string][]byte{
		"control":   body.Bytes(),
		"md5sums":   md5sums,
		"conffiles": conffiles(info),
	} {
		if err := newFileInsideTarGz(out, name, content); err != nil {
			return nil, err
		}
	}

	for script, dest := range map[string]string{
		info.Scripts.PreInstall:             "preinst",
		info.Scripts.PostInstall:            "postinst",
		info.Scripts.PreRemove:              "prerm",
		info.Scripts.PostRemove:             "postrm",
		info.Overridables.Deb.Scripts.Rules: "rules",
	} {
		if script != "" {
			if err := newScriptInsideTarGz(out, script, dest); err != nil {
				return nil, err
			}
		}
	}

	if err := out.Close(); err != nil {
		return nil, errors.Wrap(err, "closing control.tar.gz")
	}
	if err := compress.Close(); err != nil {
		return nil, errors.Wrap(err, "closing control.tar.gz")
	}
	return buf.Bytes(), nil
}

func newItemInsideTarGz(out *tar.Writer, content []byte, header tar.Header) error {
	if err := out.WriteHeader(&header); err != nil {
		return errors.Wrapf(err, "cannot write header of %s file to control.tar.gz", header.Name)
	}
	if _, err := out.Write(content); err != nil {
		return errors.Wrapf(err, "cannot write %s file to control.tar.gz", header.Name)
	}
	return nil
}

func newFileInsideTarGz(out *tar.Writer, name string, content []byte) error {
	return newItemInsideTarGz(out, content, tar.Header{
		Name:     filepath.ToSlash(name),
		Size:     int64(len(content)),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
		Format:   tar.FormatGNU,
	})
}

func newScriptInsideTarGz(out *tar.Writer, path string, dest string) error {
	file, err := os.Open(path) //nolint:gosec
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	return newItemInsideTarGz(out, content, tar.Header{
		Name:     filepath.ToSlash(dest),
		Size:     int64(len(content)),
		Mode:     0755,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
		Format:   tar.FormatGNU,
	})
}

// this is needed because the data.tar.gz file should have the empty folders
// as well, so we walk through the dst and create all subfolders.
func createTree(tarw *tar.Writer, dst string, created map[string]bool) error {
	for _, path := range pathsToCreate(dst) {
		if created[path] {
			// skipping dir that was previously created inside the archive
			// (eg: usr/)
			continue
		}
		if err := tarw.WriteHeader(&tar.Header{
			Name:     filepath.ToSlash(path + "/"),
			Mode:     0755,
			Typeflag: tar.TypeDir,
			Format:   tar.FormatGNU,
			ModTime:  time.Now(),
		}); err != nil {
			return errors.Wrap(err, "failed to create folder")
		}
		created[path] = true
	}
	return nil
}

func pathsToCreate(dst string) []string {
	var paths = []string{}
	var base = dst[1:]
	for {
		base = filepath.Dir(base)
		if base == "." {
			break
		}
		paths = append(paths, base)
	}
	// we don't really need to create those things in order apparently, but,
	// it looks really weird if we don't.
	var result = []string{}
	for i := len(paths) - 1; i >= 0; i-- {
		result = append(result, paths[i])
	}
	return result
}

func conffiles(info nfpm.Info) []byte {
	// nolint: prealloc
	var confs []string
	for _, dst := range info.ConfigFiles {
		confs = append(confs, dst)
	}
	return []byte(strings.Join(confs, "\n") + "\n")
}

const controlTemplate = `
{{- /* Mandatory fields */ -}}
Package: {{.Info.Name}}
{{- if .Info.Epoch}}
Version: {{ .Info.Epoch }}:{{.Info.Version}}
{{- else }}
Version: {{.Info.Version}}
{{- end }}
Section: {{.Info.Section}}
Priority: {{.Info.Priority}}
Architecture: {{.Info.Arch}}
{{- /* Optional fields */ -}}
{{- if .Info.Maintainer}}
Maintainer: {{.Info.Maintainer}}
{{- end }}
{{- if .Info.Vendor}}
Vendor: {{.Info.Vendor}}
{{- end }}
Installed-Size: {{.InstalledSize}}
{{- with .Info.Replaces}}
Replaces: {{join .}}
{{- end }}
{{- with .Info.Provides}}
Provides: {{join .}}
{{- end }}
{{- with .Info.Depends}}
Depends: {{join .}}
{{- end }}
{{- with .Info.Recommends}}
Recommends: {{join .}}
{{- end }}
{{- with .Info.Suggests}}
Suggests: {{join .}}
{{- end }}
{{- with .Info.Conflicts}}
Conflicts: {{join .}}
{{- end }}
{{- if .Info.Homepage}}
Homepage: {{.Info.Homepage}}
{{- end }}
{{- /* Mandatory fields */}}
Description: {{.Info.Description}}
`

type controlData struct {
	Info          nfpm.Info
	InstalledSize int64
}

func writeControl(w io.Writer, data controlData) error {
	var tmpl = template.New("control")
	tmpl.Funcs(template.FuncMap{
		"join": func(strs []string) string {
			return strings.Trim(strings.Join(strs, ", "), " ")
		},
	})
	return template.Must(tmpl.Parse(controlTemplate)).Execute(w, data)
}
