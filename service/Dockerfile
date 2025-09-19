FROM ubuntu:22.04
WORKDIR /app

ARG BIN_FILE=service.linux
COPY ${BIN_FILE} /app/service.linux
COPY .env /app/.env

RUN mkdir -p /configs
COPY configs/config.yaml /configs/config.yaml

EXPOSE 8000 9000
CMD ["/app/service.linux","-conf","/configs"]