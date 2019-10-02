# FROM golang:latest as builder
# ENV GO111MODULE=auto
# ENV WORK_ENV=production
# RUN echo $GOPATH
# RUN echo $WORK_ENV
# RUN echo %GOROOT%
# ADD . /go/src/igfeed
# WORKDIR /go/src/igfeed
# RUN go mod download
# # WORKDIR /go/src
# # RUN go build .
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

# FROM alpine:latest
# WORKDIR /go/src/igfeed
# COPY --from=builder /go/src/igfeed /go/src/igfeed/
# CMD ["./main"]
FROM golang

ENV GO111MODULE=on

WORKDIR /app

COPY . .

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

RUN go build -o main .

CMD ["./main"]
