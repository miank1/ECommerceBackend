PROTO_DIR=proto
GO_OUT_DIR=proto

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

.PHONY: clean
clean:
	rm -f $(GO_OUT_DIR)/*.pb.go
