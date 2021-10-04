package database

import (
	"fmt"
	"sync"

	pf "Project1/bookService"

	jd "Project1"

	"github.com/couchbase/gocb/v2"
)

var CollectionH syncedcbBucketHandle

type syncedcbBucketHandle struct {
	bucketLock       sync.Mutex
	scopeHandler     *gocb.Scope
	CollectionHandle *gocb.Collection
}

const (
	BUCKET_NAME     string = "BookReview_name"
	SCOPE_NAME      string = "bookdata"
	COLLECTION_NAME string = "BookInfo"
)

func CBInitialize() {

	cluster, err := gocb.Connect(
		"localhost",
		gocb.ClusterOptions{
			Username: "Administrator",
			Password: "Administrator",
		})
	if err != nil {
		panic(err)
	}

	bucketMgr := cluster.Buckets()
	if _, err := bucketMgr.GetBucket(BUCKET_NAME, nil); err != nil {
		err = bucketMgr.CreateBucket(gocb.CreateBucketSettings{
			BucketSettings: gocb.BucketSettings{
				Name:                 BUCKET_NAME,
				FlushEnabled:         false,
				ReplicaIndexDisabled: true,
				RAMQuotaMB:           150,
				NumReplicas:          0,
				BucketType:           gocb.CouchbaseBucketType,
			},
			ConflictResolutionType: gocb.ConflictResolutionTypeSequenceNumber,
		}, nil)
		if err != nil {
			panic(err)
		}
	}

	bucket := cluster.Bucket(BUCKET_NAME)
	bucketmgr := bucket.Collections()
	AllScopes, _ := bucketmgr.GetAllScopes(nil)
	var test bool = false
	for _, v := range AllScopes {
		if v.Name == SCOPE_NAME {
			test = true
		}
	}
	if test != true {
		if err := bucketmgr.CreateScope(SCOPE_NAME, nil); err != nil {
			panic(err)
		}
	}
	if AllScopes[0].Collections[0].Name != COLLECTION_NAME {
		if err := bucketmgr.CreateCollection(
			gocb.CollectionSpec{
				Name:      COLLECTION_NAME,
				ScopeName: SCOPE_NAME,
			}, nil); err != nil {
			panic(err)
		}
		query := fmt.Sprintf("create primary index on `%s`.%s.%s ", BUCKET_NAME, SCOPE_NAME, COLLECTION_NAME)
		_, err := bucket.Scope(SCOPE_NAME).Query(query, nil)
		if err != nil {
			panic("can't create index")
		}
	}
	CollectionH.scopeHandler = bucket.Scope(SCOPE_NAME)
	CollectionH.CollectionHandle = CollectionH.scopeHandler.Collection(COLLECTION_NAME)
}

func AddBookToDB(bk jd.BookReview) pf.StatusType {

	fmt.Println("In addbook couchbase")

	t := CheckForBookEntry(bk.ID)
	if t != false {
		return pf.StatusType_UNKNOWN
	}

	CollectionH.bucketLock.Lock()
	_, err2 := CollectionH.CollectionHandle.Insert(fmt.Sprintf("%d", bk.ID), &bk, nil)
	CollectionH.bucketLock.Unlock()
	if err2 != nil {
		return pf.StatusType_DB_OPERATION_FAIL
	}
	return pf.StatusType_OP_SUCCESS
}

func GetBookFromDB(lt int32) ([]pf.Book, pf.StatusType) {

	fmt.Println("In getbook couchbase")

	var books []pf.Book
	query1 := fmt.Sprintf("select name,id,author,shortdesc from `%s`.%s.%s limit $1", BUCKET_NAME, SCOPE_NAME, COLLECTION_NAME)
	CollectionH.bucketLock.Lock()
	rows, err := CollectionH.scopeHandler.Query(query1, &gocb.QueryOptions{PositionalParameters: []interface{}{lt}})
	CollectionH.bucketLock.Unlock()
	if err != nil {
		return books, pf.StatusType_DB_OPERATION_FAIL
	}

	for rows.Next() {
		var bk pf.Book
		err := rows.Row(&bk)
		if err != nil {
			return books, pf.StatusType_DB_OPERATION_FAIL
		}
		books = append(books, bk)
	}
	return books, pf.StatusType_OP_SUCCESS
}

func AddReviewToDB(id uint32, rv jd.Review) pf.StatusType {

	fmt.Println("In addreview couchbase")

	t := CheckForBookEntry(id)
	if t == false {
		return pf.StatusType_UNKNOWN
	}

	mops := []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("reviews", rv, nil),
	}
	CollectionH.bucketLock.Lock()
	_, err2 := CollectionH.CollectionHandle.MutateIn(fmt.Sprintf("%d", id), mops, &gocb.MutateInOptions{})
	CollectionH.bucketLock.Unlock()
	if err2 != nil {
		return pf.StatusType_DB_OPERATION_FAIL
	}
	return pf.StatusType_OP_SUCCESS
}

func GetReviewFromDB(id uint32) ([]pf.Review, pf.StatusType) {

	fmt.Println("In get review couchbase")
	var reviews []pf.Review

	t := CheckForBookEntry(id)
	if t == false {
		return reviews, pf.StatusType_UNKNOWN
	}

	query1 := fmt.Sprintf("select reviews from `%s`.%s.%s where id=$1", BUCKET_NAME, SCOPE_NAME, COLLECTION_NAME)
	CollectionH.bucketLock.Lock()
	rows, err1 := CollectionH.scopeHandler.Query(query1, &gocb.QueryOptions{PositionalParameters: []interface{}{id}})
	CollectionH.bucketLock.Unlock()
	if err1 != nil {
		return reviews, pf.StatusType_DB_OPERATION_FAIL
	}

	var rv jd.ReviewArr
	for rows.Next() {
		err := rows.Row(&rv)
		if err != nil {
			return reviews, pf.StatusType_DB_OPERATION_FAIL
		}
	}

	var test bool = false
	for _, r := range rv.Reviews {
		reviews = append(reviews, pf.Review{Name: r.Name, Score: r.Score, Text: r.Text})
		test = true
	}
	if test == false {
		return reviews, pf.StatusType_NOT_FOUND
	}

	return reviews, pf.StatusType_OP_SUCCESS
}

func CheckForBookEntry(id uint32) bool {

	query1 := fmt.Sprintf("select name,author,shortdesc from `%s`.%s.%s  where id=$1", BUCKET_NAME, SCOPE_NAME, COLLECTION_NAME)
	CollectionH.bucketLock.Lock()
	rows, _ := CollectionH.scopeHandler.Query(query1, &gocb.QueryOptions{PositionalParameters: []interface{}{id}})
	CollectionH.bucketLock.Unlock()
	var br jd.Book
	err := rows.One(&br)
	if err != nil {
		return false
	}
	return true

}
