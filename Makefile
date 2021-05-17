install:
	go install doxctl

add_commit_push:
	git add .
	git commit -m "Makefile commit"
	git push





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
