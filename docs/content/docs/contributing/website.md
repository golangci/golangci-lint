---
title: This Website
weight: 7
---

## Technology

We use [Hugo](https://gohugo.io/) for static site generation because sites built with it are very fast.

## Source Code

The website lives in `docs/` directory of [golangci-lint repository](https://github.com/golangci/golangci-lint).

## Theme

The site is based on [hextra](https://github.com/imfing/hextra) theme.

## Templating

We use templates like `{.SomeField}` inside our `md` files.

There templates are expanded by running `make website_expand_templates` in the root of the repository.  
It runs script `scripts/website/expand_templates/` that rewrites `md` files with replaced templates.

## Hosting

We use GitHub Pages as static website hosting and CD.

GitHub deploys the website to production after merging anything to a `main` branch.

## Local Testing

Install Hugo (v0.148.1 or newer).

Run:

```bash
 hugo server --buildDrafts --disableFastRender
```

And navigate to `http://localhost:1313` after successful build.

There is no need to restart Hugo server almost for all changes: it supports hot reload.  
Also, there is no need to refresh a webpage: hot reload updates changed content on the open page.

## Website Build

To do it run:

```bash
make website_expand_templates
```
