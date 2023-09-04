FROM golang:latest

USER root

WORKDIR /

COPY . .

CMD ["go run app.go"]