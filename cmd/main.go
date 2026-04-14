package main

import (
	"spread_service/spreads"

	"github.com/gin-gonic/gin"
)

var symbols = []string{
	"EURUSD",
	"EURCAD",
	"USDJPY",
	"BTCUSD",
	"XAUUSD",
}
var spreadHandler *spreads.SpreadHandler

func main() {
	spreadHandler = spreads.NewHandler(symbols)

	r := gin.Default() // gin rest api handler

	r.GET("/symbols", spreadHandler.GetSymbols)
	r.GET("/spreads/:symbol", spreadHandler.GetSpread)
	r.PATCH("/spreads/:symbol", spreadHandler.SetSpread)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
