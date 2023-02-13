package main

import (
	"encoding/json"
	"net/http"
	"w0/data"
	"w0/database"
	"w0/pub"
	"w0/sub"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var m_orders map[string]data.Order

func getOrder(ctx *gin.Context) {
	uid := ctx.Query("order_uid")
	order := m_orders[uid]
	j_order, err := json.Marshal(order)
	if err != nil {
		logrus.Error(err)
	}
	_, err = ctx.Writer.Write(j_order)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func main() {

	pool := database.PoolConnect()

	orders := database.RecoverCash(pool)

	m_orders = make(map[string]data.Order)

	for _, v := range orders {
		m_orders[v.OrderUID] = v
	}

	router := gin.Default()
	router.LoadHTMLFiles("interface/index.html")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "main page",
		})
	})
	router.GET("/getorder/", getOrder)

	go sub.StertSub(pool)
	go pub.StartPub()

	router.Run()

}
