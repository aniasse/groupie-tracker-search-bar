package main

import (
	"fmt"
	"groupie-tracker-filters/pkg"
	"net/http"
)

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", pkg.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			pkg.HandleFilter(w, r)
		case "/artists":
			pkg.HandleArtist(w, r)
		case "/artist-details":
			pkg.HandleArtistDeatail(w, r)
		case "/filter":
			pkg.HandleFilterDetail(w, r)
		case "/search":
			pkg.HandleSearch(w, r)
		case "/locations":
			pkg.HandleLocation(w, r)
		case "/location-detail":
			pkg.HandleLocationDetail(w, r)
		case "/dates":
			pkg.HandleDAte(w, r)
		case "/date-infos":
			pkg.HandleDateInfo(w, r)
		default:
			pkg.Error404Handler(w, r)
		}

	})))

	fmt.Println("The program is running on http://localhost:4040")
	http.ListenAndServe("localhost:4040", nil)

}
