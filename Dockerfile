FROM golang:latest

RUN useradd -m -s /bin/bash -u 10001 eyes

USER eyes

WORKDIR /home/eyes

COPY . .

CMD ["go run app.go"]