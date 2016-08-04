package main

import (
    "fmt"
    // "time"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "github.com/nlopes/slack"
)

func getPugURL(url string) string {
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println("Error getting response: ", err)
  }
  bytes, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    fmt.Println("Read response error: ", err)
  }
  puggy := make(map[string]string)
  err = json.Unmarshal(bytes, &puggy)
  if err != nil {
    fmt.Println("Json unmarshal error: ", err)
  }
  return puggy["pug"]
}

func main() {
    // chSender := make(chan slack.OutgoingMessage)
    chReceiver := make(chan slack.SlackEvent)

    api := slack.New("xoxb-66220411602-wZgsZhcITvQ3YEJMxFEn69RG")
    // api.SetDebug(true)
    wsAPI, err := api.StartRTM("", "http://example.com")
    if err != nil {
        fmt.Printf("%s\n", err)
    }
    go wsAPI.HandleIncomingEvents(chReceiver)
    // go wsAPI.Keepalive(20 * time.Second)
    // go func(wsAPI *slack.SlackWS, chSender chan slack.OutgoingMessage) {
    //     for {
    //         select {
    //         case msg := <-chSender:
    //             wsAPI.SendMessage(&msg)
    //         }
    //     }
    // }(wsAPI, chSender)
    for {
        select {
        case msg := <-chReceiver:
            fmt.Print("Event Received: ")
            switch msg.Data.(type) {
            case *slack.MessageEvent:
                a := msg.Data.(*slack.MessageEvent)
                fmt.Printf("Message: %s\n", a.Msg.Text)

                if a.Msg.Text == "message all" {
                  params := slack.PostMessageParameters{}
                  imChannels, err := api.GetIMChannels()
                  if err != nil {
                    fmt.Println(err)
                  }
                  for _, im := range imChannels {
                    user, err := api.GetUserInfo(im.UserId)
                    if err != nil {
                      fmt.Println(err)
                    }

                    channelID, timestamp, err := api.PostMessage(im.BaseChannel.Id, "Hi " + user.RealName,  params)
                    if err != nil {
                  		fmt.Printf("%s\n", err)
                  		return
                  	}
    	              fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
                  }
                }
                if a.Msg.Text == "pug me" {
                  params := slack.PostMessageParameters{
                    UnfurlLinks: true,
                    UnfurlMedia: true,
                  }
                  channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", getPugURL("http://pugme.herokuapp.com/random"), params)
                	if err != nil {
                		fmt.Printf("%s\n", err)
                		return
                	}
  	              fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
                }
            default:
                //fmt.Printf("Unexpected: %v\n", msg.Data)
            }
        }
    }
}
