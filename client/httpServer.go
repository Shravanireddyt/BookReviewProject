package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	jd "Project1"
	pf "Project1/bookService"

	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var HttpServer *http.Server

func StartHttpServer() bool {

	mxr := mux.NewRouter()
	mxr.HandleFunc("/api/v1/book", HandleAddBookReq).Methods("POST")
	mxr.HandleFunc("/api/v1/book/{limit}", HandleGetBookReq).Methods("GET")
	mxr.HandleFunc("/api/v1/review", HandleAddReviewReq).Methods("POST")
	mxr.HandleFunc("/api/v1/review/{id}", HandleGetReviewReq).Methods("GET")

	h2s := &http2.Server{}
	HttpServer = &http.Server{
		Addr:    ":6000",
		Handler: h2c.NewHandler(mxr, h2s),
	}
	return true
}

func HandleAddBookReq(w http.ResponseWriter, req *http.Request) {

	fmt.Println("In Handle add book req")

	headerContentTtype := req.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		fmt.Println("content is not json")
		w.WriteHeader(http.StatusNotAcceptable)
	}
	bytes, _ := ioutil.ReadAll(req.Body)
	bk := jd.AddBookReq{}
	err := json.Unmarshal(bytes, &bk)
	if err != nil {
		fmt.Println("Error in Unmarshal")
	}

	addbookreq := &pf.AddBookRequest{BookInfo: &pf.Book{Id: bk.Id, Name: bk.Name, Author: bk.Author, ShortDesc: bk.ShortDesc}}
	if rcv, err1 := GrpcClient.AddBook(req.Context(), addbookreq); err1 == nil {
		json.NewEncoder(w).Encode(rcv)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func HandleGetBookReq(w http.ResponseWriter, req *http.Request) {

	fmt.Println("In HandleGetBookReq")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	lt, _ := strconv.Atoi(vars["limit"])

	getrequest := &pf.GetBookRequest{Limit: int32(lt)}

	getRsp, _ := GrpcClient.GetBook(context.Background(), getrequest)
	if getRsp != nil {
		fmt.Println("Receive response => ", getRsp)
		json.NewEncoder(w).Encode(getRsp)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func HandleAddReviewReq(w http.ResponseWriter, req *http.Request) {

	fmt.Println("In Handle Add Review req")

	headerContentTtype := req.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		fmt.Println("content is not json")
		w.WriteHeader(http.StatusNotAcceptable)
	}

	bytes, _ := ioutil.ReadAll(req.Body)
	addreviewreq := jd.AddReviewReq{}
	err := json.Unmarshal(bytes, &addreviewreq)
	if err != nil {
		fmt.Println("Error in Unmarshal")
	}
	addrequest := &pf.AddReviewRequest{Id: addreviewreq.Id, Review: &pf.Review{Name: addreviewreq.Name, Score: addreviewreq.Score, Text: addreviewreq.Text}}
	addRsp, _ := GrpcClient.AddReview(context.Background(), addrequest)
	if addRsp != nil {
		fmt.Println("Receive response => ", addRsp)
		json.NewEncoder(w).Encode(addRsp)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func HandleGetReviewReq(w http.ResponseWriter, req *http.Request) {

	fmt.Println("In HandleGetBookReq")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	lt, _ := strconv.Atoi(vars["id"])
	getrequest := &pf.GetReviewRequest{BookId: uint32(lt)}
	getRsp, _ := GrpcClient.GetReview(context.Background(), getrequest)
	if getRsp != nil {
		fmt.Println("Receive response => ", getRsp)
		json.NewEncoder(w).Encode(getRsp)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
