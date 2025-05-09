package controllers

import (
	"net/http"

	"rwa.abc/views"
)

type Stock struct {
	Name      string
	Templates struct {
		New   views.Template
		Query views.Template
	}
}

type StockData struct {
	StockName string
}

func (s Stock) NewQuery(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	// Render the template with the parsed form data
	s.Templates.New.Execute(w, nil)
}

func (s Stock) StockQuery(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	data := StockData{
		StockName: r.FormValue("ticker"),
	}
	println("Stock Name You want to query is: ", data.StockName)
	//fmt.Fprint(w, "Stock Name You want to query is: ", stockName)
	s.Templates.Query.Execute(w, data)
}
