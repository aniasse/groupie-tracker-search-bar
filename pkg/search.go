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
	relation := GetRelationData("https://groupietrackers.herokuapp.com/api/relation")
	tabId := []int{}
	ind := 0

	for i := ind; i < len(Artist); i++ {
		for i = ind; i < len(relation.Relat); i++ {
			if search == Artist[i].Aname {
				tabId = append(tabId, Artist[i].Aid)
			} else if search == strconv.Itoa(Artist[i].Acread) {
				tabId = append(tabId, Artist[i].Aid)
			} else if search == Artist[i].Afalbum {
				tabId = append(tabId, Artist[i].Aid)
			} else {
				for _, v := range Artist[i].Amember {
					if v == search {
						tabId = append(tabId, Artist[i].Aid)
					}
				}
				for loc := range relation.Relat[i].IRdatloc {
					if loc == search {
						if loc == search {
							tabId = append(tabId, relation.Relat[i].IRid)
						}
					}
				}
			}
		}
		ind++
	}

	for _, id := range tabId {
		Artsearch = append(Artsearch, Artist[id-1])
	}

	temp := template.Must(template.ParseFiles("templates/resultsearch.html", "templates/navbar.html", "templates/form.html"))
	var err error
	if len(Artsearch) == 0 {
		NoResultSearch := Filter{
			Boul:      true,
			NoResult:  "NO RESULTS FOR THE INFORMATION ENTERED",
			LocatFilt: TabLoc(GetLocationData("https://groupietrackers.herokuapp.com/api/locations")),
		}
		err = temp.Execute(w, NoResultSearch)
	} else {
		ResultSearch := Filter{
			Artists:   Artsearch,
			LocatFilt: TabLoc(GetLocationData("https://groupietrackers.herokuapp.com/api/locations")),
		}
		err = temp.Execute(w, ResultSearch)
	}

	if err != nil {
		fmt.Println("Erreur lors de l'execution du template", err)
		Error500Handler(w, r)
		return
	}
}
