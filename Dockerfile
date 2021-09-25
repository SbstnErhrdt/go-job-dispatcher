# syntax=docker/dockerfile:1

### build go app
FROM golang:1.16-alpine as go
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY . ./
# set env
RUN export CGO_ENABLED=0
RUN export GOOS=linux
RUN export GOARCH=amd64
# download and build
RUN go mod download
RUN go build -o main


### get certs
FROM alpine:latest as certs
RUN apk --update add ca-certificates


### combine all
FROM scratch
ENV PATH=/bin
# copy certs
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# copy go app
WORKDIR /app
COPY --from=go /app /app
# start app
CMD ["./main"]