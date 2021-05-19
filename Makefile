export GITHUB_TOKEN = ${GO_RELEASER_GITHUB_TOKEN}

install:
	go install doxctl

add_commit_push:
	git add .
	git commit -m "Makefile commit"
	git push


dryrun:
	goreleaser --snapshot --skip-publish --rm-dist --debug
build:
	goreleaser build --rm-dist --debug
tag:
	scripts/version-up.sh --patch --apply
release:
	goreleaser release --rm-dist
commit:
	git add . ; git commit -m "Makefile commit" ; git push ; make tag


### Useful for debugging ###
#goreleaser release --skip-validate --rm-dist --debug #--skip-publish



##########################################################3
#   $ mkdir doxctl && cd doxctl
#   $ go mod init doxctl
#   $ go get -u github.com/spf13/cobra/cobra
#   $ cobra init --pkg-name doxctl
#   $ go get github.com/mitchellh/go-homedir@v1.1.0
#   $ go install doxctl
#
#
#
#   REFS:
#     - https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177
