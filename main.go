package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

//Bot struct  hold config
type Bot struct {
	Name        string
	Idtoken     string
	Accesstoken string
}

func initbot() *Bot {
	var bot Bot
	if _, err := toml.DecodeFile("config.toml", &bot); err != nil {
		log.Fatal(err)
	}
	return &bot
}

func main() {
	b := initbot()
	fmt.Println(b.Accesstoken)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "damon's cisco spark bot ")
	})
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OKOKOK",
		})
	})

	r.GET("/printrequest", func(c *gin.Context) {
		requestDump, _ := httputil.DumpRequest(c.Request, true)
		fmt.Println(requestDump)
		c.String(http.StatusOK, string(requestDump))
	})

	//send message with post json body like {"roomid":"Y2lzY29zcGFyazovL3VzL1JPT00vMWFhZjg2NWUtNTc4Ni0zYjRmLWEwNGQtMTg5Yzg0MTIyYmIx","text":"Sometext" }
	r.POST("/bot/sendmessage", func(c *gin.Context) {

	})

	//send message with post json body like {"name": "receivemessage","targetUrl": "http://173.39.230.83:18180/printrequest","resource": "messages","event": "created"}
	r.POST("/bot/createwebhook", func(c *gin.Context) {
		//https://api.ciscospark.com/v1/messages

	})

	r.GET("/bot/listwebhook", func(c *gin.Context) {
		//https://api.ciscospark.com/v1/webhooks
		_, body, _ := gorequest.New().Get("https://api.ciscospark.com/v1/webhooks").Set("Authorization", "Bearer "+b.Accesstoken).End()
		//_, body, _ := gorequest.New().Set("Authorization", "Bearer "+b.Accesstoken).Get("http://127.0.0.1:8080/printrequest").End()

		fmt.Println(body)
		c.String(200, body)

	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
