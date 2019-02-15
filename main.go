package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

//Bot struct  hold config
type Bot struct {
	Name        string
	Idtoken     string
	Accesstoken string
	Apibaseurl  string
}

type Msgitem struct {
	Items []struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		TargetURL string    `json:"targetUrl"`
		Resource  string    `json:"resource"`
		Event     string    `json:"event"`
		OrgID     string    `json:"orgId"`
		CreatedBy string    `json:"createdBy"`
		AppID     string    `json:"appId"`
		OwnedBy   string    `json:"ownedBy"`
		Status    string    `json:"status"`
		Created   time.Time `json:"created"`
	} `json:"items"`
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

	r.POST("/printrequest", func(c *gin.Context) {
		requestDump, _ := httputil.DumpRequest(c.Request, true)
		fmt.Println(string(requestDump))
		c.String(http.StatusOK, string(requestDump))
	})

	//send message with post json body like {"roomid":"Y2lzY29zcGFyazovL3VzL1JPT00vMWFhZjg2NWUtNTc4Ni0zYjRmLWEwNGQtMTg5Yzg0MTIyYmIx","text":"Sometext" }
	r.POST("/bot/sendmessage", func(c *gin.Context) {

	})

	r.POST("/bot/createwebhook", func(c *gin.Context) {
		//Todo
	})

	r.POST("/bot/webhook", func(c *gin.Context) {
		var m Msgitem
		if c.ShouldBind(&m) == nil {
			fmt.Println(m.Items[0])
		}
		c.String(200, "get message")

	})

	r.GET("/bot/getmsg/msgid", func(c *gin.Context) {
		//https://api.ciscospark.com/v1/messages/Y2lzY29zcGFyazovL3VzL01FU1NBR0UvYWM0YTFkZDAtMzBmNS0xMWU5LThiMDYtOWRhZjExZjViNWMy
		_, body, _ := gorequest.New().Get(b.Apibaseurl+"messages/").Set("Authorization", "Bearer "+b.Accesstoken).End()
		//_, body, _ := gorequest.New().Set("Authorization", "Bearer "+b.Accesstoken).Get("http://127.0.0.1:8080/printrequest").End()

		fmt.Println(body)
		c.String(200, body)

	})

	r.GET("/bot/listwebhook", func(c *gin.Context) {
		//https://api.ciscospark.com/v1/webhooks
		_, body, _ := gorequest.New().Get(b.Apibaseurl+"webhooks").Set("Authorization", "Bearer "+b.Accesstoken).End()
		//_, body, _ := gorequest.New().Set("Authorization", "Bearer "+b.Accesstoken).Get("http://127.0.0.1:8080/printrequest").End()

		fmt.Println(body)
		c.String(200, body)

	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
