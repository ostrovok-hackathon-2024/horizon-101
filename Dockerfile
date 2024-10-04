FROM golang:1.23.2

RUN apt-get update
RUN apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest 
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ADD . /project
WORKDIR /project
RUN make
ENTRYPOINT ["/project/entrypoint.sh"]
