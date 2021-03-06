WHAT ?= dash

BRANCH      ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILDTIME   ?= $(shell date '+%Y-%m-%d@%H:%M:%S')
BUILDUSER   ?= $(shell id -un)
REPO        ?= github.com/ricoberger/dash
REVISION    ?= $(shell git rev-parse HEAD)
VERSION     ?= $(shell git describe --tags)

.PHONY: build build-darwin-amd64 build-linux-amd64 build-windows-amd64 clean release release-major release-minor release-patch

build:
	for target in $(WHAT); do \
		go build -ldflags "-X ${REPO}/pkg/version.Version=${VERSION} \
			-X ${REPO}/pkg/version.Revision=${REVISION} \
			-X ${REPO}/pkg/version.Branch=${BRANCH} \
			-X ${REPO}/pkg/version.BuildUser=${BUILDUSER} \
			-X ${REPO}/pkg/version.BuildDate=${BUILDTIME}" \
			-o ./bin/$$target ./cmd/$$target; \
	done

build-darwin-amd64:
	for target in $(WHAT); do \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -a -installsuffix cgo -ldflags "-X ${REPO}/pkg/version.Version=${VERSION} \
			-X ${REPO}/pkg/version.Revision=${REVISION} \
			-X ${REPO}/pkg/version.Branch=${BRANCH} \
			-X ${REPO}/pkg/version.BuildUser=${BUILDUSER} \
			-X ${REPO}/pkg/version.BuildDate=${BUILDTIME}" \
			-o ./bin/$$target-darwin-amd64 ./cmd/$$target; \
	done

build-linux-amd64:
	for target in $(WHAT); do \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags "-X ${REPO}/pkg/version.Version=${VERSION} \
			-X ${REPO}/pkg/version.Revision=${REVISION} \
			-X ${REPO}/pkg/version.Branch=${BRANCH} \
			-X ${REPO}/pkg/version.BuildUser=${BUILDUSER} \
			-X ${REPO}/pkg/version.BuildDate=${BUILDTIME}" \
			-o ./bin/$$target-linux-amd64 ./cmd/$$target; \
	done

build-windows-amd64:
	for target in $(WHAT); do \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -a -installsuffix cgo -ldflags "-X ${REPO}/pkg/version.Version=${VERSION} \
			-X ${REPO}/pkg/version.Revision=${REVISION} \
			-X ${REPO}/pkg/version.Branch=${BRANCH} \
			-X ${REPO}/pkg/version.BuildUser=${BUILDUSER} \
			-X ${REPO}/pkg/version.BuildDate=${BUILDTIME}" \
			-o ./bin/$$target-windows-amd64.exe ./cmd/$$target; \
	done

clean:
	rm -rf ./bin

release: clean build-darwin-amd64 build-linux-amd64 build-windows-amd64

release-major:
	$(eval MAJORVERSION=$(shell git describe --tags --abbrev=0 | sed s/v// | awk -F. '{print $$1+1".0.0"}'))
	git checkout master
	git pull
	git tag -a $(MAJORVERSION) -m 'Release $(MAJORVERSION)'
	git push origin --tags

release-minor:
	$(eval MINORVERSION=$(shell git describe --tags --abbrev=0 | sed s/v// | awk -F. '{print $$1"."$$2+1".0"}'))
	git checkout master
	git pull
	git tag -a $(MINORVERSION) -m 'Release $(MINORVERSION)'
	git push origin --tags

release-patch:
	$(eval PATCHVERSION=$(shell git describe --tags --abbrev=0 | sed s/v// | awk -F. '{print $$1"."$$2"."$$3+1}'))
	git checkout master
	git pull
	git tag -a $(PATCHVERSION) -m 'Release $(PATCHVERSION)'
	git push origin --tags
