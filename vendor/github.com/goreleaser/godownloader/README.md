<p align="center">
  <img alt="GoReleaser Logo" src="https://avatars2.githubusercontent.com/u/24697112?v=3&s=200" height="140" />
  <h3 align="center">GoDownloader</h3>
  <p align="center">Download Go binaries as fast and easily as possible.</p>
  <p align="center">
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
    <a href="https://travis-ci.org/goreleaser/godownloader"><img alt="Travis" src="https://img.shields.io/travis/goreleaser/godownloader/master.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/goreleaser/godownloader"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/goreleaser/godownloader?style=flat-square"></a>
    <a href="http://godoc.org/github.com/goreleaser/goreleaser"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>
    <a href="https://github.com/goreleaser"><img alt="Powered By: GoReleaser" src="https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square"></a>
  </p>
</p>

---

This is the inverse of [goreleaser](https://github.com/goreleaser/goreleaser).  The goreleaser YAML file is read and creates a custom shell script that can download the right package and the right version for the existing machine.

If you use goreleaser already, this will create scripts suitable for "curl bash" style downloads.

* Run godownloader on your `goreleaser.yaml` file
* Add the `godownloader.sh` file to your repo
* Tell your users to use https://raw.githubusercontent.com/YOU/YOURAPP/master/godownloader.sh to install

This is also useful in CI/CD systems such as [travis-ci.org](https://travis-ci.org).

* Much faster then 'go get' (sometimes up to 100x)
* Make sure your local environment (macOS) and the CI environment (Linux) are using the exact same versions of your go binaries.

## CI/CD Example

Let's say you are using [hugo](https://gohugo.io), the static website generator, with [travis-ci](https://travis-ci.org).

Your old `.travis.yml` file might have 

```yaml
install:
  - go get github.com/gohugoio/hugo
```

This can take up to 30 seconds! 

Hugo doesn't have (yet) a `godownloader.sh` file.  So we will make our own:

```
# create a godownloader script
godownloader --repo=gohugoio/hugo > ./godownloader-hugo.sh
```

and add `godownloader-hugo.sh` to your GitHub repo.  Edit your `.travis.yml` as such

```yaml
install:
  - ./godownloader-hugo.sh v0.37.1
```

Without a version number, GitHub is queried to get the latest version number.

```yaml
install:
  - ./godownloader-hugo.sh
```

Typical download time is 0.3 seconds, or 100x improvement. 

Your new `hugo` binary is in `./bin`, so change your Makefie or scripts to use `./bin/hugo`. 

The default installation directory can be changed with the `-b` flag or the `BINDIR` environment variable.

## Notes on Functionality

* Only GitHub Releases are supported right now.
* Checksums are checked.
* Binares are installed using `tar.gz` or `zip`. 
* No OS-specific installs such as homebrew, deb, rpm.  Everything is installed locally via a `tar.gz` or `zip`.  Typically OS installs are done differently anyways (e.g. brew, apt-get, yum, etc).

## Experimental support

Some people do not use Goreleaser (why!), so there is experimental support for the following alterative distributions.

### "naked" releases on GitHub

A naked release is just the raw binary put on GitHub releases.  Limited support can be done by

```bash
./goreleaser -source raw -repo [owner/repo] -exe [name] -nametpl [tpl]
```

Where `exe` is the final binary name, and `tpl` is the same type of name template that Goreleaser uses.

An example repo is at [mvdan/sh](https://github.com/mvdan/sh/releases). Note how the repo `sh` is different than the binary `shfmt`.

### Equinox.io

[Equinox.io](https://equinox.io) is a really interesting platform.  Take a look.

There is no API, so godownloader screen scrapes to figure out the latest release.  Likewise, checksums are not verified.

```bash
./goreleaser -source equinoxio -repo [owner/repo]
```

While Equinox.io supports the concept of different release channels, only the `stable` channel is supported by godownloader.

## Yes, it's true.

It's a go program that reads a YAML file that uses a template to make a posix shell script.

## Other Resources and Inspiration

Other applications have written custom shell downloaders and installers:

### golang/dep

The [golang/dep](https://github.com/golang/dep) package manager has a nice downloader, [install.sh](https://github.com/golang/dep/blob/master/install.sh). Their trick to extract a version number from GitHub Releases is excellent:

```sh
$(echo "${LATEST_RELEASE}" | tr -s '\n' ' ' | sed 's/.*"tag_name":"//' | sed 's/".*//' )
```

This is probably based on [masterminds/glide](https://github.com/Masterminds/glide) and its installer at https://glide.sh/get

### kubernetes/helm

[kubernetes/helm](https://github.com/kubernetes/helm) is a "tool for managing Kubernetes charts. Charts are packages of pre-configured Kubernetes resources."

It has a [get script](https://github.com/kubernetes/helm/blob/master/scripts/get). Of note is that it won't re-install if the desired version is already present.

### chef

[Chef](https://www.chef.io) has the one of the most complete installers at https://omnitruck.chef.io/install.sh. In particular it has support for

* Support for solaris and aix, and some other less common platforms
* python or perl as installers if curl or wget isn't present
* http proxy support

### Caddy

[Caddy](https://caddyserver.com) is "the HTTP/2 web server with automatic HTTPS" and a NGINX replacement.  It has a clever installer at https://getcaddy.com. Of note is GPG signature verification.
