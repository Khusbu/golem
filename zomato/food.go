package zomato

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "regexp"
    "net/url"
    "strconv"
)

const (
  locationExp = "{\"entity_type\":\"city\",(.*?),\"country_name\":\"(.*?)\"}"
  restaurantsExp = "{\"restaurant\":{(.*?)}}"
  zomatoToken = "ZOMATO TOKEN HERE"
)

type RestaurantDetails struct {
	  Restaurant struct {
    Name string `json:"name"`
    Location struct  {
      Address string `json:"address"`
    } `json:"location"`
    Cuisines string `json:"cuisines"`
    UserRating struct {
      AggregateRating string `json:"aggregate_rating"`
      RatingText string  `json:"rating_text"`
      Votes string `json:"votes"`
    } `json:"user_rating"`
  } `json:"restaurant"`
}

func GetRestaurants(id int, entityType, q string) []RestaurantDetails {
  query := make(url.Values)
  query.Add("entity_id", strconv.Itoa(id))
  query.Add("entity_type", entityType)
  query.Add("q", q)

	url := &url.URL{RawQuery: query.Encode(), Host: "developers.zomato.com", Path: "/api/v2.1/search", Scheme: "https"}

  client := &http.Client{}
  req, err := http.NewRequest("GET", url.String(), nil)
  if err != nil {
      fmt.Println(err)
  }

  //Zomato API Credentials
  req.Header.Add("user-key", zomatoToken)

  resp, err := client.Do(req)
  if err != nil {
      fmt.Println(err)
  }

  bytes, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
      fmt.Println("Read response error: ", err)
  }

  restaurantsRegex := regexp.MustCompile(restaurantsExp)
  restaurants := restaurantsRegex.FindAll(bytes, -1)

  var restoArray []RestaurantDetails
  for _, restaurant := range restaurants {
      var restaurantData RestaurantDetails
      err := json.Unmarshal(restaurant, &restaurantData)
      if err != nil {
          fmt.Println(err)
      }
      restoArray = append(restoArray, restaurantData)
  }
  return restoArray
}

func GetCityID(city string) int {
    client := &http.Client{}
    req, err := http.NewRequest("GET", fmt.Sprintf("https://developers.zomato.com/api/v2.1/locations?query=%s", city), nil)
    if err != nil {
        fmt.Println(err)
    }

    //Zomato API Credentials
    req.Header.Add("user-key", zomatoToken)

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
    }

    bytes, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        fmt.Println("Read response error: ", err)
    }

    locationRegex := regexp.MustCompile(locationExp)
    details := locationRegex.Find(bytes)

    location := make(map[string]interface{})
    err = json.Unmarshal(details, &location)
    if err != nil {
        fmt.Println("Json unmarshal error: ", err)
    }
    return int(location["city_id"].(float64))

}
