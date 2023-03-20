FROM --platform=linux/amd64 golang:1.16-alpine as builder
RUN apk update && apk add git
COPY . /go/src/github.com/robbertnoordzij/nerve-centre-webhook
WORKDIR /go/src/github.com/robbertnoordzij/nerve-centre-webhook
ENV CGO_ENABLED 0
RUN go get ./...
RUN go vet ./... && \
    go test ./... && \
    go build

FROM --platform=linux/amd64 alpine:3.8
COPY --from=builder /go/src/github.com/robbertnoordzij/nerve-centre-webhook/nerve-centre-webhook \
	/usr/local/bin/nerve-centre-webhook

ENTRYPOINT [ "/usr/local/bin/nerve-centre-webhook" ]
CMD [ ]