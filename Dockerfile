FROM alpine:3.12

RUN apk add --no-cache git go make musl-dev

ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV PATH=/go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

WORKDIR $GOPATH/src

COPY ./src .

RUN go build -o /go/bin/app .

ENTRYPOINT ["/go/bin/app"]
