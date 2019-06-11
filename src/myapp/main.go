package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

//Константы подключения к БД
const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "280883"
	dbname   = "postgres"
)

/* Определение структуры, в которую будут записываться значения из БД, для дальнейшего вывода на экран.
Комментарий типа "value is..." обязательно, без него будет выводиться ошибка создания структуры.*/
//Book is...
type Book struct {
	isbn   string
	title  string
	author string
	price  float32
}

//Основная функции
func main() {

	//Определяем строку с данными для подключения к БД
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	//Если подключение не удалось, выводим код и описание ошибки пользователю. Используем 'error' предоставляемый библиотекой 'pq'.
	if err, ok := err.(*pq.Error); ok {
		fmt.Println("pq error:", err.Code.Name())
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err, ok := err.(*pq.Error); ok {

		fmt.Println("pq error:", err.Code.Name())

		panic(err)
	}
	//Радуемся поддключению
	fmt.Println("Successfully connected!")

	//Функции работы с БД
	getValues(&db)
	addValues(&db)

}

//Функция получения данныйх из БД и вывод их на экран
func getValues(pdb **sql.DB) {

	db := *pdb

	//Выбираем запросом все поля и записываем их в переменную
	rows, err := db.Query("SELECT * FROM books")

	if err, ok := err.(*pq.Error); ok {
		fmt.Println("pq error:", err.Code.Name())
	}
	//Утечку памяти никто не отменял
	defer rows.Close()

	//Записывыем данные в слайс
	bks := make([]*Book, 0)
	for rows.Next() {
		bk := new(Book)
		err := rows.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price)
		if err != nil {
			log.Fatal(err)
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	//И в цикле выводим их на экран
	for _, bk := range bks {
		fmt.Printf("%s, %s, %s, £%.2f\n", bk.isbn, bk.title, bk.author, bk.price)
	}

}

//Функция добавления данных
func addValues(pdb **sql.DB) {

	db := *pdb

	//за добавление данных в БД со стороны Go отвечает функция 'Exec', а со стороны БД 'insert'
	_, err := db.Exec("insert into books (isbn, title, author, price) values ('101-1503261960', 'St.Pavel', 'Sten Maers', 6.5)")

	if err, ok := err.(*pq.Error); ok {
		fmt.Println("pq error:", err.Code.Name())
	} else {
		fmt.Println("String added to SQL.")
	}
}
