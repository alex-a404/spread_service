package spreads

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// setup gin to unit test endpoints
func setupRouter() (*gin.Engine, *SpreadHandler) {
	gin.SetMode(gin.TestMode)
	symbols := []string{"EURUSD", "EURCAD", "USDJPY", "BTCUSD", "XAUUSD"}
	h := NewHandler(symbols)
	r := gin.Default()
	r.GET("/symbols", h.GetSymbols)
	r.GET("/spreads/:symbol", h.GetSpread)
	r.PATCH("/spreads/:symbol", h.SetSpread)
	return r, h
}

func TestGetSymbols(t *testing.T) {
	r, _ := setupRouter()

	req, _ := http.NewRequest("GET", "/symbols", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var res []string
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("invalid json response")
	}

	if len(res) != 5 {
		t.Fatalf("expected 5 symbols, got %d", len(res))
	}
}

func TestSetSpread_InvalidSpread_400(t *testing.T) {
	r, _ := setupRouter()
	body := []byte(`{"spread": abcxyz}`)
	req, _ := http.NewRequest("PATCH", "/spreads/EURUSD", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSetSpread_InvalidSymbol_404(t *testing.T) {
	r, _ := setupRouter()
	body := []byte(`{"spread": 1}`)
	req, _ := http.NewRequest("PATCH", "/spreads/ABCXYZ", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestSetSpread_Correct_200(t *testing.T) {
	r, _ := setupRouter()
	body := []byte(`{"spread": 1}`)
	req, _ := http.NewRequest("PATCH", "/spreads/EURUSD", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetSpread_NotSet_404(t *testing.T) {
	r, _ := setupRouter()
	req, _ := http.NewRequest("GET", "/spreads/ABCXYZ", nil) // spread for symbol that does not exist
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetSpread_Full_200(t *testing.T) {
	r, _ := setupRouter()
	// set spread
	body := []byte(`{"spread": 5}`)
	req, _ := http.NewRequest("PATCH", "/spreads/EURUSD", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// get spread
	req, _ = http.NewRequest("GET", "/spreads/EURUSD", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// decode json resp
	var resp Spread
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if resp.Spread != 5 {
		t.Fatalf("expected spread 5, got %f", resp.Spread)
	}

	if resp.Symbol != "EURUSD" {
		t.Fatalf("expected symbol EURUSD, got %s", resp.Symbol)
	}

	if resp.UpdatedAt.IsZero() {
		t.Fatalf("expected UpdatedAt to be set")
	}
}
