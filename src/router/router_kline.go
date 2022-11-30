package router

import (
	"demo/src/logger"
	"demo/src/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func getKLine(c *gin.Context) {
	var params *kLineParams
	err := c.ShouldBind(&params)
	if err != nil {
		logger.Error("klines bind failed", err)
		c.JSON(200, err)
	}
	if params.Size == 0 {
		params.Size = 10
	}
	if params.Begin == "" {
		params.Begin = strconv.FormatInt(time.Now().Add(-10*time.Minute).Unix(), 10)
	}

	gets, err := models.KLinesGets(params.Symbol, params.Period, params.Begin, params.Page, params.Size)
	if err != nil {
		logger.Error("klines query failed", err)
		c.JSON(200, err)
	}
	c.JSON(200, gets)
}

type kLineParams struct {
	Symbol string `json:"symbol" form:"symbol"`
	Period string `json:"period" form:"period"`
	Page   int    `json:"page" form:"page"`
	Size   int    `json:"size" form:"size"`
	Begin  string `json:"begin" form:"begin"`
}
