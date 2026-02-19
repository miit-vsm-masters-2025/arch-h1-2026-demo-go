package main

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var counter uint64

	counterGauge := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "app_counter_value",
			Help: "Current in-memory counter value.",
		},
		func() float64 {
			return float64(atomic.LoadUint64(&counter))
		},
	)
	prometheus.MustRegister(counterGauge)

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.POST("/increment", func(c *gin.Context) {
		newVal := atomic.AddUint64(&counter, 1)
		c.JSON(http.StatusOK, gin.H{"value": newVal})
	})

	r.GET("/value", func(c *gin.Context) {
		val := atomic.LoadUint64(&counter)
		c.JSON(http.StatusOK, gin.H{"value": val})
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	_ = r.Run(":8080")
}
