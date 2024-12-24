
FROM golang:1.23.4-alpine

WORKDIR /app

COPY . .

RUN go build -o mock-cell .

CMD ["./mock-cell"]