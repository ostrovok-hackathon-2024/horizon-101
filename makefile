all: bin/client
bin/client: proto
proto: proto/yea.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$<
	touch $@

bin/%: cmd/%/*.go
	mkdir -p $(@D)
	go build -o $@ ./$(<D)

