FROM golang:1.13.4-alpine3.10 as builder
RUN mkdir /go/src/bundle-server
RUN apk add git
RUN go get github.com/gorilla/mux
COPY . /go/src/bundle-server
WORKDIR /go/src/bundle-server
RUN go build -o /go/bin/bundle-server .

FROM alpine:latest
RUN mkdir /bundle-server
WORKDIR /bundle-server
COPY --from=builder /go/bin/bundle-server .
RUN chmod 755 bundle-server
CMD ["./bundle-server"]