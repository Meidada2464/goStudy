package testifyTest

import (
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
)

type HttpSuite struct {
	suite.Suite
	r   *httptest.ResponseRecorder
	mux *http.ServeMux
}

func (t *HttpSuite) SetupTest() {
	t.r = httptest.NewRecorder()
	t.mux = http.NewServeMux()
	t.mux.HandleFunc("/", index)
	t.mux.HandleFunc("/greeting", great)
}

func (t *HttpSuite) TestIndex() {
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return
	}
	t.mux.ServeHTTP(t.r, request)
	t.Assert().Equal(t.r.Code, 200, "")
}
