FROM golang:latest

USER 10001

WORKDIR /home/10001

COPY . .

CMD ["go run app.go"]