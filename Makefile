BRANCH ?= "master"
REPONAME ?= "terraform-linter"
VERSION ?= $(shell cat ./VERSION)
PACKAGES ?= "./..."
CMD_PATH ?= "github.com/vidsy/terraform-linter/cmd/terraform/linter"
BUILD_TIME ?= "$(shell date +'%d/%m/%YT%H:%M:%S%z')"

DEFAULT: test

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
	@GITHUB_TOKEN=${VIDSY_GOBOT_GITHUB_TOKEN} VERSION=${VERSION} BUILD_TIME=%${BUILD_TIME} CMD_PATH=${CMD_PATH} goreleaser

run:
	@go build -i -ldflags "-X ${CMD_PATH}/main.Version=${VERSION}-dev -X ${CMD_PATH}/main.BuildTime=17/01/2017T14:12:35+0000" -o ${REPONAME} ${CMD_PATH}
	@./${REPONAME} ${ARGS}

test:
	@go test "${PACKAGES}" -cover

vet:
	@go vet "${PACKAGES}"

