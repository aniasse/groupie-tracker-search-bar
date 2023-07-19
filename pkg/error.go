package pkg

import (
	"fmt"
	"html/template"
	"net/http"
)

// error500Handler manage the HTTP Response with an error 500 (Internal Server Error).
func Error500Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl := template.Must(template.ParseFiles("templates/error500.html"))
	tmpl.Execute(w, nil)
}

// error404Handler manage the HTTP Response with an error 404 (Not Found).
func Error404Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tmpl := template.Must(template.ParseFiles("templates/error404.html"))
	tmpl.Execute(w, nil)
}

// error400Handler manage the HTTP Response with an error 400 (Bad Request).
func Error400Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl := template.Must(template.ParseFiles("templates/error400.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError)
		return
	}
}

// error405Handler manage the HTTP Response with an error 405 (Method Not Allowed).
func Error405Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	tmpl := template.Must(template.ParseFiles("templates/error405.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		HandleError(w, r, http.StatusMethodNotAllowed)
		return
	}
}

// handleError handles HTTP errors and redirects to the appropriate error handlers.
func HandleError(w http.ResponseWriter, r *http.Request, statusCode int) {

	switch statusCode {
	case http.StatusNotFound:
		Error404Handler(w, r)
	case http.StatusInternalServerError:
		// Retrieve the server's internal error code
		errorCode := http.StatusInternalServerError

		// Your treatment according to the internal error code
		switch errorCode {
		case http.StatusInternalServerError:
			// Manage the 500 error (Internal Server Error)
			Error500Handler(w, r)
		default:
			// Manage the other errors internes
			http.Error(w, "Autre erreur interne", errorCode)
		}
	case http.StatusBadRequest:
		Error400Handler(w, r)
	case http.StatusMethodNotAllowed:
		Error405Handler(w, r)
	default:
		// Manage the unexpected error
		Error500Handler(w, r)
	}
}

// errorHandler is a middleware application for handling errors when processing HTTP requests.
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Erreur interne du serveur:", r)
				Error500Handler(w, nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
