export GITHUB_TOKEN = ${GO_RELEASER_GITHUB_TOKEN}

.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null \
		| awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' \
		| sort \
		| egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

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
	make tag ; git add . ; git commit -m "Makefile commit" ; git push #; make tag
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

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
