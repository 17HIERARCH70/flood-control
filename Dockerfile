FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/main .
COPY config/config-local.yaml ./config/

ENV CONFIG_PATH=/app/config/config-local.yaml
EXPOSE 8080

CMD ["./main"]
