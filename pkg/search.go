package pkg

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func HandleSearch(w http.ResponseWriter, r *http.Request) {

	search := r.FormValue("search")

	Artsearch := []Artist{}

	Artist := GetArtistData("https://groupietrackers.herokuapp.com/api/artists")
	tabId := []int{}
	tabloc := TabLoc(GetLocationData("https://groupietrackers.herokuapp.com/api/locations"))
	temp := template.Must(template.ParseFiles("templates/resultsearch.html", "templates/navbar.html", "templates/form.html"))
	var err error
	if !CheckURL(tabloc, search) {
		MyResult := ResultLocation(search)
		err = temp.Execute(w, MyResult)
	} else {

		for _, Art := range Artist {

			if search == Art.Aname {
				tabId = append(tabId, Art.Aid)
			} else if search == strconv.Itoa(Art.Acread) {
				tabId = append(tabId, Art.Aid)
			} else if search == Art.Afalbum {
				tabId = append(tabId, Art.Aid)
			} else {
				for _, v := range Art.Amember {
					if v == search {
						tabId = append(tabId, Art.Aid)
					}
				}
			}
		}
		for _, id := range tabId {
			Artsearch = append(Artsearch, Artist[id-1])
		}

		if len(Artsearch) == 0 {
			NoResultSearch := Filter{
				Boul:      true,
				NoResult:  "NO RESULTS FOR THE INFORMATION ENTERED",
				LocatFilt: tabloc,
			}
			err = temp.Execute(w, NoResultSearch)
		} else {
			ResultSearch := Filter{
				Artboul:   true,
				Artists:   Artsearch,
				LocatFilt: tabloc,
			}
			err = temp.Execute(w, ResultSearch)
		}
	}

	if err != nil {
		fmt.Println("Erreur lors de l'execution du template", err)
		Error500Handler(w, r)
		return
	}
}
