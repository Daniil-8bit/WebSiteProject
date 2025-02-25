package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func main() {
	apiKey = flag.String("apikey", "", "api key for newsapi.org")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("api key must be set!")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "3333"
	}

	mux := http.NewServeMux()

	//fs := http.FileServer(http.Dir(""))
	//mux.Handle("/", http.StripPrefix("", fs))

	mux.HandleFunc("/", startPageHandler)
	mux.HandleFunc("/general", generalPageHandler)
	mux.HandleFunc("/contacts", contactsPageHandler)
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/login", loginPageHandler)
	mux.HandleFunc("/checklogin", checkLoginHandler)
	mux.HandleFunc("/register", registerPageHandler)
	mux.HandleFunc("/registerData", registerUserHandler)
	mux.HandleFunc("/index", indexHandler)

	http.ListenAndServe(":"+port, mux)
}

var tpl = template.Must(template.ParseFiles("index.html"))
var tplContacts = template.Must(template.ParseFiles("contacts.html"))
var tplGeneral = template.Must(template.ParseFiles("general.html"))
var tplLogin = template.Must(template.ParseFiles("login.html"))
var tplRegister = template.Must(template.ParseFiles("register.html"))
var apiKey *string

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("<h1>Hello new web site!</h1>"))

	tpl.Execute(w, nil)
}

func generalPageHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("<p>This is general page!</p>"))

	tplGeneral.Execute(w, nil)
}

func contactsPageHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("<p>Contacts page!</p>"))

	tplContacts.Execute(w, nil)
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	fmt.Printf("User login: %s, password: %s\n", login, password)

	success := 0
	if login != "" && password != "" {
		if checkLoginData(login, password) {
			success = 1
		} else {
			success = 2
		}
	} else {
		success = 0
	}
	fmt.Println("Success: ", success)

	if success == 1 {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	tplLogin.Execute(w, success)
}

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	//err := r.ParseForm()

	/*if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}*/

	login := r.FormValue("login")
	password := r.FormValue("password")

	if login != "" && password != "" {
		usersLoginData[login] = password
	}

	fmt.Println(login, ":", usersLoginData[login])

	tplRegister.Execute(w, nil)
}

func startPageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	result, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := result.Query()
	search := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	//fmt.Println("Search: ", search)
	//fmt.Println("Page: ", page)

	searchLine := &Search{}
	searchLine.SearchKey = search

	nextPage, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Internal server error: ", http.StatusInternalServerError)
	}

	searchLine.NextPage = nextPage
	pageSize := 20

	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%d&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(searchLine.SearchKey), pageSize, searchLine.NextPage, *apiKey)
	resp, err := http.Get(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&searchLine.Results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	searchLine.TotalPages = int(math.Ceil(float64(searchLine.Results.TotalResults / pageSize)))

	if ok := !searchLine.IsLastPage(); ok {
		searchLine.NextPage++
	}

	err = tpl.Execute(w, searchLine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type Source struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

type Article struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Search struct {
	SearchKey  string
	NextPage   int
	TotalPages int
	Results    Results
}

func (a *Article) FormatPublishedDate() string {
	year, month, day := a.PublishedAt.Date()
	return fmt.Sprintf("%d %v, %d", day, month, year)
}

func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}

func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}

/*type loginData struct {
	Login    string
	password string
}*/

var LoginStatus string = "Some text!"

func checkLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	fmt.Printf("User login: %s, password: %s\n", login, password)

	//success := checkLoginData(login, password)

	http.Redirect(w, r, "/login", http.StatusOK)

	/*if success {
		http.Redirect(w, r, "/general", http.StatusOK)
		LoginStatus = ""
		fmt.Println("Success: ", success)
		fmt.Println("login status: ", LoginStatus)
	} else {
		http.Redirect(w, r, "/login", http.StatusOK)
		LoginStatus = "Login or password is incorrect!"
		fmt.Println("Success: ", success)
		fmt.Println("login status: ", LoginStatus)
	}*/
}

var usersLoginData map[string]string = map[string]string{"admin": "1234"}

func checkLoginData(login string, password string) bool {

	value, ok := usersLoginData[login]
	if !ok {
		fmt.Println("No user with login: ", login)
		return false
	}

	if password == value {
		fmt.Println("Correct!")
		return true
	}

	return false
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	usersLoginData[login] = password
	fmt.Println(login, ":", usersLoginData[login])

	//w.Write([]byte("<h1>You have successfully registered!</h1>"))

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
