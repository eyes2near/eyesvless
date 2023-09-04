FROM golang:latest

USER 10001

WORKDIR /home/10001

COPY . .

RUN chmod +x ./start.sh

CMD ["./start.sh"]