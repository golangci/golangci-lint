.PHONY: default clean serve build

default: clean serve

clean:
	rm -rf public/

serve: clean
	hugo server --disableFastRender --buildDrafts --enableGitInfo

# IMPORTANT used by the CI to deploy documentation
build: clean
	HUGO_ENVIRONMENT=production \
	HUGO_ENV=production \
	hugo \
		--gc --minify \
		--baseURL "https://golangci-lint.run/"
