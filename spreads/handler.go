package spreads

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type SpreadHandler struct {
	symbols []string
	spreads map[string]*Spread
	mu      sync.RWMutex
}

func NewHandler(symbols []string) *SpreadHandler {
	spreads := make(map[string]*Spread)
	for _, s := range symbols {
		spreads[s] = nil
	}
	return &SpreadHandler{
		symbols: symbols,
		spreads: spreads,
	}
}

func (h *SpreadHandler) GetSpread(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if spread, ok := h.spreads[c.Param("symbol")]; ok {
		if spread.Spread >= 0 {
			c.JSON(http.StatusOK, spread)
		} else {
			c.JSON(http.StatusNotFound, "Spread not set")
		}
	} else {
		c.JSON(http.StatusNotFound, "Symbol not listed")
	}
}

func (h *SpreadHandler) GetSymbols(c *gin.Context) {
	c.JSON(http.StatusOK, h.symbols)
}

func (h *SpreadHandler) SetSpread(c *gin.Context) {
	var body struct {
		Spread float64 `json:"spread"`
	} //expected from POST request
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request body")
		return
	}
	if body.Spread <= 0 {
		c.JSON(http.StatusBadRequest, "Bad request: spread should be > 0")
		return
	}

	symbol := c.Param("symbol")
	h.mu.Lock()
	defer h.mu.Unlock()
	if spread, ok := h.spreads[symbol]; ok {
		spread.UpdatedAt = time.Now()
		spread.Spread = body.Spread
	} else {
		c.JSON(http.StatusNotFound, "Symbol not found")
		return
	}

	c.Status(http.StatusOK)
}
