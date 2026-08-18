package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEvaluateEndpoint_OK(t *testing.T) {
	srv := NewServer(nil)
	body := `{"goodsValue":10000,"gstPct":18,"freight":500,"loadingCoolie":200,"localTransport":300,"wastagePct":5,"units":100,"regular":false}`
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/evaluate", strings.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
	var r Result
	json.Unmarshal(rec.Body.Bytes(), &r)
	if !almost(r.LandedTotal, 12800) {
		t.Fatalf("landed=%v want 12800", r.LandedTotal)
	}
}

func TestEvaluateEndpoint_ValidationError(t *testing.T) {
	srv := NewServer(nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/evaluate", strings.NewReader(`{"units":0}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d want 400", rec.Code)
	}
}

func TestHistoryEndpoint_UnavailableWithoutDB(t *testing.T) {
	srv := NewServer(nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/history", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status %d want 503", rec.Code)
	}
}
