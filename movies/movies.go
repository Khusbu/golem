package movies

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
)

const movieExpression = "{\"name\":\"(.*?)\",(.*​?),\"dimension13\":\"(.*?)\"}"

//Movie structure for storing movie details
type Movie struct {
  Name string `json:"name"`
  Language string `json:"dimension13"`
}

//MovieIMDB structure for storing movie IMDB rating
type MovieIMDB struct {
  IMDBRating string
  Plot string
}

func getMovies(movies [][]byte) []Movie {
    var (
      movieArray []Movie
      detailMap Movie
    )
    for _, movie := range movies {
        movieDetailRegex := regexp.MustCompile(movieExpression)
        details := movieDetailRegex.Find(movie)

        err := json.Unmarshal(details, &detailMap)
        if err != nil {
            fmt.Println(err)
        }
        movieArray = append(movieArray, detailMap)
    }
    return movieArray
}

//GetComingSoon coming soon
func GetComingSoon(location string) []Movie {
    comingSoonRegex := regexp.MustCompile("(.*?)\"list\":\"category(.*?)coming soon\"}(.*?)}}}}")

    resp, err := http.Get(fmt.Sprintf("https://in.bookmyshow.com/%s/movies", location))
    response, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        fmt.Println(err)
    }
    movies := comingSoonRegex.FindAll(response, -1)
    return getMovies(movies)
}

//GetNowShowing now showing
func GetNowShowing(location string) []Movie {
    nowShowingRegex := regexp.MustCompile("(.*?)\"list\":\"Filter Impression:category(.*?)now showing\"}(.*?)}]}}}")

    resp, err := http.Get(fmt.Sprintf("https://in.bookmyshow.com/%s/movies", location))
    response, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        fmt.Println(err)
    }
    movies := nowShowingRegex.FindAll(response, -1)
    return getMovies(movies)
}

//GetIMDBMovieRating imdb movie rating
func GetIMDBMovieRating(movie string) MovieIMDB {
    resp, err := http.Get(fmt.Sprintf("http://www.omdbapi.com/?t=%s&y=&plot=full&r=json", movie))
    response, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        fmt.Println(err)
    }
    var ratingMap MovieIMDB
    err = json.Unmarshal(response, &ratingMap)
    if err != nil {
        fmt.Println(err)
    }
    return ratingMap
}
