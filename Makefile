BRANCH ?= "master"
REPONAME ?= "terraform-linter"
VERSION ?= $(shell cat ./VERSION)
PACKAGES ?= "./..."
CMD_PATH ?= "github.com/vidsy/terraform-linter/cmd/terraform/linter"
BUILD_TIME ?= "$(shell date +'%d/%m/%YT%H:%M:%S%z')"

DEFAULT: test

build:
	@go build -i -o ${REPONAME} ${CMD_PATH}

install:
	@echo "=> Installing dependencies"
	@dep ensure

push-tag:
	@echo "=> New tag version: ${VERSION}"
	git checkout ${BRANCH}
	git pull origin ${BRANCH}
	git tag ${VERSION}
	git push origin ${BRANCH} ${VERSION}

release:
	rm -rf dist
	@GITHUB_TOKEN=${VIDSY_GOBOT_GITHUB_TOKEN}  goreleaser

run: build
	@./${REPONAME} ${ARGS}

test:
	@go test "${PACKAGES}" -cover

vet:
	@go vet "${PACKAGES}"

