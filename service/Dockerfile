FROM ubuntu:22.04
WORKDIR /
ADD alertor.linux .
ADD .env .
RUN mkdir configs
WORKDIR /configs
ADD config.yaml .
EXPOSE 8000
EXPOSE 9000
WORKDIR /
CMD ["./service.linux","-conf","/configs"]