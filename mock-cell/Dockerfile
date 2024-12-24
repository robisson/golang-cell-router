
FROM golang:1.23.4-alpine

WORKDIR /app

COPY . .

RUN go build -o mock-cell .

EXPOSE 8080

CMD ["./mock-cell"]