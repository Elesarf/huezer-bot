package main

import (
	"database/sql"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// parse message
func messageParse(message *InternalMessage, bot *tgbotapi.BotAPI) {
	// collect string
	str := "ChatID: "
	str = str + strconv.FormatInt(message.chatID, 10)
	str = str + " " + message.userName + "\n\t" + message.messageText
	log.Print(str)

	// regex find bash key. If found, get bash citate and return
	niceParse, _ := regexp.Compile("/cheer")
	if niceParse.MatchString(message.messageText) {
		message.messageText = GetCitate()
		return
	}

	// some\nmessage -> some message
	message.messageText = strings.Replace(message.messageText, "\n", " ", -1)
	list := strings.Split(message.messageText, " ")
	s := ""
	// exclude all @some substring
	for i := range list {
		tmp := list[i]
		if !strings.Contains(tmp, "@") {
			s = s + list[i] + " "
		}
	}
	message.messageText = s

	str = strings.ToLower(s)

	// open base
	db, err := sql.Open("sqlite3", SystemConfig.dbPath)
	_check(err)
	defer db.Close()

	// get aliases
	if strings.Contains(message.messageText, "ебни чо знаешь по") {
		str = strings.Replace(message.messageText, "ебни чо знаешь по", "", -1)
		if len(str) == 0 {
			message.messageText = "Косячишь, фуфел"
			return
		}

		// ' some string ' -> 'some string'
		str = strings.TrimSpace(str)
		log.Println("Get aliases" + "_" + str + "_")

		// collect aliases
		var aliases []string
		if !strings.Contains(str, "мне") {
			aliases = getAliases(str, db, "")
		} else {
			aliases = getAliases("", db, message.userName)
		}
		aliasesString := "Цени:\n"
		for i := range aliases {
			aliasesString = aliasesString + aliases[i] + "\n"
		}

		message.messageText = aliasesString
		return
	}

	// make aliases
	if strings.Contains(message.messageText, "запомни мудак") {
		str = strings.Replace(message.messageText, "запомни мудак", "", -1)
		strList := strings.Split(str, ":")
		if len(strList) != 2 {
			message.messageText = "Косячишь, бубен."
			return
		}

		done := tryAddUserCity(strList[0], strList[1], db, message.userName)
		if done == true {
			message.messageText = "Заебок. Отныне " + strList[0] + " : " +
				strings.Title(strList[1]) + " на твоей совести."
		} else {
			message.messageText = "Херня какая - то. Я умываю умывальник."
		}
		return
	}

	// get weater section
	// try check weather, and add new city in base
	if !strings.Contains(message.messageText, "ебни погодку в") {
		message.messageText = ""
		return
	}

	str = strings.Replace(message.messageText, "ебни погодку в", "", -1)
	str = strings.TrimSpace(str)

	if checkKey(str, db) == true {
		w := requestWeather(getValue(str, db), "0", bot)
		message.messageText = weatherMessageBuilder(&w)
		return
	}

	log.Println("Not found city in base. Check on openweather")
	exist, cityName := checkCity(str)

	if exist != true {
		log.Println("Openweather check false.")
		message.messageText = "А вот твой злоебучий город нихуя #ненайденбля"
		return
	}

	w := requestWeather(cityName, "0", bot)
	message.messageText = weatherMessageBuilder(&w)

	if len(w.city) != 0 {
		insertIntoDB(strings.Title(str), w.city, db, message.userName)
		message.messageText = message.messageText +
			"\nЗаебок, я еще твою помойку запомнил.\n" +
			strings.Title(str) + " : " + w.city
	}
}

// try add user city to db
func tryAddUserCity(alias string, city string, db *sql.DB, user string) bool {
	// defaut value. fuck go
	if user == "" {
		user = "system"
	}
	city = strings.ToLower(city)
	city = strings.Title(city)
	city = strings.TrimSpace(city)
	alias = strings.ToLower(alias)
	alias = strings.Title(alias)
	alias = strings.TrimSpace(alias)

	// shit style. Delete later
	if city[0] == ' ' {
		city = city[1:len(city)]
	}
	if alias[0] == ' ' {
		alias = alias[1:len(alias)]
	}

	if checkKey(alias, db) == true {
		log.Println("Error add user city: city is already exist")
		return false
	}

	exist, cityName := checkCity(city)
	if exist == false {
		log.Println("Error add user city: city not found on openweather.com")
		return false
	}

	err := insertIntoDB(alias, cityName, db, user)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

//IsLetter ...
func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
