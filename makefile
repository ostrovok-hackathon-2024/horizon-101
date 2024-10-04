all: bin/client bin/server
bin/client: proto
bin/server: proto
proto: yea.proto
	mkdir -p $@
	protoc --go_out=proto --go_opt=paths=source_relative \
		--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
		$<
	touch $@

bin/%: cmd/%/*.go
	mkdir -p $(@D)
	go build -o $@ ./$(<D)

clean:
	rm -rf bin proto
