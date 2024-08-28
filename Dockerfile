FROM golang:1.22.5-alpine

RUN apk add --no-cache gcc musl-dev sqlite

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]
