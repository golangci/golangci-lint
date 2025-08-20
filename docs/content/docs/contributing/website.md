---
title: This Website
weight: 7
aliases:
  - /contributing/website/
---

## Technology

We use [Hugo](https://gohugo.io/) for static site generation because sites built with it are very fast.

## Source Code

The website lives in `docs/` directory of [golangci-lint repository](https://github.com/golangci/golangci-lint).

## Theme

The site is based on [hextra](https://github.com/imfing/hextra) theme.

## Templating

We use [shortcodes](https://gohugo.io/templates/types/#shortcode) and [partials](https://gohugo.io/templates/types/#partial) based on files from `./docs/.tmp/` and `./docs/data/`. 

- The files in `./docs/.tmp/` are used to be embedded with the shortcode `{{%/* golangci/embed file="filename.ext" */%}}`.
- The files in `./docs/data/` are used as [data sources](https://gohugo.io/content-management/data-sources/). 

These files are created by running:

- `make website_expand_templates` in the root of the repository.  
- `make website_dump_info` in the root of the repository. (only during a release)

### Some Notes

[shortcodes](https://gohugo.io/templates/types/#shortcode):
- cannot be used inside another shortcode
- can only be used inside a page
- can contain Markdown or HTML, but the tag is different: `{{%/* shortcode */%}}` vs `{{</* shortcode */>}}`

[partials](https://gohugo.io/templates/types/#partial):
- are reusable HTML blocks or "functions"
- cannot be used inside a page
- can be used inside another partial
- can be used inside a shortcode
- can be used inside a layout

## Hosting

We use GitHub Pages for static website hosting and CD.

GitHub deploys the website to production after merging anything to a `main` branch.

## Local Testing

Install Hugo Extended (v0.148.1 or newer).

Run:

```bash
# (in the root of the repository)
make docs_serve
```

or

```bash
# (in the root of the repository)
make website_expand_templates

cd docs/

# (inside the docs/ folder)
make serve
```

And navigate to `http://localhost:1313` after a successful build.

There is no need to restart the Hugo server for almost all changes: it supports hot reload.  
Also, there is no need to refresh a webpage: hot reload updates changed content on the open page.

## Website Build

To do this, run:

```bash
# (in the root of the repository)
make docs_build
```
or

```bash
# (in the root of the repository)
make website_copy_jsonschema website_expand_templates

cd docs/

# (inside the docs/ folder)
make build
```
