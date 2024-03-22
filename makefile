create-pb:
	protoc --go_out=. --go-grpc_out=. proto/greet.proto
gen-ssl-cert:
	openssl genrsa -out ssl/server.key 2048
	openssl req -nodes -new -x509 -sha256 -days 1825 -config ssl/cert.conf -extensions 'req_ext' -key ssl/server.key -out ssl/server.crt
run-server:
	go run server/main.go
run-client:
	go run client/main.go