VERSION 0.7

ARG --global GOLANG_VERSION=1.19.7
ARG --global AMTRPC_VERSION=v2.6.0
ARG --global GOLINT_VERSION=v1.51.2
ARG --global IMAGE_REPOSITORY=ghcr.io/kairos-io/openamt

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

    RUN --no-cache echo $(git describe --always --tags) > VERSION

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

    RUN go build -o agent-provider-amt cmd/main.go && upx agent-provider-amt

    SAVE ARTIFACT agent-provider-amt AS LOCAL artifacts/agent-provider-amt
    SAVE ARTIFACT /usr/local/lib/librpc.so AS LOCAL artifacts/librpc.so

image:
    FROM +version

    ARG VERSION=$(cat VERSION)

    FROM scratch

    COPY --chmod 0777 +amt-rpc-lib/librpc.so  librpc.so
    COPY --chmod 0777 +build/agent-provider-amt agent-provider-amt

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