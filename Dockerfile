# docker build "$PWD" --build-arg commit="$(git rev-parse --short HEAD)" --build-arg version=v1.0.0 -t a-wasm/awasm:1.0.0
# docker build "$PWD" --build-arg commit="$(git rev-parse --short HEAD)" --build-arg version="$(git rev-parse --short HEAD)" -t a-wasm/awasm-pre-release:"$(git rev-parse --short HEAD)"

# Stage 1: modules caching
FROM golang:1.21-alpine as modules
LABEL maintainer="htquangg@gmail.com"

WORKDIR /a-wasm

COPY go.* .

RUN go mod download

# Stage 2: build
FROM golang:1.21-alpine as builder
LABEL maintainer="htquangg@gmail.com"

ARG version
ARG commit

ENV GOOSE linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
ENV GO111MODULE=on

WORKDIR /a-wasm

COPY --from=modules /go/pkg /go/pkg
COPY . .

RUN go build -o ./awasm -trimpath -ldflags "-s -w -X main.version=$version -X main.commitID=$commit" . \
    && cp ./config/config.docker.yaml ./config.yaml

RUN chmod 755 awasm
RUN cp ./awasm /usr/bin/awasm

# Stage 3: deploy
FROM alpine:3 as runtime
LABEL maintainer="htquangg@gmail.com"

ARG version

LABEL version=$version

WORKDIR /a-wasm

RUN apk update \
    && apk --no-cache add \
        bash \
    && echo "UTC" > /etc/timezone

ENV TZ UTC
ENV CONFIG_PATH /a-wasm/config.yaml

COPY --from=builder /usr/bin/awasm /usr/bin/awasm
COPY --from=builder /a-wasm/i18n ./i18n
COPY --from=builder /a-wasm/config.yaml ./config.yaml

EXPOSE 3000

CMD ["awasm", "run"]
