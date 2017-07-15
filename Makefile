GO_LINKER_SYMBOL := main.version
GOOS := linux
GOARCH := amd64

glv:
	$(eval GO_LINKER_VALUE := $(shell git describe --tags --always))

ldflags: glv
	$(eval LDFLAGS := -ldflags "-X ${GO_LINKER_SYMBOL}=${GO_LINKER_VALUE}")

ver: glv
	$(eval VERSION := $(shell echo ${GO_LINKER_VALUE} | sed s/^v//))

build: glv ldflags ver
	GOOS=${GOOS} GOARCH=${GOARCH} go build ${LDFLAGS} -v -o tq-${VERSION}-${GOOS}-${GOARCH}