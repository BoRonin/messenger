FROM golang:1.21.1-alpine

RUN mkdir /app

COPY . /app

WORKDIR /app

CMD ["go", "run", "cmd/server/main.go"]
