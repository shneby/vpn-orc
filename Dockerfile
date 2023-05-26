FROM golang:1.20

WORKDIR /app

COPY . .

RUN go build -o app

CMD ["./app"]

EXPOSE 8080/tcp
