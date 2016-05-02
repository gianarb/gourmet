package api

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	l := log.Logger{}
	r := http.Request{}
	w := httptest.NewRecorder()
	h := PingHandler(&l)
	h(w, &r)
	if "{\"Status\":\"ok\"}" != fmt.Sprintf("%s", w.Body) {
		t.Fail()
	}
}
