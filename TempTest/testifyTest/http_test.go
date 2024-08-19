package testifyTest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r", r)

	fmt.Println("index")
	// put something into w
	w.Write([]byte("index"))
}

func great(w http.ResponseWriter, r *http.Request) {
	fmt.Println("great")
}

// how to mock http request

//func TestHttp(t *testing.T) {
//	mux := http.NewServeMux()
//	// open source
//	mux.HandleFunc("/", index)
//	mux.HandleFunc("/greeting", great)
//
//	// open http server
//	server := &http.Server{
//		// address
//		Addr: ":9999",
//		// route
//		Handler: mux,
//	}
//	// open http server
//	if err := server.ListenAndServe(); err != nil {
//		t.Fatal(err)
//	}
//}

func TestIndexHttp(t *testing.T) {
	// this recorder is used to record the http response and not to router
	recorder := httptest.NewRecorder()
	// mock request
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/greeting", great)
	mux.ServeHTTP(recorder, request)

	contains := assert.Contains(t, recorder.Body.String(), "index")
	fmt.Println("recorder", recorder)
	fmt.Println("is contain", contains)
}
