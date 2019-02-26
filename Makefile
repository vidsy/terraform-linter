BRANCH ?= "master"
REPONAME ?= "terraform-linter"
VERSION ?= $(shell cat ./VERSION)
PACKAGES ?= "./..."
CMD_PATH ?= "github.com/vidsy/terraform-linter/cmd/terraform/linter"

DEFAULT: test

build:
	@go build -i -o ${REPONAME} ${CMD_PATH}

build-image:
	@docker build -t vidsyhq/${REPONAME} --build-arg VERSION=${VERSION} .

docker-login:
	@docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}

install:
	@echo "=> Installing dependencies"
	@dep ensure

push-tag:
	@echo "=> New tag version: ${VERSION}"
	git checkout ${BRANCH}
	git pull origin ${BRANCH}
	git tag ${VERSION}
	git push origin ${BRANCH} ${VERSION}

push-to-docker-hub: docker-login
	@docker tag vidsyhq/${REPONAME}:latest vidsyhq/${REPONAME}:${CIRCLE_TAG}
	@docker push vidsyhq/${REPONAME}:${CIRCLE_TAG}
	@docker push vidsyhq/${REPONAME}

release:
	rm -rf dist
	@GITHUB_TOKEN=${VIDSY_GOBOT_GITHUB_TOKEN}  goreleaser

run: build
	@./${REPONAME} ${ARGS}

test:
	@go test "${PACKAGES}" -cover

vet:
	@go vet "${PACKAGES}"

