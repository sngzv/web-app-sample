package main

import (
	"html/template" //шаблонизация
	"net/http"      //реализация клиента и сервера HTTP
	"os"
)

// определение шаблона
var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("<h1>Hello world!</h1>")) //w - отвечает за отпрвку ответов на HTTP-запрос
	tpl.Execute(w, nil) //записываем выходные данные в интерфейс
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	mux := http.NewServeMux() //новый мультиплексор HTTP-запросов

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", indexHandler)  //регистрация функции для пути "/"
	http.ListenAndServe(":"+port, mux) //запуск сервера
}
