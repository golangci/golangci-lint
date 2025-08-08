.PHONY: default clean serve build

default: clean serve

clean:
	rm -rf public/

serve: clean
	hugo server --disableFastRender --buildDrafts --enableGitInfo

build: clean
	HUGO_ENVIRONMENT=production \
	HUGO_ENV=production \
	hugo \
		--baseURL "https://golangci-lint.run/"
