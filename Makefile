# MakeFile

.PHONY: clean

export GOPATH=$(PWD)/../..


THE_GO_FILES = \
	filehandler.go \
	filePacker.go \
	main.go \
	utils.go \
	response.go

all: FilesOnWebRun
	@echo "Make FilesOnWebRun Succeed"

FilesOnWebRun: $(THE_GO_FILES)
	go build FilesOnWebRun

run: FilesOnWebRun
	@./FilesOnWebRun || true

clean: FilesOnWebRun
	@rm FilesOnWebRun
