package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "github.com/nlopes/slack"
    "golem/movies"
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
    chReceiver := make(chan slack.SlackEvent)

    api := slack.New("YOUR TOKEN HERE")
    wsAPI, err := api.StartRTM("", "http://example.com")
    if err != nil {
        fmt.Printf("%s\n", err)
    }
    go wsAPI.HandleIncomingEvents(chReceiver)
    for {
        select {
        case msg := <-chReceiver:
            fmt.Print("Event Received: ")
            switch msg.Data.(type) {
            case *slack.MessageEvent:
                a := msg.Data.(*slack.MessageEvent)
                fmt.Printf("Message: %s\n", string(a.Msg.Text))

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
                if a.Msg.Text == "pug bomb" {
        					params := slack.PostMessageParameters{
        						UnfurlLinks: true,
        						UnfurlMedia: true,
        					}
        					for i := 1; i <= 2; i++ {
        						channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", getPugURL("http://pugme.herokuapp.com/random"), params)
        						if err != nil {
        							fmt.Printf("%s\n", err)
        							return
        						}
        						fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
        					}
      				  }
                if a.Msg.Text == "movies now showing" {
                  var postMsg string
                  movies := movies.GetNowShowing("kolkata")

                  for _, movie := range movies {
                    postMsg = fmt.Sprintf("%s %s(%s)\n", postMsg, movie.Name, movie.Language)
                  }
                  channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", postMsg, slack.PostMessageParameters{})
                  if err != nil {
                    fmt.Printf("%s\n", err)
                    return
                  }
                  fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
                }

                if a.Msg.Text == "movies coming soon" {
                  var postMsg string
                  movies := movies.GetComingSoon("kolkata")

                  for _, movie := range movies {
                    postMsg = fmt.Sprintf("%s %s(%s)\n", postMsg, movie.Name, movie.Language)
                  }
                  channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", postMsg, slack.PostMessageParameters{})
                  if err != nil {
                    fmt.Printf("%s\n", err)
                    return
                  }
                  fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
                }

                if a.Msg.Text == "movie imdb" {
                  rating := movies.GetIMDBMovieRating("sultan")

                  channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", fmt.Sprintf("IMDB Rating: %s\nPlot: %s", rating.IMDBRating, rating.Plot), slack.PostMessageParameters{})
                  if err != nil {
                    fmt.Printf("%s\n", err)
                    return
                  }
                  fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
                }

                if a.Msg.Text == "help" {
                  helpDoc := `help doc:
                  message all - Sends 'Hi <name>' to all users
                  pug me - Receive a pug
                  pug bomb - Receive 2 pugs
                  movies now showing - Movies showing now in Kolkata
                  movies coming soon - Movies coming soon in Kolkata
                  movie imdb - Shows IMDB rating for Sultan`


                  channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", helpDoc, slack.PostMessageParameters{})
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
