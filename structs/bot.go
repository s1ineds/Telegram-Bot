/*
	Метод getUpdates.
	Если в строке запроса использовать параметр offset,
	то метод getUpdates вернет только то обновление, идентификатор которого указан в offset.
	Поэтому в строке запроса, нужно указывать offset увеличенный на 1. Он как бы будет ждать нового обновления с
	переданным offset'ом.
	При этом, старые обновления будут утеряны.

*/

package structs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Bot struct {
	offset   int
	url      string
	upd      Update
	db       Database
	userObjs []*Users
	answers  []string
}

func (b *Bot) Go() {
	b.userObjs = make([]*Users, 0)
	b.db = Database{driver: "postgres", connectionString: "user=postgres password=P@ssw0rd! dbname=gobotdb sslmode=disable"}

	// main loop
	for {
		b.url = "https://api.telegram.org/bot" + os.Getenv("TOKEN") + "/getUpdates?offset=" + strconv.Itoa(b.offset) + "&timeout=15"
		fmt.Println("URL: ", b.url)

		client := http.Client{}
		req, err := http.NewRequest(http.MethodGet, b.url, nil)
		if err != nil {
			log.Fatal(err)
		}

		resp, _ := client.Do(req)

		respBody, _ := io.ReadAll(resp.Body)
		un_err := json.Unmarshal([]byte(respBody), &b.upd)
		if un_err != nil {
			fmt.Println(un_err)
		}

		// if Result array is empty then there is no new updates
		if len(b.upd.Result) > 0 {
			fmt.Println(b.upd.Result[0].Message.Text)
			//fmt.Println(b.upd.Result[0].Callback_query.Data)

			if b.upd.Result[0].Message.Text == "/start" {
				b.sendTextMessage(b.upd.Result[0].Message.Chat.Id, "Привет, "+b.upd.Result[0].Message.From.First_name+"!")

				b.db.InsertUser(b.upd.Result[0].Message.Chat.Id, b.upd.Result[0].Message.From.First_name)

				b.userObjs = b.db.GetAllUsers()
				fmt.Println("[b.userObjs] ", b.userObjs)
			} else {
				if len(b.userObjs) < 1 {
					b.userObjs = b.db.GetAllUsers()
					fmt.Println("[b.userObjs] ", b.userObjs)
				}
			}

			for _, obj := range b.userObjs {
				if obj.chatId == b.upd.Result[0].Message.Chat.Id {
					if b.upd.Result[0].Message.Text == "/expenses" || obj.userTypingExpenses {
						if !obj.userTypingExpenses {
							obj.userTypingExpenses = true
							obj.state = "category"
						}
						b.handleExpenses(obj)
					} else if b.upd.Result[0].Message.Text == "/income" || obj.userTypingIncome {
						if !obj.userTypingIncome {
							obj.userTypingIncome = true
							obj.state = "category"
						}
						b.handleIncome(obj)
					} else if b.upd.Result[0].Message.Text == "/report" || obj.userGettingReport {
						if !obj.userGettingReport {
							obj.userGettingReport = true
							obj.state = "typing-start-date"
						}
						b.getReport(obj)
					}
					if b.upd.Result[0].Message.Text == "/help" {
						var helpTextChan chan string = make(chan string)

						go b.getHelp("help.txt", helpTextChan)

						b.sendTextMessage(obj.chatId, string(<-helpTextChan))
					}
				}
			}

			b.offset = b.upd.Result[0].Update_id
			b.offset += 1
		}
	}
}

func (b *Bot) handleIncome(usr *Users) {
	if usr.state == "category" {
		client := http.Client{}

		keyboard_obj := `{"keyboard":[[{"text":"Заработная плата"}],[{"text":"Доп. заработок"}],[{"text":"Прочее"}]]}`

		tmp_url := "https://api.telegram.org/bot" + os.Getenv("TOKEN") + "/sendMessage?chat_id=" +
			strconv.Itoa(usr.chatId) +
			"&text=Выбери категорию из меню с кнопками." +
			"&reply_markup=" + keyboard_obj

		req, _ := http.NewRequest(http.MethodPost, tmp_url, nil)
		_, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		usr.state = "item"
	} else {
		b.dialogue(usr)
	}
}

func (b *Bot) handleExpenses(usr *Users) {
	if usr.state == "category" {
		client := http.Client{}
		keyboard_obj := `{"keyboard":[[{"text":"Продукты питания"},{"text":"Медикаменты"}],[{"text":"Ком. услуги"},{"text":"Транспорт"}],[{"text":"Техника"},{"text":"Образование"}],[{"text":"Развлечения"},{"text":"Подарки"}]]}`

		tmp_url := "https://api.telegram.org/bot" + os.Getenv("TOKEN") + "/sendMessage?chat_id=" +
			strconv.Itoa(usr.chatId) +
			"&text=Выбери категорию из меню с кнопками." +
			"&reply_markup=" + keyboard_obj

		req, _ := http.NewRequest(http.MethodPost, tmp_url, nil)
		_, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		usr.state = "item"
	} else {
		b.dialogue(usr)
	}
}

