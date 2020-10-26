package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fastjson"
)

// WeatherPack contains weather data
type WeatherPack struct {
	coord struct {
		lon float64
		lat float64
	}
	main struct {
		temp     float64
		pressure float64
	}
	wind struct {
		speed float64
	}
	city string
	date string
}

// request weather from openweathermap
func requestWeather(city string, interval string, bot *tgbotapi.BotAPI) WeatherPack {
	if len(city) < 0 {
		log.Print("Error: empty city name")
		return WeatherPack{}
	}

	if len(interval) < 0 {
		log.Print("Error: wrong interval")
		return WeatherPack{}
	}

	// TODO: use trim! not it
	if city[0] == ' ' {
		city = city[1:len(city)]
	}
	// TODO: make api key outer parameter
	apiKey := SystemConfig.weatherAPIKey
	var weatherURL string

	// parse 1 or 2 or 3 to weather period
	i, _ := strconv.ParseInt(interval, 10, 64)

	switch i {
	case 0:
		weatherURL = "http://api.openweathermap.org/data/2.5/weather?q=" +
			city + "&appid=" + apiKey + "&units=metric"
	case 1:
		fallthrough
	case 3:
		weatherURL = "http://api.openweathermap.org/data/2.5/forecast?q=" +
			strings.TrimRight(city, " ") + "&appid=" + apiKey + "&units=metric"
	}
	log.Printf("Send request to openweathermap \n\t(%s)", weatherURL)

	// http request section
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	log.Print("sending")
	req, err := http.NewRequest(http.MethodGet, weatherURL, nil)
	_check(err)

	log.Print("wait")
	res, err := spaceClient.Do(req)
	_check(err)

	log.Print("read")
	body, err := ioutil.ReadAll(res.Body)
	_check(err)

	// parse section
	weather := WeatherPack{}
	var dat map[string]interface{}
	coeffHg := 0.750062

	switch i {
	case 0:
		weather.main.temp = fastjson.GetFloat64(body, "main", "temp")
		weather.main.pressure = fastjson.GetFloat64(body, "main", "pressure") * coeffHg
		weather.coord.lat = fastjson.GetFloat64(body, "coord", "lat")
		weather.coord.lon = fastjson.GetFloat64(body, "coord", "lon")
		weather.wind.speed = fastjson.GetFloat64(body, "wind", "speed")
		weather.city = fastjson.GetString(body, "name")
		weather.date = "сейчас"
	case 1:
		if err := json.Unmarshal(body, &dat); err != nil {
			panic(err)
		}

		code, _ := strconv.ParseInt(dat["cod"].(string), 10, 32)
		if code != 200 {
			log.Println("Bad request")
			return weather
		}

		list := dat["list"].([]interface{})

		sudList := list[4].(map[string]interface{}) // 12
		weather.date = sudList["dt_txt"].(string)
		weather.city = dat["city"].(map[string]interface{})["name"].(string)
		weather.main.temp = sudList["main"].(map[string]interface{})["temp"].(float64)
		weather.wind.speed = sudList["wind"].(map[string]interface{})["speed"].(float64)
		weather.main.pressure = sudList["main"].(map[string]interface{})["pressure"].(float64) * coeffHg
	case 3:

		if err := json.Unmarshal(body, &dat); err != nil {
			panic(err)
		}

		code, _ := strconv.ParseInt(dat["cod"].(string), 10, 32)
		if code != 200 {
			log.Println("Bad request")
			return weather
		}

		list := dat["list"].([]interface{})
		sudList := list[12].(map[string]interface{}) // 12
		weather.date = sudList["dt_txt"].(string)
		weather.city = dat["city"].(map[string]interface{})["name"].(string)
		weather.main.temp = sudList["main"].(map[string]interface{})["temp"].(float64)
		weather.wind.speed = sudList["wind"].(map[string]interface{})["speed"].(float64)
		weather.main.pressure = sudList["main"].(map[string]interface{})["pressure"].(float64) * coeffHg
	}

	return weather
}

func checkCity(city string) (bool, string) {
	URL := "http://api.openweathermap.org/data/2.5/find?callback=?&q="
	URL = URL + city + "&type=like&sort=population&cnt=30&appid=6afb386f44c3179f6e9706365bc60d99"

	log.Println("Check city :" + URL)

	log.Print("sending")
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	_check(err)

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	log.Print("wait")
	res, err := spaceClient.Do(req)
	_check(err)

	log.Print("read")
	body, err := ioutil.ReadAll(res.Body)
	_check(err)

	if res.StatusCode != http.StatusOK {
		return false, ""
	}

	body = body[2 : len(body)-1]

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	code, _ := strconv.ParseInt(dat["cod"].(string), 10, 32)
	if code != 200 {
		log.Println("Bad request")
	} else {
		log.Println("Good request")
	}

	count := dat["count"].(float64)
	if count == 0 {
		log.Println("Zero cities found")
		return false, ""
	}

	list := dat["list"].([]interface{})
	subList := list[0].(map[string]interface{})
	cityTrueName := subList["name"].(string)
	log.Println("Found city name:" + cityTrueName)

	return true, cityTrueName
}

func weatherMessageBuilder(w *WeatherPack) string {

	if len(w.city) == 0 {
		return "Такого злоебучего города не существует. Пораскинь мозгами бля!!!	"
	}
	return "В засраном городе " + w.city +
		"\nЕбучая температура: " + strconv.FormatFloat(w.main.temp, 'f', 3, 64) + " °C" +
		"\nДавление ебошит как проклятое: " + strconv.FormatFloat(w.main.pressure, 'f', 2, 64) + " мм рт.ст" +
		"\nСкорость гребанного ветра составляет: " + strconv.FormatFloat(w.wind.speed, 'f', 2, 64) + " м/с" +
		"\nПрогноз вежливо предоставлен ботом Хуезер." +
		"\nДанные предоставлены на злоебучее " + w.date
}
