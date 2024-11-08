ARG GO_VERSION=1.23

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

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
