FROM golang:1.20.5-bullseye
RUN apt-get update && apt-get -y install gosu git
WORKDIR /go/app
COPY ./src /go/app
COPY ./docker/go/entrypoint-go.sh /usr/local/bin/entrypoint-go.sh
RUN chmod +x /usr/local/bin/entrypoint-go.sh
ENTRYPOINT ["/usr/local/bin/entrypoint-go.sh"]
