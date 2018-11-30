.PHONY: build install snapshot dist test vet lint fmt run clean docker
OUT := ngxcpd
PKG := github.com/maetthu/ngxcpd
VERSION := $(shell git describe --always --dirty --tags)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: build

build:
	CGO_ENABLED=0 GOOS=linux go build -a -v -o ${OUT} ${PKG}

install:
	CGO_ENABLED=0 GOOS=linux go install -a -v -o ${OUT} ${PKG}

snapshot:
	goreleaser --snapshot --rm-dist

dist:
	goreleaser --rm-dist

test:
	@go test -v -tags=fixtures ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

fmt:
	@gofmt -l -w -s ${GO_FILES}

run: build
	./${OUT} listen

clean:
	-@rm ${OUT}

docker: 
	docker build \
		-t maetthu/ngxcpd:${VERSION} \
		-t maetthu/ngxcpd:latest\
		.
	docker push maetthu/ngxcpd:${VERSION}
	docker push maetthu/ngxcpd:latest


