FROM golang:1.18 as build-env

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o /go/bin/app

FROM gcr.io/distroless/static

COPY --from=build-env /go/bin/app /
EXPOSE 8080

ENTRYPOINT ["/app"]
