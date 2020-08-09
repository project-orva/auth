FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY /src .

RUN go mod download

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

EXPOSE 5258

CMD ["/dist/main"]