FROM golang:latest

USER 10001

WORKDIR /home/10001

COPY . .

ENV PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

CMD ["go run app.go"]