FROM golang:1.21-bullseye as build

WORKDIR /code

ENV CGO_ENABLED=1
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build make

FROM gcr.io/distroless/base-debian11

COPY --from=build /code/bin/* /

ENV CGO_ENABLED=1
ENTRYPOINT ["/ipinfo-server"]
