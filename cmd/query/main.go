package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "$41yohul22"
	dbname   = "web_8"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

// Обработчики HTTP-запросов
func (h *Handlers) GetHello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprintf(w, "Hello, sweetheart<3")
		return
	}
	msg, err := h.dbProvider.SelectHello(name)
	if msg == false {
		fmt.Fprintf(w, "this user dont exist((")
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello, " + name + "<3"))
}
func (h *Handlers) PostHello(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprintf(w, "Hello, sweetheart<3")
		return
	}

	err := h.dbProvider.InsertHello(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
}

// Методы для работы с базой данных
func (dp *DatabaseProvider) SelectHello(name string) (bool, error) {
	var yeano string

	// Получаем одно сообщение из таблицы hello, отсортированной в случайном порядке
	querry := "SELECT name FROM names WHERE  name = ($1)"
	err := dp.db.QueryRow(querry, name).Scan(&yeano)
	if err != nil {
		return false, err
	}

	return true, nil
}
func (dp *DatabaseProvider) InsertHello(msg string) error {
	_, err := dp.db.Exec("INSERT INTO names(name) VALUES ($1)", msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Считываем аргументы командной строки
	address := flag.String("address", "127.0.0.1:8081", "адрес для запуска сервера")
	flag.Parse()

	// Формирование строки подключения для postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Создание соединения с сервером postgres
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем провайдер для БД с набором методов
	dp := DatabaseProvider{db: db}
	// Создаем экземпляр структуры с набором обработчиков
	h := Handlers{dbProvider: dp}

	// Регистрируем обработчики
	http.HandleFunc("/api/user/get", h.GetHello)
	http.HandleFunc("/api/user/post", h.PostHello)

	// Запускаем веб-сервер на указанном адресе
	err = http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
