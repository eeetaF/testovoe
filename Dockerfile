FROM golang:1.23.7-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./src/main

CMD ["./server"]
