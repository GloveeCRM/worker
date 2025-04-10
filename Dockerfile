FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go build -o worker cmd/worker/main.go

EXPOSE 7460

CMD ["./worker"]