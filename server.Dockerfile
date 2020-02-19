FROM golang

RUN yarn

ADD . /go/src/github.com/kil0meters/acolyte

RUN go install github.com/kil0meters/cmd/acolyte

ENTRYPOINT /go/bin/acolyte

EXPOSE 8080

