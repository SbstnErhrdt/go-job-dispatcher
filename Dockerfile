# syntax=docker/dockerfile:1

### get certs
FROM alpine:latest as certs
RUN apk --update add ca-certificates

### build go app
FROM golang:1.16-alpine as builder
RUN apk add --no-cache gcc musl-dev
RUN mkdir /build
WORKDIR /build
COPY . ./
# download and build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

### combine all
FROM scratch
ENV PATH=/bin
# copy certs
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# copy go app
COPY --from=builder /build .
# start app
# executable
ENTRYPOINT [ "./main" ]
# arguments that can be overridden
CMD [ "" ]