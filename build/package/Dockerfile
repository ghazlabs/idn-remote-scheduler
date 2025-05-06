FROM golang:1.23.8-alpine3.20 AS build
WORKDIR /go/src/github.com/ghazlabs/wa-scheduler

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./cmd ./cmd
COPY ./internal ./internal

WORKDIR /go/src/github.com/ghazlabs/wa-scheduler/cmd
RUN go build -o server

FROM alpine:3.20
RUN apk add ca-certificates tzdata

COPY --from=build /go/src/github.com/ghazlabs/wa-scheduler/cmd/server ./server
COPY --from=build /go/src/github.com/ghazlabs/wa-scheduler/cmd/web ./web

ENTRYPOINT [ "./server" ]
