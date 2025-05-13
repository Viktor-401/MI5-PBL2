FROM golang:1.24

WORKDIR /go/src/app

COPY . .

EXPOSE 8080

RUN go build -o main cmd/main.go

CMD ["./main"]