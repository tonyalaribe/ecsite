package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"html/template"
	"log"
	"net/http"
)

type Product struct {
	Site        string
	Name        string
	Link        string
	Description string
	Price       int
	Image       string
	Features    []string
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("in search handler")

	q := r.FormValue("q")
	log.Println(q)

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	c := session.DB("echai").C("products")
	results := []Product{}
	c.Find(bson.M{
		"$text": bson.M{
			"$search": q,
		},
	}).Sort("price").All(&results)

	t, err := template.ParseFiles("template/index.html")

	if err != nil {
		log.Println(err)
	}

	dat := struct {
		Q       string
		Results []Product
	}{
		Q:       q,
		Results: results,
	}
	err = t.Execute(w, dat)
	if err != nil {
		log.Println(err)
	}

}

func main() {

	fs := http.FileServer(http.Dir("template"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", searchHandler)

	http.ListenAndServe(":8080", nil)

}
