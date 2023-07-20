package pkg

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type DatIndex struct {
	Id       int      `json:"id"`
	IndexDat []string `json:"dates"`
}

type Date struct {
	Dat []DatIndex `json:"index"`
}

type DateEvent struct {
	Date      map[string][]string
	data      []string
	Artists   []Artist
	LocatFilt []string
}

type DateInfos struct {
	Artimg    []string
	Artname   []string
	Loc       []string
	GlobDat   map[string][]string
	Glob      []string
	Even      string
	Artists   []Artist
	LocatFilt []string
}

var (
	NewDatEvent DateEvent
	uniq_date   []string
	years       map[string][]string
)

func GetDateData(api_date string) Date {

	data_api, err := http.Get(api_date)
	if err != nil {
		fmt.Println("Erreur lors de la recuperation des donnees", err)
	}
	defer data_api.Body.Close()

	scan, er := ioutil.ReadAll(data_api.Body)
	if er != nil {
		fmt.Println("Erreur lors de la lecture des donnees", er)
	}
	var Dates Date

	erreur := json.Unmarshal([]byte(scan), &Dates)
	if erreur != nil {
		fmt.Println("Erreur lors du decodage JSON", erreur)
	}

	return Dates
}

func HandleDAte(w http.ResponseWriter, r *http.Request) {

	dates := GetDateData("https://groupietrackers.herokuapp.com/api/dates")

	for _, v := range dates.Dat {
		for i := 0; i < len(v.IndexDat); i++ {
			stock := strings.ReplaceAll(v.IndexDat[i], "*", "")
			if !NoRepeatLoc(uniq_date, stock) {
				uniq_date = append(uniq_date, stock)
			}
		}
	}
	year := []string{}

	years = make(map[string][]string)

	day_and_month := []string{}
	for _, v := range uniq_date {
		split := strings.Split(v, "-")
		stock := split[len(split)-1]
		fmt.Println(stock)
		if !NoRepeatLoc(year, stock) {
			year = append(year, stock)
		}
	}
	for _, s := range year {
		for _, v := range uniq_date {
			split := strings.Split(v, "-")
			if split[len(split)-1] == s {
				day_and_month = append(day_and_month, v)
			}
			years[s] = day_and_month
		}
		day_and_month = []string{}
	}

	NewDatEvent = DateEvent{
		Date:      years,
		data:      uniq_date,
		Artists:   GetArtistData("https://groupietrackers.herokuapp.com/api/artists"),
		LocatFilt: TabLoc(GetLocationData("https://groupietrackers.herokuapp.com/api/locations")),
	}

	err := template.Must(template.ParseFiles("templates/date.html", "templates/navbar.html")).Execute(w, NewDatEvent)
	if err != nil {
		fmt.Println("Erreur lors de l'execution du template", err)
	}

}

func HandleDateInfo(w http.ResponseWriter, r *http.Request) {

	user_date := r.URL.Query().Get("Dat")

	date := NewDatEvent.data

	if CheckURL(date, user_date) { //Checked the date of the client
		Error404Handler(w, r)
		return
	} else {

		Relation := GetRelationData("https://groupietrackers.herokuapp.com/api/relation")

		var (
			IdDat []int
			Locat []string
		)

		for _, v := range Relation.Relat {
			for key, dat := range v.IRdatloc {
				for i := 0; i < len(dat); i++ {
					if user_date == dat[i] {
						Locat = append(Locat, key)
						IdDat = append(IdDat, v.IRid)
					}
				}
			}
		}
		Artists := GetArtistData("https://groupietrackers.herokuapp.com/api/artists")
		var Artistsdat []Artist
		for _, id := range IdDat {
			Artistsdat = append(Artistsdat, Artists[id-1])
		}

		artistimg := []string{}
		artistname := []string{}

		for _, v := range Artistsdat {
			artistimg = append(artistimg, v.AImg)
			artistname = append(artistname, v.Aname)
		}

		NewDateInfo := DateInfos{
			Artimg:    artistimg,
			Artname:   artistname,
			Loc:       Locat,
			GlobDat:   years,
			Glob:      NewDatEvent.data,
			Even:      user_date,
			Artists:   GetArtistData("https://groupietrackers.herokuapp.com/api/artists"),
			LocatFilt: TabLoc(GetLocationData("https://groupietrackers.herokuapp.com/api/locations")),
		}

		erreur := template.Must(template.ParseFiles("templates/date_detail.html", "templates/navbar.html")).Execute(w, NewDateInfo)
		if erreur != nil {
			fmt.Println("Erreur lors de l'execution du template", erreur)
		}
	}

}
