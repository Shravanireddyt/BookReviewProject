package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	db "Project1/database"

	jd "Project1"

	pf "Project1/bookService"

	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) AddBook(ctx context.Context, req *pf.AddBookRequest) (*pf.AddBookResponse, error) {

	fmt.Println("In AddBook")

	var bk jd.BookReview
	bk.ID = req.BookInfo.Id
	bk.Name = req.BookInfo.Name
	bk.Author = req.BookInfo.Author
	bk.ShortDesc = req.BookInfo.ShortDesc
	bk.Reviews = []jd.Review{}

	status := db.AddBookToDB(bk)
	if status == pf.StatusType_UNKNOWN {
		return &pf.AddBookResponse{Status: http.StatusBadRequest, Msg: "The Book id already present in database"}, nil
	} else if status == pf.StatusType_DB_OPERATION_FAIL {
		return &pf.AddBookResponse{Status: http.StatusBadRequest, Msg: "Error in database operation"}, nil
	}
	return &pf.AddBookResponse{Status: http.StatusOK, Msg: "AddBook Success"}, nil
}

func (s *server) GetBook(ctx context.Context, req *pf.GetBookRequest) (*pf.GetBookResponse, error) {

	fmt.Println("In GetBook")

	var Response *pf.GetBookResponse
	var book []*pf.Book

	books, status := db.GetBookFromDB(req.Limit)

	var test bool = false
	if status == pf.StatusType_OP_SUCCESS {
		for _, bk := range books {
			book = append(book, &pf.Book{Name: bk.Name, Author: bk.Author, ShortDesc: bk.ShortDesc})
			test = true
		}
		if test == false {
			return &pf.GetBookResponse{Status: http.StatusNoContent, Msg: "Database empty"}, nil
		}
		Response = &pf.GetBookResponse{Status: http.StatusOK, Books: book}
	} else {
		Response = &pf.GetBookResponse{Status: http.StatusNoContent, Msg: "Database operation failed"}
	}
	return Response, nil
}

func (s *server) AddReview(ctx context.Context, req *pf.AddReviewRequest) (*pf.AddReviewResponse, error) {

	fmt.Println("In AddReview")
	var rv jd.Review
	rv.Name = req.Review.Name
	rv.Score = uint32(req.Review.Score)
	rv.Text = req.Review.Text

	status := db.AddReviewToDB(req.Id, rv)
	if status == pf.StatusType_UNKNOWN {
		return &pf.AddReviewResponse{Status: http.StatusNoContent, Msg: "Book Id not present in database"}, nil
	} else if status == pf.StatusType_DB_OPERATION_FAIL {
		return &pf.AddReviewResponse{Status: http.StatusNotAcceptable, Msg: "Error in Database operation"}, nil
	}
	Response := &pf.AddReviewResponse{Status: http.StatusOK, Msg: "AddReview Success"}
	return Response, nil
}

func (s *server) GetReview(ctx context.Context, req *pf.GetReviewRequest) (*pf.GetReviewResponse, error) {

	fmt.Println("In GetReview")

	var Response *pf.GetReviewResponse
	var reviews []*pf.Review

	rvs, status := db.GetReviewFromDB(req.BookId)
	if status == pf.StatusType_UNKNOWN {
		return &pf.GetReviewResponse{Status: http.StatusNoContent, Msg: "Book Id not present in database"}, nil
	} else if status == pf.StatusType_DB_OPERATION_FAIL {
		return &pf.GetReviewResponse{Status: http.StatusNotAcceptable, Msg: "Error in Database operation"}, nil
	} else if status == pf.StatusType_NOT_FOUND {
		return &pf.GetReviewResponse{Status: http.StatusOK, Msg: "The reviews not present for this book id"}, nil
	}

	for _, rv := range rvs {
		reviews = append(reviews, &pf.Review{Name: rv.Name, Score: rv.Score, Text: rv.Text})
	}
	Response = &pf.GetReviewResponse{Status: http.StatusOK, Reviews: reviews}
	return Response, nil
}

func main() {

	db.CBInitialize()
	address := "0.0.0.0:8090"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Listening failed")
		panic(err)
	}
	fmt.Printf("Server is listening on %v ", address)

	s := grpc.NewServer()
	pf.RegisterBookServiceServer(s, &server{})

	s.Serve(lis)

}
