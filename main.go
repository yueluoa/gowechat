package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var bot *openwechat.Bot

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/open-wechat/send-message", sendMessage)

	return router
}

func main() {
	bot = openwechat.DefaultBot(openwechat.Desktop)

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			_, _ = msg.ReplyText("pong")
		}
	}

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	router := setupRouter()
	err := router.Run(":8122")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	_ = bot.Block()
}

type Message struct {
	Contents []string `json:"contents" form:"contents"`
}

func sendMessage(ctx *gin.Context) {
	req := &Message{}
	if err := ctx.ShouldBind(req); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("请求参数: %v\n", req)

	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("获取登陆的用户: %v\n", self)

	// 获取所有的好友
	friends, err := self.Friends()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	friend := friends.SearchByNickName(1, "7喜")

	text := strings.Join(req.Contents, "\n")
	if err := friend.SendText(text, 1); err != nil {
		return
	}

	ctx.JSON(http.StatusOK, struct{}{})
}
