FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./pricecollector ./pricecollector
COPY ./lib ./lib

ENV GOOS=linux
ENV CGO_ENABLED=0

RUN go build \
  -o main ./pricecollector/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Bangkok

COPY --from=builder /app/main .

CMD ["/app/main"]
