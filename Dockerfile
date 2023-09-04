FROM golang:latest

USER 10001

WORKDIR /home/10001

COPY . .

USER root

RUN chmod +x /home/10001/start.sh

USER 10001

CMD ["./start.sh"]