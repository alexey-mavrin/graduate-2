ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

client:
	mkdir -p bin
	go build -o bin/client cmd/client/main.go

server:
	mkdir -p bin
	go build -o bin/server cmd/server/main.go

key:
	openssl req -new -subj "/C=RU/ST=SPb/CN=localhost" \
		-newkey rsa:2048 -nodes -keyout keys/server.key \
		-out keys/server.csr
	openssl x509 -req -days 365 -in keys/server.csr \
		-signkey keys/server.key -out keys/server.crt \
		-extfile keys/self-signed-cert.ext

db_clean:
	@rm -f $(ROOT_DIR)/cache_store.db
	@rm -f $(ROOT_DIR)/server_storage.db
	@rm -f $(ROOT_DIR)/internal/store/secret_storage.db
	@rm -f $(ROOT_DIR)/internal/server/secret_storage.db
	@rm -f $(ROOT_DIR)/internal/client/cache_storage.db
	@rm -f $(ROOT_DIR)/server_storage.db
	@rm -f $(ROOT_DIR)/cache_store.db
	@rm -f $(ROOT_DIR)/tests/cache_inttest_store.db
	@rm -f $(ROOT_DIR)/tests/server_inttest_storage.db
