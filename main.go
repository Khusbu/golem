package main

import (
	"encoding/json"
	"fmt"
	"golem/movies"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
  "strings"
	"github.com/nlopes/slack"
  "golem/zomato"
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
				messageallRegex, _ := regexp.Compile("message all (.*?)")
				pugbombRegex, _ := regexp.Compile("pug bomb (.*?)")
				movieshowingRegex, _ := regexp.Compile("movies showing in (.*?)")
				moviecomingsoonRegex, _ := regexp.Compile("movies coming soon in (.*?)")
				movieimdbRegex, _ := regexp.Compile("movie imdb (.*?)")
				if messageallRegex.MatchString(a.Msg.Text) {
					loc := messageallRegex.FindStringIndex(a.Msg.Text)
					if loc[0] == 0 {
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
							fmt.Println(user.RealName)
							channelID, timestamp, err := api.PostMessage(im.BaseChannel.Id, a.Msg.Text[loc[1]:len(a.Msg.Text)]+" "+user.RealName, params)
							if err != nil {
								fmt.Printf("%s\n", err)
								return
							}
							fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
						}
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
				if pugbombRegex.MatchString(a.Msg.Text) {
					params := slack.PostMessageParameters{
						UnfurlLinks: true,
						UnfurlMedia: true,
					}
					loc := pugbombRegex.FindStringIndex(a.Msg.Text)
					if loc[0] == 0 {
						bomb, err := strconv.Atoi(a.Msg.Text[loc[1]:len(a.Msg.Text)])
						if err != nil {
							channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", "You have to bomb pugs with numbers (1 to 10)", params)
							if err != nil {
								fmt.Printf("%s\n", err)
								return
							}
							fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
						}
						if bomb > 10 {
							channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", "Guys bomb pugs less than or equal to 10", params)
							if err != nil {
								fmt.Printf("%s\n", err)
								return
							}
							fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
						} else {
							for i := 1; i <= bomb; i++ {
								channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", getPugURL("http://pugme.herokuapp.com/random"), params)
								if err != nil {
									fmt.Printf("%s\n", err)
									return
								}
								fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
							}
						}
					}
				}
				if movieshowingRegex.MatchString(a.Msg.Text) {
					loc := movieshowingRegex.FindStringIndex(a.Msg.Text)
					if loc[0] == 0 {
						var postMsg string
						movies := movies.GetNowShowing(a.Msg.Text[loc[1]:len(a.Msg.Text)])

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
				}

				if moviecomingsoonRegex.MatchString(a.Msg.Text) {
					loc := moviecomingsoonRegex.FindStringIndex(a.Msg.Text)
					if loc[0] == 0 {
						var postMsg string
						movies := movies.GetComingSoon(a.Msg.Text[loc[1]:len(a.Msg.Text)])

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
				}

				if movieimdbRegex.MatchString(a.Msg.Text) {
					loc := movieimdbRegex.FindStringIndex(a.Msg.Text)
					if loc[0] == 0 {
						rating := movies.GetIMDBMovieRating(a.Msg.Text[loc[1]:len(a.Msg.Text)])

						channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", fmt.Sprintf("IMDB Rating: %s\nPlot: %s", rating.IMDBRating, rating.Plot), slack.PostMessageParameters{})
						if err != nil {
							fmt.Printf("%s\n", err)
							return
						}
						fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
					}
				}
        if strings.HasPrefix(a.Msg.Text, "zomato") {
          cmd := strings.Fields(a.Msg.Text)
          city, query := cmd[1], cmd[2:]
          cityID := zomato.GetCityID(city)
          restaurants := zomato.GetRestaurants(cityID, "city", strings.Join(query, " "))
          var postMsg string
          for _, restaurant := range restaurants {
            postMsg = fmt.Sprintf("%sName: %s(%s)\nAddress: %s\nRating: %s (%s)\tVotes: %s\n\n", postMsg, restaurant.Restaurant.Name, restaurant.Restaurant.Cuisines, restaurant.Restaurant.Location.Address, restaurant.Restaurant.UserRating.AggregateRating, restaurant.Restaurant.UserRating.RatingText, restaurant.Restaurant.UserRating.Votes)
          }
          channelID, timestamp, err := api.PostMessage("C1Y7PBU9X", postMsg, slack.PostMessageParameters{})
          if err != nil {
            fmt.Printf("%s\n", err)
            return
          }
          fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
        }

				if a.Msg.Text == "help" {
					helpDoc := `help doc:
                  message all <text> - Sends '<text>' to all users
                  pug me - Receive a pug
                  pug bomb N - Receive N pugs
                  movies showing in <city> - Movies showing in <city>
                  movies coming soon in <city> - Movies coming soon in <city>
                  movie imdb <movie name> - Shows IMDB rating for <movie name>
                  zomato <city> <restaurant-name>`

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
