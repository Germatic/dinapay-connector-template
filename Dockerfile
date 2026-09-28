FROM golang:1.24-alpine AS build
ARG SERVICE=dinapay-connector-template
ARG REPOSITORY=github.com/Germatic/dinapay-connector-template
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILT_AT=unknown
ARG ENVIRONMENT=container
WORKDIR /src
COPY go.mod ./
COPY . .
RUN MODULE="$(go list -m)" && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w \
    -X ${MODULE}/internal/buildinfo.Service=${SERVICE} \
    -X ${MODULE}/internal/buildinfo.Repository=${REPOSITORY} \
    -X ${MODULE}/internal/buildinfo.Version=${VERSION} \
    -X ${MODULE}/internal/buildinfo.Commit=${COMMIT} \
    -X ${MODULE}/internal/buildinfo.BuiltAt=${BUILT_AT} \
    -X ${MODULE}/internal/buildinfo.Environment=${ENVIRONMENT}" -o /connector ./cmd/connector

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /connector /connector
EXPOSE 8092
ENTRYPOINT ["/connector"]
