package main

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func handleErr(err interface{}, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}

func sendMsg(pushToken string, pushUser string, msg string) {
	pushUrl := "https://api.pushover.net:443/1/messages.json"
	postMsg := url.Values{
		"token":   {pushToken},
		"user":    {pushUser},
		"message": {msg},
		"sound":   {"cosmic"},
	}
	resp, err := http.PostForm(pushUrl, postMsg)
	defer resp.Body.Close()
	handleErr(err, "Error sending push notification")
}

func main() {
	app := cli.NewApp()
	app.Name = "rain-check"
	app.Usage = "so you never leave without an umbrella!"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "wt",
			Usage: "OpenWeatherMap API token",
		},
		cli.StringFlag{
			Name:  "pt",
			Usage: "Pushover API app token",
		},
		cli.StringFlag{
			Name:  "pu",
			Usage: "Pushover API user token",
		},
		cli.StringFlag{
			Name:  "lat",
			Usage: "Latitude for weather status location. Cannot be specified if city flag is provided.",
		},
		cli.StringFlag{
			Name:  "long",
			Usage: "Longitude for weather status location. Cannot be specified if city flag is provided.",
		},
		cli.StringFlag{
			Name:  "city",
			Usage: "City for weather status location. Cannot be specified if latitude/longitude flags are provided",
		},
	}

	app.Action = func(c *cli.Context) error {
		weatherToken := c.String("wt")
		pushAppToken := c.String("pt")
		pushUserToken := c.String("pu")
		if weatherToken == "" || pushAppToken == "" || pushUserToken == "" {
			log.Fatal("Error: Rain-check requires Pushover API tokens and an OpenWeatherAPI token.  Please provide these using the appropriate cli flags.")
		}
		var lat string
		var long string
		var city string
		var url string
		if (c.String("lat") != "") || (c.String("long") != "") {
			if c.String("city") != "" {
				log.Fatal("Error: Cannot provide latitude/longitude values with city value")
			}
			lat = c.String("lat")
			long = c.String("long")
			if lat == "" || long == "" {
				log.Fatal("Error: Latitude and longitude must both be provided when using geographic coordinates!")
			}
			url = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", lat, long, weatherToken)
		} else {
			city = c.String("city")
			url = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, weatherToken)
		}
		currentWeather := 700
		for {
			resp, err := http.Get(url)
			handleErr(err, "Error fetching weather data.  Be sure you have an internet connection and you provided the correct flags")
			body, _ := ioutil.ReadAll(resp.Body)
			parsed, err := gabs.ParseJSON(body)
			handleErr(err, "Error processing OpenWeatherMap API request.  Please be sure your API token is valid")
			weatherCode, _ := parsed.Path("weather").Index(0).Path("id").Data().(int)
			if weatherCode < 700 {
				if currentWeather >= 700 {
					sendMsg(pushAppToken, pushUserToken, "It's coming down out there!")
					currentWeather = weatherCode
				} else {
					currentWeather = weatherCode
				}
			} else {
				if currentWeather < 700 {
					sendMsg(pushAppToken, pushUserToken, "It cleared up")
					currentWeather = weatherCode
				} else {
					currentWeather = weatherCode
				}
			}
			resp.Body.Close()
			time.Sleep(5 * time.Minute)
		}
		return nil
	}
	app.Run(os.Args)
}
