FROM golang:latest

RUN useradd -m -s /bin/bash -u 10001 10001

USER 10001

WORKDIR /home/10001

COPY . .

USER root

RUN chmod +x /home/10001/start.sh

RUN chmod +x /home/10001/app

USER 10001

CMD ["./start.sh"]