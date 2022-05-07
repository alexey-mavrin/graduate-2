client:
	go build -o client cmd/client/main.go

server:
	go build -o server cmd/server/main.go

key:
	openssl req -new -subj "/C=RU/ST=SPb/CN=localhost" \
		-newkey rsa:2048 -nodes -keyout keys/server.key \
		-out keys/server.csr
	openssl x509 -req -days 365 -in keys/server.csr \
		-signkey keys/server.key -out keys/server.crt \
		-extfile keys/self-signed-cert.ext
