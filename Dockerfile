FROM golang:1-alpine as builder

RUN apk update && \
    apk add curl git file && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/niranjan94/s3-archiver/
COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure -v -vendor-only
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" && file s3-archiver


FROM alpine
WORKDIR /data/s3-archiver
COPY --from=builder /go/src/github.com/niranjan94/s3-archiver/s3-archiver .
RUN apk add --no-cache ca-certificates && ln -s $(pwd)/s3-archiver /usr/bin/s3-archiver
ENTRYPOINT ["s3-archiver"]
