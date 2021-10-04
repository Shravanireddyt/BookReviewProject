Service file generation:
protoc --proto_path=bookService/ --go_out=plugins=grpc:bookService/ bookprt.proto

client execution:
go run client/client.go client/httpServer.go client/grpcClient.go

Server execution:
go run server/server.go 

Add Book:
curl --http2-prior-knowledge -H "Content-Type: application/json" -X POST --data @data.json http://localhost:6000/api/v1/book

Get Book:
curl --http2-prior-knowledge -H "Content-Type: application/json" -X GET http://localhost:6000/api/v1/book/{1}

Add Review
curl --http2-prior-knowledge -H "Content-Type: application/json" -X POST --data @review.json http://localhost:6000/api/v1/review

Get Review
curl --http2-prior-knowledge -H "Content-Type: application/json" -X GET http://localhost:6000/api/v1/review/{1}