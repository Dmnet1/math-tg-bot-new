package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
	"os"
	"strings"
)

type Config struct {
	TelegramBotToken string
}

type Equations struct {
	XMLName   xml.Name   `xml:"compendium"`
	Equations []Equation `xml:"equation"`
	//Equation  Equation   `xml:"equation"`
}

type Equation struct {
	XMLName     xml.Name `xml:"equation"`
	Name        string   `xml:"name"`
	Odds        string   `xml:"odds"`
	GeneralForm string   `xml:"generalForm"`
	SupportInfo []string `xml:"supportInfo"`
}

func Filter(equations []Equation, fn func(equation Equation) bool) []Equation {
	var filtered []Equation
	for _, equation := range equations {
		if fn(equation) {
			filtered = append(filtered, equation)
		}
	}
	return filtered
}

func parseEquations() (Equations, error) {
	file, err := os.Open("types.xml")
	if err != nil {
		log.Panic(err)
	} else {
		log.Println(file)
	}

	fi, err := file.Stat()
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	var data = make([]byte, fi.Size())
	_, err = file.Read(data)
	if err != nil {
		log.Panic(err)
	}

	var v Equations
	err = xml.Unmarshal(data, &v)

	if err != nil {
		log.Println(err)
		return v, err
	}
	return v, nil
}

func main() {
	file, err := os.Open("tsconfig.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(configuration.TelegramBotToken)

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	equations, err := parseEquations()
	if err != nil {
		log.Panic(err)
	}

	// В канал updates будут приходить все новые сообщения.
	for update := range updates {
		query := update.Message.Text
		filtered := Filter(equations.Equations, func(equation Equation) bool {
			return strings.Index(strings.ToLower(equation.Name), strings.ToLower(query)) >= 0
		})

		if len(filtered) == 0 && query != "/start" && query != "уравнения" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Попробуй еще разок!")
			bot.Send(msg)
		}

		for _, equation := range filtered {
			text := ""
			for _, t := range equation.SupportInfo {
				text = text + t + "\n"
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n%s", equation.Name, text))
			bot.Send(msg)
		}

		if query == "уравнения" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n", equations.Equations))
			bot.Send(msg)
		}

		if query == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Это математический бот для решения разного вида уравнений. Для поиска интересующего вида уравнения отправь в чат: 'уравнение'. Все, что содержит данный запрос, будет выведено новым сообщением")
			bot.Send(msg)
		}
	}
}

//новая ветка
