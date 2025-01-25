package main

import (
	"html/template"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3333"
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(""))
	mux.Handle("/", http.StripPrefix("", fs))

	//mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/general", generalPageHandler)
	mux.HandleFunc("/contacts", contactsPageHandler)

	http.ListenAndServe(":"+port, mux)
}

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("<h1>Hello new web site!</h1>"))

	tpl.Execute(w, nil)
}

func generalPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h3>This is general page!</h3>"))
}

func contactsPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<p>Contacts page!</p>"))
}
