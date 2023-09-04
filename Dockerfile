FROM golang:latest

WORKDIR /

COPY . .

CMD ["go run app.go"]