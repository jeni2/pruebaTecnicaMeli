FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app cmd/api/main.go


EXPOSE 8080

CMD ["./app"]
