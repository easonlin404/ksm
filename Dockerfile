FROM golang:1.15

WORKDIR /go/src/github.com/easonlin404/ksm

ADD . /go/src/github.com/easonlin404/ksm

RUN go get -t -v ./...

CMD ["go", "run", "example/basic.go"]

EXPOSE 8080
