FROM alpine:3.12

RUN apk add --no-cache git make musl-dev \
    && wget https://golang.org/dl/go1.24.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz \
    && rm go1.24.0.linux-amd64.tar.gz

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

WORKDIR $GOPATH/src

COPY ./src .

RUN go build -o /go/bin/app .

ENTRYPOINT ["/go/bin/app"]
