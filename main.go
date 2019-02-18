package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

//Msgitem is the what webhook receive message
type Msgitem struct {
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
	ActorID   string    `json:"actorId"`
	Data      struct {
		ID          string    `json:"id"`
		RoomID      string    `json:"roomId"`
		RoomType    string    `json:"roomType"`
		PersonID    string    `json:"personId"`
		PersonEmail string    `json:"personEmail"`
		Created     time.Time `json:"created"`
	} `json:"data"`
}

type Msgcontent struct {
	ID          string    `json:"id"`
	RoomID      string    `json:"roomId"`
	RoomType    string    `json:"roomType"`
	Text        string    `json:"text"`
	PersonID    string    `json:"personId"`
	PersonEmail string    `json:"personEmail"`
	Created     time.Time `json:"created"`
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
	//fmt.Println(b.Accesstoken)
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
		fmt.Println(string(requestDump))
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
		body, _ := ioutil.ReadAll(c.Request.Body)
		errs := json.Unmarshal(body, &m)
		if errs != nil {
			fmt.Println("unmarshal webhook text fail")
		}
		msg := getmessage(b, m.Data.ID)
		fmt.Println("Got spark message : " + msg)
		c.String(200, msg)

		// Q? is c ShouldBindJSON not ok?
		// if c.ShouldBindJSON(&m) == nil {
		// 		fmt.Println(m.ActorID)
		// 		fmt.Println(m.Data.ID)
		// 		//Y2lzY29zcGFyazovL3VzL01FU1NBR0UvNzFkODViYzAtMzBmYi0xMWU5LWFiNmEtNzVjODMzNjViYTFm
		// 		msg := getmessage(b, m.Data.ID)
		// 		fmt.Println(m)
		// 		fmt.Println(&m)
		// 		c.String(200, msg)
		// }

	})

	// r.GET("/bot/getmsg/msgid", func(c *gin.Context) {
	// })

	r.GET("/bot/listwebhook", func(c *gin.Context) {
		//https://api.ciscospark.com/v1/webhooks
		_, body, _ := gorequest.New().Get(b.Apibaseurl+"webhooks").Set("Authorization", "Bearer "+b.Accesstoken).End()
		//_, body, _ := gorequest.New().Set("Authorization", "Bearer "+b.Accesstoken).Get("http://127.0.0.1:8080/printrequest").End()

		fmt.Println(body)
		c.String(200, body)

	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func getmessage(b *Bot, msgid string) string {
	//https://api.ciscospark.com/v1/messages/Y2lzY29zcGFyazovL3VzL01FU1NBR0UvYWM0YTFkZDAtMzBmNS0xMWU5LThiMDYtOWRhZjExZjViNWMy
	//fmt.Println("Authorization: Bearer " + b.Accesstoken)
	//fmt.Println(b.Apibaseurl + "messages/" + msgid)
	var content Msgcontent
	_, bodyBytes, err := gorequest.New().Get(b.Apibaseurl+"messages/"+msgid).Set("Authorization", "Bearer "+b.Accesstoken).EndBytes()
	if err != nil {
		fmt.Println("call to get msg content fail")
	}

	errs := json.Unmarshal(bodyBytes, &content)
	if errs != nil {
		fmt.Println("unmarshal msg text fail")
	}
	//fmt.Println(content)
	return content.Text

}
