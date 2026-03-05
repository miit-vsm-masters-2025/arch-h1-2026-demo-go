package main

import (
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	counter := uint64(rng.Intn(100001))

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

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			atomic.AddUint64(&counter, uint64(rng.Intn(101)))
		}
	}()

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