func (b *Bot) dialogue(usr *Users) {
	if usr.state == "item" {
		b.sendTextMessage(usr.chatId, "Категория: "+b.upd.Result[0].Message.Text)
		b.answers = append(b.answers, b.upd.Result[0].Message.Text)
		b.sendTextMessage(usr.chatId, "Введи свой расход или доход.")
		usr.state = "price"
	} else if usr.state == "price" {
		b.sendTextMessage(usr.chatId, "Покупка: "+b.upd.Result[0].Message.Text)
		b.answers = append(b.answers, b.upd.Result[0].Message.Text)
		b.sendTextMessage(usr.chatId, "Введи цену.")
		usr.state = "date"
	} else if usr.state == "date" {
		b.sendTextMessage(usr.chatId, "Цена: "+b.upd.Result[0].Message.Text)
		b.answers = append(b.answers, b.upd.Result[0].Message.Text)
		b.sendTextMessage(usr.chatId, "Введи дату.")
		usr.state = "check-date"
	} else if usr.state == "check-date" {
		if tmpResult := b.checkDate(usr, b.upd.Result[0].Message.Text); tmpResult {
			b.sendTextMessage(usr.chatId, "Данные, которые ты ввел: Категория: "+b.answers[0]+", Предмет покупки: "+b.answers[1]+", Цена: "+b.answers[2]+", Дата покупки: "+b.answers[3])
			prc, _ := strconv.Atoi(b.answers[2])
			if usr.userTypingExpenses {
				b.db.InsertEntry("expenses", b.answers[0], b.answers[1], prc, b.answers[3], usr.dbId)
			} else if usr.userTypingIncome {
				b.db.InsertEntry("income", b.answers[0], b.answers[1], prc, b.answers[3], usr.dbId)
			}
			usr.state = "reset"
			b.resetUser(usr)
		}
	}
}

func (b *Bot) getReport(usr *Users) {
	if usr.state == "typing-start-date" {
		b.sendTextMessage(usr.chatId, "Сейчас тебе нужно будет ввести период даты, за которую ты хочешь получить отчет.")
		b.sendTextMessage(usr.chatId, "Введи начало периода в формате ГГГГ-мм-дд")
		usr.state = "check-start-date"
		return
	}
	if usr.state == "check-start-date" && b.checkDate(usr, b.upd.Result[0].Message.Text) {
		b.sendTextMessage(usr.chatId, "Теперь введи конец периода в формате ГГГГ-мм-дд")
		usr.state = "check-end-date"
		return
	}
	if usr.state == "check-end-date" && b.checkDate(usr, b.upd.Result[0].Message.Text) {
		usr.state = "getting-report"
	}
	if usr.state == "getting-report" {
		client := http.Client{}

		if b.upd.Result[0].Callback_query.Data == "xls" {
			b.db.getReport(b.answers[0], b.answers[1], usr.dbId)

			b.resetUser(usr)
			return
		}

		inline_keyboard := `{"inline_keyboard":[[{"text":"xls", "callback_data":"xls"},{"text":"txt", "callback_data":"txt"}]]}`

		tmp_url := "https://api.telegram.org/bot" + os.Getenv("TOKEN") + "/sendMessage?chat_id=" +
			strconv.Itoa(usr.chatId) +
			"&text=Твой отчет за период с " + b.answers[0] + " по " + b.answers[1] + ". " +
			"Выбери формат, в котором вывести отчет." +
			"&reply_markup=" + inline_keyboard

		req, _ := http.NewRequest(http.MethodPost, tmp_url, nil)
		_, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (b *Bot) checkDate(usr *Users, date string) bool {
	if match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, date); match {
		b.answers = append(b.answers, date)
		return true
	} else {
		b.sendTextMessage(usr.chatId, "Не правильно ты вводишь дату.")
		b.sendTextMessage(usr.chatId, "Попробуй в таком формате ГГГГ-ММ-ДД")
		usr.state = "check-date"
		return false
	}
}

func (b *Bot) resetUser(usr *Users) {
	usr.userTypingExpenses = false
	usr.userTypingIncome = false
	usr.userGettingReport = false
	usr.state = ""
	b.answers = make([]string, 0)
}

func (b *Bot) sendTextMessage(chat_id int, text string) {
	client := http.Client{}

	tmp_url := "https://api.telegram.org/bot" + os.Getenv("TOKEN") + "/sendMessage?chat_id=" +
		strconv.Itoa(chat_id) + "&text=" + text
	send_req, _ := http.NewRequest(http.MethodPost, tmp_url, nil)
	client.Do(send_req)
}

func (b *Bot) getHelp(filepath string, result chan string) {
	filepath = path.Join("./", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}
	defer b.recoverPanic()

	var buffer []byte = make([]byte, 834)
	file.Read(buffer)

	text := strings.ReplaceAll(string(buffer), string('\r'), "")
	text = strings.ReplaceAll(text, string('\n'), "")

	result <- text
}

func (b *Bot) recoverPanic() {
	err := recover()
	if err != nil {
		fmt.Printf("RECOVERED %v\n", err)
	}
}

// func (b *Bot) spinner() {
// 	for _, v := range `-\|/` {
// 		fmt.Printf("\r%c", v)
// 		time.Sleep(100 * time.Millisecond)
// 	}
// }
