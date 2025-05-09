package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"rwa.abc/controllers"
	"rwa.abc/templates"
	"rwa.abc/views"
)

// main function sets up the router and starts the server
func main() {
	r := chi.NewRouter()

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/stocks", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "stocks.gohtml", "tailwind.gohtml"))))

	// r.Get("/stockquery", controllers.StaticHandler(
	// 	views.Must(views.ParseFS(templates.FS, "stockquery.gohtml", "tailwind.gohtml"))))

	stockC := controllers.Stock{}
	stockC.Templates.New = views.Must(views.ParseFS(templates.FS, "stockquery.gohtml", "tailwind.gohtml"))
	stockC.Templates.Query = views.Must(views.ParseFS(templates.FS, "stockqueryresponse.gohtml", "tailwind.gohtml"))
	r.Get("/stockquery", stockC.NewQuery)
	r.Post("/stockdetails", stockC.StockQuery)

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "404 - Page Not Found")
	})
	fmt.Println("Server is running on port 3000.")
	http.ListenAndServe(":3000", r)
}
