FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:3.17  

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
