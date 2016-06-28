package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	r := http.Request{}
	w := httptest.NewRecorder()
	h := PingHandler()
	h(w, &r)
	if "{\"Status\":\"ok\"}" != fmt.Sprintf("%s", w.Body) {
		t.Fail()
	}
}
