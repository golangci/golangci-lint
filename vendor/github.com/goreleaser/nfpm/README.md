<p align="center">
  <img alt="GoReleaser Logo" src="https://avatars2.githubusercontent.com/u/24697112?v=3&s=200" height="140" />
  <h3 align="center">NFPM</h3>
  <p align="center">NFPM is Not FPM - a simple deb and rpm packager written in Go.</p>
  <p align="center">
    <a href="https://github.com/goreleaser/nfpm/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/goreleaser/nfpm.svg?style=flat-square"></a>
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
    <a href="https://travis-ci.org/goreleaser/nfpm"><img alt="Travis" src="https://img.shields.io/travis/goreleaser/nfpm/master.svg?style=flat-square"></a>
    <a href="https://codecov.io/gh/goreleaser/nfpm"><img alt="Codecov branch" src="https://img.shields.io/codecov/c/github/goreleaser/nfpm/master.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/goreleaser/nfpm"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/goreleaser/nfpm?style=flat-square"></a>
    <a href="http://godoc.org/github.com/goreleaser/nfpm"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>
    <a href="https://saythanks.io/to/caarlos0"><img alt="SayThanks.io" src="https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg?style=flat-square"></a>
    <a href="https://github.com/goreleaser"><img alt="Powered By: GoReleaser" src="https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square"></a>
  </p>
</p>

## Why

While [fpm][] is great, for me it is a bummer that it depends on Ruby, tar
and probably other software.

I wanted something that could be used as a binary and/or as a lib on go-software,
so I hacked this together and it works!

## Goals

* [x] be simple to use
* [x] provide packaging for the most common linux packaging systems (at very least deb and rpm)
* [x] be distributed as a single binary
* [x] reproducible results
  * [x] depend on the fewer external things as possible (namely `rpmbuild`)
  * [x] generate packages based on yaml files (maybe also json and toml?)
* [x] be possible to use it as a lib in other go projects (namely [goreleaser][] itself)
* [ ] support complex packages and power users

## Usage

The first steps are to run `nfpm init` to initialize a config file and edit
the generated file according to your needs:

![nfpm init](https://user-images.githubusercontent.com/245435/36346101-f81cdcec-141e-11e8-8afc-a5eb93b7d510.png)

The next step is to run `nfpm pkg --target mypkg.deb`.
NFPM will guess which packager to use based on the target file extension.

![nfpm pkg](https://user-images.githubusercontent.com/245435/36346100-eaaf24c0-141e-11e8-8345-100f4d3ed02d.png)

And that's it!

## Usage as a docker image

You can run it with docker as well:

```sh
docker run --rm \
  -v $PWD:/tmp/pkg \
  goreleaser/nfpm pkg --config /tmp/pkg/foo.yml --target /tmp/pkg/foo.rpm
```

That's it!

## Usage as lib

You can look at the code of nfpm itself to see how to use it as a library, or, take
a look at the [nfpm pipe on GoReleaser](https://github.com/goreleaser/goreleaser/blob/master/internal/pipe/nfpm/nfpm.go).

> **Attention**: GoReleaser `deb` packager only compiles with go1.10+.

## Status

* both deb and rpm packaging are working but there are some missing features.

## Special thanks

Thanks to the [fpm][] authors for fpm, which inspires nfpm a lot.

## Donate

Donations are very much appreciated! You can donate/sponsor on the main
[goreleaser opencollective](https://opencollective.com/goreleaser)! It's
easy and will surely help the developers at least buy some ☕️ or 🍺!

## Stargazers over time

[![goreleaser/nfpm stargazers over time](https://starcharts.herokuapp.com/goreleaser/nfpm.svg)](https://starcharts.herokuapp.com/goreleaser/nfpm)

---

Would you like to fix something in the documentation? Feel free to open an [issue](https://github.com/goreleaser/nfpm/issues).

[goreleaser]: https://goreleaser.com/#linux_packages.nfpm
[fpm]: https://github.com/jordansissel/fpm
