package jsondata

type Book struct {
	Name      string   `json:"name"`
	Author    []string `json:"author"`
	ShortDesc string   `json:"shortDesc"`
}
type Review struct {
	Name  string `json:"name"`
	Score uint32 `json:"score"`
	Text  string `json:"text"`
}

type BookReview struct {
	ID        uint32   `json:"id"`
	Name      string   `json:"name"`
	Author    []string `json:"author"`
	ShortDesc string   `json:"shortdesc"`
	Reviews   []Review `json:"reviews"`
}

type ReviewArr struct {
	Reviews []Review `json:"reviews"`
}

type AddBookReq struct {
	Id        uint32   `json:"id"`
	Name      string   `json:"name"`
	Author    []string `json:"author"`
	ShortDesc string   `json:"shortDesc"`
}

type AddReviewReq struct {
	Id    uint32 `json:"id"`
	Name  string `json:"name"`
	Score uint32 `json:"score"`
	Text  string `json:"text"`
}

// curl --http2-prior-knowledge -H "Content-Type: application/json" -X POST --data @data.json http://localhost:6000/api/v1/book
// curl --http2-prior-knowledge -H "Content-Type: application/json" -X GET http://localhost:6000/api/v1/book/{2}
// curl -V --header "Content-Type: application/json" --request GET --data '{"id":2,"name":"abc","author":"xyz"},"shortdesc":"nice"}' http://localhost:8090/api/v1/book
// curl --http2-prior-knowledge -X POST --data @data.json http://localhost:6000/api/v1/book
// curl --http2-prior-knowledge  "Content-Type: application/json" --request POST --data '{"id":2,"name":"abc","author":"xyz"},"shortdesc":"nice"}' http://localhost:6000/api/v1/book
//protoc --proto_path=bookService/ --go_out=plugins=grpc:bookService/ bookprt.proto
