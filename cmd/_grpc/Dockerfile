FROM golang:alpine

ENV GOPATH=

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./server/main.go .
COPY ./proto/. proto/.

RUN go build -o main .

WORKDIR /dist

RUN cp /go/main /dist

EXPOSE 5300

CMD ["./main"]
