VERSION 0.7

ARG --global GOLANG_VERSION=1.19.7
ARG --global AMTRPC_VERSION=v2.6.0
ARG --global GOLINT_VERSION=v1.51.2
ARG --global IMAGE_REPOSITORY=quay.io/kairos-io/provider-amt

builder:
    FROM golang:$GOLANG_VERSION

    RUN apt update \
        && apt install -y upx git \
        && wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin $GOLINT_VERSION

    WORKDIR /build

    COPY go.mod go.sum ./

    RUN go mod download

version:
    FROM +builder

    COPY .git .git

    RUN --no-cache echo $(git describe --always --tags --dirty) > VERSION

    ARG VERSION=$(cat VERSION)
    SAVE ARTIFACT VERSION VERSION

amt-rpc-lib:
    FROM +builder

    GIT CLONE --branch $AMTRPC_VERSION git@github.com:open-amt-cloud-toolkit/rpc-go.git .

    RUN go build -buildmode=c-shared -o librpc.so ./cmd

    SAVE ARTIFACT librpc.so
    SAVE ARTIFACT librpc.h

build:
    FROM +builder

    ENV CGO_ENABLED=1

    COPY cmd cmd
    COPY pkg pkg

    COPY +amt-rpc-lib/librpc.so  /usr/local/lib/librpc.so
    COPY +amt-rpc-lib/librpc.h  /usr/local/include/librpc.h

    RUN go build -o provider-amt cmd/main.go && upx provider-amt

    SAVE ARTIFACT provider-amt AS LOCAL artifacts/provider-amt

image:
    FROM +version

    ARG VERSION=$(cat VERSION)

    FROM scratch

    COPY +amt-rpc-lib/librpc.so  /usr/local/lib/librpc.so
    COPY +amt-rpc-lib/librpc.h  /usr/local/include/librpc.h
    COPY +build/provider-amt /system/providers/provider-amt

    SAVE IMAGE --push $IMAGE_REPOSITORY:$VERSION

test:
    FROM +builder

    COPY cmd cmd
    COPY pkg pkg

    RUN go test -v -tags fake ./...

lint:
    FROM +builder

    COPY cmd cmd
    COPY pkg pkg

    RUN golangci-lint run