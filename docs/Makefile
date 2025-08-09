.PHONY: default clean reset serve

default: clean serve

clean:
	rm -rf public/

reset:
	git checkout HEAD content/

serve: clean
	hugo server --disableFastRender --buildDrafts --enableGitInfo --watch --source .
