FROM golang:1.23.2

ADD . /project
WORKDIR /project
RUN make
EXPOSE 8080
ENTRYPOINT ["/project/bin/client"]

