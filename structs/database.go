/*
	go get github.com/lib/pq
	_ "github.com/lib/pq"

	go get github.com/xuri/excelize
	go get github.com/xuri/excelize/v2
	"github.com/xuri/excelize/v2"
*/

package structs

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

type Database struct {
	driver           string
	connectionString string
}

func (d *Database) connect() *sql.DB {
	conn, err := sql.Open(d.driver, d.connectionString)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func (d *Database) InsertUser(id int, name string) {
	conn := d.connect()
	defer conn.Close()

	result, err := conn.Exec("INSERT INTO users (chat_id, first_name) VALUES ($1, $2)", id, name)
	if err != nil {
		fmt.Println(err)
	}

	defer d.recoverPanic()

	rowsAffectedCount, _ := result.RowsAffected()
	fmt.Printf("Rows affected %d\n", rowsAffectedCount)
}

func (d *Database) GetAllUsers() []*Users {
	conn := d.connect()
	defer conn.Close()

	rows, err := conn.Query("SELECT user_id, chat_id, first_name FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var id int
	var chatId int
	var name string
	tmpSlice := make([]*Users, 0)

	for rows.Next() {
		err := rows.Scan(&id, &chatId, &name)
		if err != nil {
			fmt.Println(err)
		}
		obj := Users{dbId: id, chatId: chatId, firstName: name}
		tmpSlice = append(tmpSlice, &obj)
	}
	return tmpSlice
}

func (d *Database) InsertEntry(tablename, category string, item string, price int, date string, id int) {
	var result sql.Result
	var err error

	conn := d.connect()
	defer conn.Close()

	if tablename == "expenses" {
		result, err = conn.Exec(`INSERT INTO expenses (category, item_name, price, buy_date, user_id) VALUES ($1, $2, $3, $4, $5)`, category, item, price, date, id)
		if err != nil {
			log.Fatal(err)
		}
	} else if tablename == "income" {
		result, err = conn.Exec(`INSERT INTO income (category, item_name, price, buy_date, user_id) VALUES ($1, $2, $3, $4, $5)`, category, item, price, date, id)
		if err != nil {
			log.Fatal(err)
		}
	}

	r, _ := result.RowsAffected()
	fmt.Printf("Rows affected %d\n", r)
}

func (d *Database) getReport(startDate, endDate string, userId int) {
	expensesSlice, incomeSlice := d.createReport(startDate, endDate, userId)

	f := excelize.NewFile()

	incomeSheet, err := f.NewSheet("Income")
	if err != nil {
		fmt.Println(err)
	}
	f.SetActiveSheet(incomeSheet)

	for i := 0; i < len(incomeSlice); i++ {
		for j := 0; j < len(incomeSlice[i]); j++ {
			cell, _ := excelize.CoordinatesToCellName(i+1, j+1)
			f.SetCellValue("Income", cell, incomeSlice[i][j])
		}
	}

	expensesSheet, err := f.NewSheet("Expenses")
	if err != nil {
		fmt.Println(err)
	}
	f.SetActiveSheet(expensesSheet)

	for i := 0; i < len(expensesSlice); i++ {
		for j := 0; j < len(expensesSlice[i]); j++ {
			cell, _ := excelize.CoordinatesToCellName(i+1, j+1)
			f.SetCellValue("Expenses", cell, expensesSlice[i][j])
		}
	}

	f.SaveAs("./reports/report.xlsx")
}

func (d *Database) createReport(startDate, endDate string, userId int) ([][]string, [][]string) {
	conn := d.connect()
	defer conn.Close()

	var rows *sql.Rows
	var err error
	var size int
	var category, item_name, buy_date string
	var price, dbId, userid, index int

	// expenses report
	rows, err = conn.Query(`SELECT * FROM expenses WHERE buy_date BETWEEN $1 AND $2 AND user_id=$3`, startDate, endDate, userId)
	if err != nil {
		log.Fatal(err)
	}

	countRows := conn.QueryRow(`SELECT COUNT(*) FROM expenses WHERE buy_date BETWEEN $1 AND $2 AND user_id=$3`, startDate, endDate, userId)
	countRows.Scan(&size)

	var expensesSlice [][]string = make([][]string, size)

	for rows.Next() {
		err := rows.Scan(&dbId, &category, &item_name, &price, &buy_date, &userid)
		if err != nil {
			fmt.Println(err)
			continue
		}
		expensesSlice[index] = append(expensesSlice[index], category, item_name, strconv.Itoa(price), buy_date)
		index++
	}

	index = 0

	// income report
	rows, err = conn.Query(`SELECT * FROM income WHERE buy_date BETWEEN $1 AND $2 AND user_id=$3`, startDate, endDate, userId)
	if err != nil {
		log.Fatal(err)
	}

	countRows = conn.QueryRow(`SELECT COUNT(*) FROM income WHERE buy_date BETWEEN $1 AND $2 AND user_id=$3`, startDate, endDate, userId)
	countRows.Scan(&size)

	var incomeSlice [][]string = make([][]string, size)

	for rows.Next() {
		err = rows.Scan(&dbId, &category, &item_name, &price, &buy_date, &userid)
		if err != nil {
			fmt.Println(err)
			continue
		}
		incomeSlice[index] = append(incomeSlice[index], category, item_name, strconv.Itoa(price), buy_date)
		index++
	}
	return expensesSlice, incomeSlice
}

func (d *Database) recoverPanic() {
	if err := recover(); err != nil {
		fmt.Printf("RECOVERED %v\n", err)
	}
}
