/*
	https://api.telegram.org/bot6523339653:AAF1mrO7mTD0JBbV9vhie-FG4dmBBW6OlOI/getUpdates
*/

package structs

// === Update struct start ===
type Update struct {
	Ok     bool
	Result []Result
}

type Result struct {
	Update_id      int
	Message        Message
	Callback_query Callback_query
}

type Callback_query struct {
	Data string
}

type Message struct {
	MessageId int
	From      From
	Chat      Chat
	Text      string
}

type From struct {
	Id         int
	First_name string
	Username   string
}

type Chat struct {
	Id         int
	First_name string
	Username   string
}

// === Update struct end ===

// === Users struct start ===
type Users struct {
	dbId               int
	chatId             int
	firstName          string
	userTypingExpenses bool
	userTypingIncome   bool
	userGettingReport  bool
	state              string
}

// === Users struct end ===
