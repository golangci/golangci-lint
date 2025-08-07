.PHONY: default clean serve

default: clean serve

clean:
	rm -rf public/

serve: clean
	hugo server --disableFastRender --buildDrafts --enableGitInfo
