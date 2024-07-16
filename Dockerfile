ARG GOLANG_VERSION=1.22.5
FROM --platform=$BUILDPLATFORM golang:${GOLANG_VERSION} AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETOS TARGETARCH
RUN --mount=target=. --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /dist-main ./cmd/main
RUN --mount=target=. --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /dist-migrate ./cmd/migrate

FROM gcr.io/distroless/base:nonroot

USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /dist-main ./
COPY --from=build --chown=nonroot:nonroot /dist-migrate ./
