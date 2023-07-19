package pkg

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Filter struct {
	Artists   []Artist
	LocatFilt []string
	NoResult  string
	Boul      bool
}

func HandleFilter(w http.ResponseWriter, r *http.Request) {

	location := GetLocationData("https://groupietrackers.herokuapp.com/api/locations")

	tabloc := TabLoc(location) //Recuperation et arrangemet des locations

	NewFilter := Filter{
		Artists:   GetArtistData("https://groupietrackers.herokuapp.com/api/artists"),
		LocatFilt: tabloc,
	}

	err := template.Must(template.ParseFiles("templates/home.html", "templates/form.html", "templates/navbar.html")).Execute(w, NewFilter)
	if err != nil {
		fmt.Println("Erreur de l'execution du template", err)
	}
}

func HandleFilterDetail(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/filter" && strings.ToUpper(r.Method) != http.MethodPost {
		Error405Handler(w, r)
	} else {

		//Les champs coches par l'utiisateur
		creat_date := r.FormValue("creationdate")
		first_album := r.FormValue("firstalbum")
		checkb_members := r.FormValue("members")
		checkb_location := r.FormValue("location")

		Art := GetArtistData("https://groupietrackers.herokuapp.com/api/artists")
		var creation_date []int
		for _, v := range Art {
			if !NoRepeatInt(creation_date, v.Acread) {
				creation_date = append(creation_date, v.Acread)
			}
		}

		for i := 0; i < len(creation_date); i++ {
			for j := i + 1; j < len(creation_date); j++ {
				if creation_date[i] > creation_date[j] {
					swap := creation_date[i]
					creation_date[i] = creation_date[j]
					creation_date[j] = swap
				}
			}
		}

		album := []string{}
		album = append(album, "01-01-1962")
		for _, v := range Art {
			album = append(album, v.Afalbum)
		}
		album = append(album, "02-01-2026")

		// Convertir les dates en objets time.Time
		timeDates := make([]time.Time, len(album))
		for i, dateStr := range album {
			timeDate, _ := time.Parse("02-01-2006", dateStr)
			timeDates[i] = timeDate
		}

		// Tri du tableau de dates
		sort.Slice(timeDates, func(i, j int) bool {
			return timeDates[i].Before(timeDates[j])
		})
		formattedDate := []string{}
		// Put the date sort in a slice
		for _, date := range timeDates {
			formattedDate = append(formattedDate, date.Format("02-01-2006"))
		}

		Filt := []Artist{}
		var (
			dat_debut   int
			dat_fin     int
			debut_album time.Time
			final_album time.Time
			member1     int
			member2     int
		)
		if !Active(creat_date) {
			if len(creation_date) != 0 {
				dat_debut = creation_date[0]
				dat_fin = creation_date[len(creation_date)-1]
			}

		} else {
			cons1, _ := strconv.Atoi(r.FormValue("datdebut"))
			dat_debut = cons1
			cons2, _ := strconv.Atoi(r.FormValue("datfin"))
			dat_fin = cons2

		}
		if !Active(first_album) {

			if len(album) != 0 {
				d_al := formattedDate[0]
				f_al := formattedDate[len(formattedDate)-1]

				// Analyse de la chaîne de date en un objet time.Time
				date1, err := time.Parse("02-01-2006", d_al)
				date2, err := time.Parse("02-01-2006", f_al)
				if err != nil {
					fmt.Println("Erreur lors de l'analyse de la date :", err)
					return
				}
				// Formattez la date selon le nouveau format
				debut_album = date1
				final_album = date2
			}

		} else {
			d_al := r.FormValue("debutalbum")
			f_al := r.FormValue("finalalbum")

			// Analyse de la chaîne de date en un objet time.Time
			date1, err := time.Parse("2006-01-02", d_al)
			date2, err := time.Parse("2006-01-02", f_al)

			if err != nil {
				fmt.Println("Erreur lors de l'analyse de la date :", err)
				return
			}

			// Formattez la date selon le nouveau format
			debut_album = date1
			final_album = date2
		}

		Rel := GetRelationData("https://groupietrackers.herokuapp.com/api/relation")

		if !Active(checkb_members) && !Active(checkb_location) { //If the memmber and location are not checked

			for _, v := range Art {
				pars := v.Afalbum
				dat, err := time.Parse("02-01-2006", pars)
				if err != nil {
					fmt.Println("Erreur lors de l'analyse de la date :", err)
					return
				}

				if (v.Acread >= dat_debut && v.Acread <= dat_fin) && ((dat.After(debut_album) && dat.Before(final_album)) || (dat.Equal(debut_album) || dat.Equal(final_album))) {
					Filt = append(Filt, v)
				}
			}
		} else if Active(checkb_members) && !Active(checkb_location) { //If the member is checked and the location isn't checked

			consmem1, _ := strconv.Atoi(r.FormValue("member1"))
			consmem2, _ := strconv.Atoi(r.FormValue("member2"))

			member1 = consmem1
			member2 = consmem2

			for _, v := range Art {
				pars := v.Afalbum
				dat, err := time.Parse("02-01-2006", pars)
				if err != nil {
					fmt.Println("Erreur lors de l'analyse de la date :", err)
					return
				}
				if (v.Acread >= dat_debut && v.Acread <= dat_fin) && (len(v.Amember) >= member1 && len(v.Amember) <= member2) && ((dat.After(debut_album) && dat.Before(final_album)) || (dat.Equal(debut_album) || dat.Equal(final_album))) {

					Filt = append(Filt, v)
				}
			}

		} else if !Active(checkb_members) && Active(checkb_location) { //If the location is checked and the member isn't checked

			loca := r.FormValue("loc")
			ind := 0
			for i := ind; i < len(Art); i++ {
				for i = ind; i < len(Rel.Relat); i++ {
					pars := Art[i].Afalbum
					dat, err := time.Parse("02-01-2006", pars)
					if err != nil {
						fmt.Println("Erreur lors de l'analyse de la date :", err)
						return
					}
					for key := range Rel.Relat[i].IRdatloc {
						if (Art[i].Acread >= dat_debut && Art[i].Acread <= dat_fin) && ((dat.After(debut_album) && dat.Before(final_album)) || (dat.Equal(debut_album) || dat.Equal(final_album))) && key == loca {
							Filt = append(Filt, Art[i])

						}
					}
					ind++
				}
			}
		} else if Active(checkb_members) && Active(checkb_location) { //If the member and the location are checked

			loca := r.FormValue("loc")
			consmem1, _ := strconv.Atoi(r.FormValue("member1"))
			consmem2, _ := strconv.Atoi(r.FormValue("member2"))

			member1 = consmem1
			member2 = consmem2
			ind := 0
			for i := ind; i < len(Art); i++ {
				for i = ind; i < len(Rel.Relat); i++ {
					pars := Art[i].Afalbum
					dat, err := time.Parse("02-01-2006", pars)
					if err != nil {
						fmt.Println("Erreur lors de l'analyse de la date :", err)
						return
					}
					for key := range Rel.Relat[i].IRdatloc {
						if (Art[i].Acread >= dat_debut && Art[i].Acread <= dat_fin) && ((dat.After(debut_album) && dat.Before(final_album)) || (dat.Equal(debut_album) || dat.Equal(final_album))) && (len(Art[i].Amember) >= member1 && len(Art[i].Amember) <= member2) && (key == loca) {
							Filt = append(Filt, Art[i])

						}
						ind++
					}
				}
			}
		}
		location := GetLocationData("https://groupietrackers.herokuapp.com/api/locations")

		tabloc := TabLoc(location) //Recuperation et arrangemet des locations

		temp := template.Must(template.ParseFiles("templates/filter.html", "templates/form.html", "templates/navbar.html"))
		var err error
		if len(Filt) == 0 {
			NewNoresult := Filter{
				Boul:      true,
				NoResult:  "NO RESULTS FOR THE INFORMATION ENTERED",
				LocatFilt: tabloc,
			}
			err = temp.Execute(w, NewNoresult)
			return
		} else {
			NewFilterDetail := Filter{
				Artists:   Filt,
				LocatFilt: tabloc,
			}
			err = temp.Execute(w, NewFilterDetail)
		}
		if err != nil {
			fmt.Println("Erreur", err)
		}
	}
}

func Active(str string) bool {
	if str == "on" {
		return true
	}
	return false
}

func NoRepeatInt(tab []int, str int) bool {

	for _, v := range tab {
		if v == str {
			return true
		}
	}
	return false
}
