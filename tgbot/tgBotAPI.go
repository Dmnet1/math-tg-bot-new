package tgbot

import (
	"encoding/json"
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
	"math-tg-bot-new/discriminant"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	TelegramBotToken string
}

type Equation struct {
	Name        string
	Odds        string
	GeneralForm string
	SupportInfo string
}

var queryA, queryB, queryC, D float64
var notice string
var update tgbotapi.Update
var bot *tgbotapi.BotAPI

func inputCoefficientA(updates tgbotapi.UpdatesChannel, notice string) (float64, string) {
	var err error
	for update := range updates {
		query := update.Message.Text

		if query != "" && notice == "GoToB" {
			query := update.Message.Text
			queryA, err = strconv.ParseFloat(query, 64)
			treatmentErrors(err)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите коэффициент 'b'")
			bot.Send(msg)
			notice = "GoToC"
		}
		inputCoefficientB(updates, "GoToC")
	}
	return queryA, notice
}

func inputCoefficientB(updates tgbotapi.UpdatesChannel, notice string) (float64, string) {
	var err error
	for update := range updates {
		query := update.Message.Text
		if query != "" && notice == "GoToC" {
			query := update.Message.Text
			queryB, err = strconv.ParseFloat(query, 64)
			treatmentErrors(err)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите коэффициент 'c'")
			bot.Send(msg)
			notice = "GoToD"
		}
		inputCoefficientC(updates, "GoToD")
	}
	return queryB, notice
}

func inputCoefficientC(updates tgbotapi.UpdatesChannel, notice string) (float64, string) {
	var err error
	for update := range updates {
		query := update.Message.Text
		if query != "" && notice == "GoToD" {
			query := update.Message.Text
			queryC, err = strconv.ParseFloat(query, 64)
			treatmentErrors(err)
		}
		findingDiscriminant(queryA, queryB, queryC)
		findingRootsOfEquation()
	}
	return queryC, notice
}

func findingDiscriminant(queryA, queryB, queryC float64) {
	if queryA != 0 && queryB != 0 && queryC != 0 {
		D, notice, err := discriminant.Discriminant(queryA, queryB, queryC)

		if err != nil {
			log.Printf("Ошибка при расчете дискриминанта: \n%s", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ошибка при расчете дискриминанта")
			bot.Send(msg)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Дискриминант равен %f\n", D))
		bot.Send(msg)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, notice)
		bot.Send(msg)
	}
}

func findingRootsOfEquation() {
	switch {
	case queryA == 0:
		text := "При нулевом коэффициенте 'а' уравнение становится линейным. Воспользуйся командой 'a*x + b = 0' для решения данного типа уравнений."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		bot.Send(msg)

	case D > 0 && queryA != 0 && queryB != 0 && queryC != 0:
		x1, x2, err := discriminant.X1X2(queryA, queryB, D)

		if err != nil {
			log.Printf("Ошибка при вычислении корней уравнения: \n%s", err)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Корни уравнения1:\n x1 = %f\n x2 = %f\n", x1, x2))
		bot.Send(msg)

	case D == 0 && queryA != 0 && queryB != 0 && queryC != 0:
		x, err := discriminant.X(queryA, queryB)

		if err != nil {
			log.Printf("Ошибка при вычислении корня уравнения: \n%s", err)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Корень уравнения2: x = %f\n", x))
		bot.Send(msg)

	case queryA != 0 && queryB == 0 && queryC == 0:
		x := 0
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Корень уравнения3: x = %v\n", x))
		bot.Send(msg)

	case queryA != 0 && queryB == 0 && queryC != 0:
		x1, x2, err := discriminant.SpecialCase1(queryA, queryC)

		if err != nil {
			log.Printf("Ошибка при вычислении корней уравнения (частный случай 1): \n%s", err)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Корни уравнения4:\n x1 = √%f\n x2 = -√%f\n", x1, x2)) // получается корень из отрицательного числа
		bot.Send(msg)

	case queryA != 0 && queryB != 0 && queryC == 0:
		x1, x2, err := discriminant.SpecialCase2(queryA, queryB)

		if err != nil {
			log.Printf("Ошибка при вычислении корней уравнения (частный случай 2): \n%s", err)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Корни уравнения5:\n x1 = %f\n x2 = %f\n", x1, x2))
		bot.Send(msg)
	}
}

func treatmentErrors(err error) {
	if err != nil {
		log.Printf("Ошибка при конвертации типа: \n%s", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ошибка при вводе коэффициента")
		bot.Send(msg)
	}
}

func TgBotApi() {
	file, err := os.Open("tsconfig.json")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(configuration.TelegramBotToken)

	bot, err = tgbotapi.NewBotAPI(configuration.TelegramBotToken)

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

	equations, err := UnmarshalEquations()
	if err != nil {
		log.Panic(err)
	}

	for update = range updates {
		query := update.Message.Text
		filtered := Filter(equations, func(equation Equation) bool {
			return strings.Index(strings.ToLower(equation.Name), strings.ToLower(query)) >= 0
		})

		isUnknownQuery := len(filtered) == 0 && query != "/start" && query != "уравнения" && query != "a*x*x + b*x + c = 0"
		if isUnknownQuery {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Попробуй еще разок!")
			bot.Send(msg)
			continue
		}

		for _, equation := range filtered {
			text := ""
			for _, t := range equation.SupportInfo {
				text = text + string(t)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n%s", equation.Name, text))
			bot.Send(msg)

			text = "Отправь сообщение ниже боту"
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n", equation.GeneralForm))
			bot.Send(msg)
		}

		if query == "уравнения" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n", Equation{}))
			bot.Send(msg)
		}

		if query == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Это математический бот для решения разного вида уравнений. Для поиска интересующего вида уравнения отправь в чат: 'уравнение'. Все, что содержит данный запрос, будет выведено новым сообщением.")
			bot.Send(msg)
		}

		if query == "a*x*x + b*x + c = 0" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите коэффициент 'a'")
			bot.Send(msg)
			notice = "GoToB"
			inputCoefficientA(updates, "GoToB")
		}
		break
	}
}
