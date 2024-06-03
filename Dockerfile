# Stage 1: modules caching
FROM golang:1.21-alpine as modules
LABEL maintainer="htquangg@gmail.com"

WORKDIR /awasm

COPY go.* .

RUN go mod download

# Stage 2: build
FROM golang:1.21-alpine as builder
LABEL maintainer="htquangg@gmail.com"

ENV CGO_ENABLED 0
ENV GO111MODULE=on

COPY --from=modules /go/pkg /go/pkg

WORKDIR /awasm

COPY . .

RUN apk --no-cache add build-base git \
    && make clean build

RUN cp ./awasm /usr/bin/awasm && cp ./config/awasm.docker.yaml ./awasm.yaml

# Stage 3: deploy
FROM alpine:3 as runtime
LABEL maintainer="htquangg@gmail.com"

WORKDIR /awasm

RUN apk update \
    && apk --no-cache add \
        bash \
    && echo "UTC" > /etc/timezone

ARG TZ
ENV TZ=${TZ:-"UTC"}

COPY --from=builder /usr/bin/awasm /usr/bin/awasm
COPY --from=builder /awasm/i18n ./i18n
COPY --from=builder /awasm/awasm.yaml ./awasm.yaml

RUN chmod 755 /usr/bin/awasm

EXPOSE 8080

CMD ["awasm", "run", "--config-path", "/awasm/awasm.yaml"]
