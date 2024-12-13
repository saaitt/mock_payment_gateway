package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"strconv"
	"time"
)
import "net/http"

type Config struct {
	SystemUrl string `mapstructure:"systemUrl"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
}

func main() {
	config := ReadConfig()
	r := gin.Default()
	r.Use(cors.Default())
	r.LoadHTMLFiles("public/index.html")
	r.GET("/", func(c *gin.Context) {
		paymentId, ok := c.GetQuery("payment_id")
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "invalid_request",
				"message": "payment_id was not provided",
			})
			return
		}
		amountS, ok := c.GetQuery("amount")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "invalid_request",
				"message": "amount was not provided",
			})
			return
		}
		amount, err := strconv.Atoi(amountS)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "amount was not an acceptable integer",
			})
		}
		callbackUrl, ok := c.GetQuery("callback")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "nil",
				"message": "callback url was not provided",
			})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"amount":    amount,
			"paymentId": paymentId,
			"host":      callbackUrl,
		})
	})
	r.GET("/payment", func(c *gin.Context) {
		amount, ok := c.GetQuery("amount")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "invalid_request",
				"message": "amount was not provided",
			})
			return
		}
		callbackUrl, ok := c.GetQuery("callback_url")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "invalid_request",
				"message": "callback_url was not provided",
			})
			return
		}
		paymentId := rand.Int()
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"id":     paymentId,
			"url":    fmt.Sprintf("%s/?payment_id=%v&amount=%s&callback=%s", config.SystemUrl, paymentId, amount, callbackUrl),
		})
	})

	r.GET("/callback", func(c *gin.Context) {
		status, ok := c.GetQuery("status")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "invalid_request",
				"message": "status was not provided",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	})

	r.POST("/payment", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%v", config.Host, config.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println(
		fmt.Sprintf("%+v", config))
	s.ListenAndServe()
}

func ReadConfig() Config {
	v := viper.New()
	v.SetConfigType("yaml")
	f, err := os.Open("config.yaml")
	if err != nil {
		panic("cannot read config: " + err.Error())
	}
	err = v.ReadConfig(f)
	if err != nil {
		panic("cannot read config" + err.Error())
	}
	var configs Config
	if err := v.Unmarshal(&configs); err != nil {
		fmt.Println(err)
		return configs
	}
	return configs
}
