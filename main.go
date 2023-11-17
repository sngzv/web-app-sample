package main

import (
	"encoding/json"
	"flag" //анализ базовых флогов командой строки
	"fmt"
	"html/template" //шаблонизация
	"log"
	"math"
	"net/http" //реализация клиента и сервера HTTP
	"net/url"
	"os"
	"strconv"
	"time"
)

// определение шаблона
var tpl = template.Must(template.ParseFiles("index.html"))
var apiKey *string

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("<h1>Hello world!</h1>")) //w - отвечает за отпрвку ответов на HTTP-запрос
	tpl.Execute(w, nil) //записываем выходные данные в интерфейс
}

// обработчик поисковых запросов
func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}

	params := u.Query()
	searchKey := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	search := &Search{}
	search.SearchKey = searchKey

	next, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Unexpected server error", http.StatusInternalServerError)
	}

	search.NextPage = next
	pageSize := 20

	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%d&apiKey=%s&sortBy=publishedAt&language=en",
		url.QueryEscape(search.SearchKey), pageSize, search.NextPage, *apiKey)

	resp, err := http.Get(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = json.NewDecoder(resp.Body).Decode(&search.Results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	search.TotalPages = int(math.Ceil(float64(search.Results.TotalResults / pageSize)))
	err = tpl.Execute(w, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func main() {
	apiKey = flag.String("apiKey", "", "Newsapi.org access key") //определение строкового флага
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("apiKey must be set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	mux := http.NewServeMux() //новый мультиплексор HTTP-запросов

	fs := http.FileServer(http.Dir("assets"))                //создание файлового сервера
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs)) //использования для всех путей начинающихся с префика /assets/

	mux.HandleFunc("/", indexHandler) //регистрация функции для пути "/"
	mux.HandleFunc("/search", searchHandler)
	http.ListenAndServe(":"+port, mux) //запуск сервера
}

type Source struct {
	ID   interface{}
	Name string
}

type Article struct {
	Source      Source
	Author      string
	Title       string
	Description string
	URL         string
	URLToImage  string
	PublishedAt time.Time
	Content     string
}

type Results struct {
	Status       string
	TotalResults int
	Articles     []Article
}

type Search struct {
	SearchKey  string
	NextPage   int
	TotalPages int
	Results    Results
}
