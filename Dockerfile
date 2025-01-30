FROM golang:1.23.4-alpine3.20 AS build-stage

WORKDIR /app

COPY . ./

ENV GOARCH=amd64
ENV GOOS=linux
ENV CGO_ENABLED=0

RUN go build -ldflags="-s -w" -trimpath -o /api ./cmd/api


FROM gcr.io/distroless/base-debian12 AS release-stage

WORKDIR /

COPY --from=build-stage /api /bin/api

USER nonroot:nonroot

ENTRYPOINT ["/bin/api"]