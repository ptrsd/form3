FROM golang:latest
LABEL maintainer="Piotr Siudy"
RUN mkdir /app

WORKDIR /app

CMD ["test"]
ENTRYPOINT ["go"]