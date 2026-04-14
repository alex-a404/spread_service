package spreads

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// SpreadHandler a struct that tracks spreads for a given currency pair symbols manifest
type SpreadHandler struct {
	symbols map[string]struct{} //currency pairs as a set
	spreads map[string]*Spread
	mu      sync.RWMutex
}

// NewHandler create a SpreadHandler for a given manifest of symbols
func NewHandler(symbols []string) *SpreadHandler {
	spreads := make(map[string]*Spread, len(symbols))
	symbolSet := make(map[string]struct{}, len(symbols))

	for _, s := range symbols {
		symbolSet[s] = struct{}{}
		spreads[s] = nil
	}

	return &SpreadHandler{
		symbols: symbolSet,
		spreads: spreads,
	}
}

// GetSpread returns value of a given spread symbol
func (h *SpreadHandler) GetSpread(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	spread, ok := h.spreads[c.Param("symbol")]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Symbol not listed"})
		return
	}

	if spread == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Spread not set"})
		return
	}

	c.JSON(http.StatusOK, spread)
}

// GetSymbols returns symbols for which spreads are tracked
func (h *SpreadHandler) GetSymbols(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	symbolsAsList := make([]string, 0, len(h.symbols))
	for s := range h.symbols {
		symbolsAsList = append(symbolsAsList, s)
	}
	c.JSON(http.StatusOK, symbolsAsList)
}

// SetSpread sets a spread for a symbol and updates the updated time
func (h *SpreadHandler) SetSpread(c *gin.Context) {
	var body struct {
		Spread float64 `json:"spread"`
	} //expected from POST request

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if body.Spread <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request: spread should be > 0"})
		return
	}

	symbol := c.Param("symbol")
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.spreads[symbol]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Symbol not found"})
		return
	}

	h.spreads[symbol] = &Spread{
		Symbol:    symbol,
		Spread:    body.Spread,
		UpdatedAt: time.Now(),
	}

	c.Status(http.StatusOK)
}
