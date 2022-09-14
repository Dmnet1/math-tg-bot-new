package main

import (
	"encoding/json"
	"encoding/xml"
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

		if len(filtered) == 0 && query != "/start" && query != "уравнения" && query != "a*x*x + b*x + c = 0" {
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

			text = "Отправь сообщение ниже боту"
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n", equation.GeneralForm))
			bot.Send(msg)
		}

		if query == "уравнения" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n", equations.Equations))
			bot.Send(msg)
		}

		if query == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Это математический бот для решения разного вида уравнений. Для поиска интересующего вида уравнения отправь в чат: 'уравнение'. Все, что содержит данный запрос, будет выведено новым сообщением.")
			bot.Send(msg)
		}

		var queryA, queryB, queryC, D float64
		var notice string

		// Логика получения дискриминанта и корней квадратного уравнения.
		if query == "a*x*x + b*x + c = 0" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите коэффициент 'a'")
			bot.Send(msg)
			notice = "GoToB"

			for update := range updates {
				query := update.Message.Text

				if query != "" && notice == "GoToB" {
					query := update.Message.Text
					queryA, _ = strconv.ParseFloat(query, 64)

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите коэффициент 'b'")
					bot.Send(msg)
					notice = "GoToC"
				}

				for update := range updates {
					query := update.Message.Text

					if query != "" && notice == "GoToC" {
						query := update.Message.Text
						queryB, _ = strconv.ParseFloat(query, 64)
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите коэффициент 'c'")
						bot.Send(msg)
						notice = "GoToD"
					}

					for update := range updates {
						query := update.Message.Text

						if query != "" && notice == "GoToD" {
							query := update.Message.Text
							queryC, _ = strconv.ParseFloat(query, 64) // проверить ошибку

							if queryA != 0 && queryB != 0 && queryC != 0 {
								D, notice, err = discriminant.Discriminant(queryA, queryB, queryC)
							}

							if queryA != 0 && queryB != 0 && queryC != 0 {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Дискриминант равен %f\n", D))
								bot.Send(msg)
								msg = tgbotapi.NewMessage(update.Message.Chat.ID, notice)
								bot.Send(msg)
							}
							notice = "end"

							if D > 0 && queryA != 0 && queryB != 0 && queryC != 0 {
								x1, x2, _ := discriminant.X1X2(queryA, queryB, D) // проверить ошибку
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Корни уравнения: x1 = %f\n, x2 = %f\n", x1, x2))
								bot.Send(msg)
							}

							if D == 0 && queryA != 0 && queryB != 0 && queryC != 0 {
								x, _ := discriminant.X(queryA, queryB) // проверить ошибку
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Корень уравнения: x = %f\n", x))
								bot.Send(msg)
							}

							if queryA != 0 && queryB == 0 && queryC == 0 {
								x := 0
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Корень уравнения: x = %v\n", x))
								bot.Send(msg)
							}

							if queryA != 0 && queryB == 0 && queryC != 0 {
								x1, x2, _ := discriminant.SpecialCase1(queryA, queryC)
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Корни уравнения:\n x1 = √%f\n x2 = -√%f\n", x1, x2)) // получается корень из отрицательного числа
								bot.Send(msg)
							}

							if queryA != 0 && queryB != 0 && queryC == 0 {
								x1, x2, _ := discriminant.SpecialCase2(queryA, queryB)
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Корни уравнения:\n x1 = %f\n x2 = %f\n", x1, x2))
								bot.Send(msg)
							}

							if queryA == 0 {
								text := "При нулевом коэффициенте 'а' уравнение становится линейным. Воспользуйся командой 'a*x + b = 0' для решения данного типа уравнений."
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
								bot.Send(msg)
							}
						}
						break
					}
					break
				}
				break
			}
		}
	}
}
