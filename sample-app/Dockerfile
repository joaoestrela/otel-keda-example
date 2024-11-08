ARG GO_VERSION=1.23

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build

ARG PROTOBUF_VERSION=24.4-r1
ARG PROTOC_GEN_GO_VERSION=1.35.1
ARG PROTOC_GEN_GO_GRPC_VERSION=1.5.1

RUN apk add --no-cache protobuf=${PROTOBUF_VERSION}
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION} \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GO_GRPC_VERSION}

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN protoc --go_out=. --go-grpc_out=. counter.proto

ARG TARGETOS TARGETARCH

RUN GOOS="$TARGETOS" GOARCH="$TARGETARCH" CGO_ENABLED=0 go build \
		-installsuffix 'static' \
		-o ./tmp/server \
		-a .

FROM gcr.io/distroless/static-debian12:nonroot AS final
USER nonroot

WORKDIR /

COPY --from=build --chown=root:root --chmod=755 /src/tmp/server /

ENTRYPOINT ["/server"]
