package main

import (
	"log"
	"net/http"

	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

func init() {
	botHandler, err := httphandler.New(
		"Channel Secret",
		"アクセストークン",
	)
	if err != nil {
		log.Fatal(err)
	}
	
	botHandler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		ctx := appengine.NewContext(r)
		bot, err := botHandler.NewClient(linebot.WithHTTPClient(urlfetch.Client(ctx)))
		if err != nil {
			aelog.Errorf(ctx, "%v", err)
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).WithContext(ctx).Do(); err != nil {
						aelog.Errorf(ctx, "%v", err)
					}
				}
			}
		}
	})
	
	http.Handle("/", botHandler)
}
