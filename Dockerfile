FROM golang:1.23.4-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o cell-router .

EXPOSE 8080

CMD ["./cell-router"]
