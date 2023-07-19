package pkg

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Index struct {
	Loc_id int      `json:"id"`
	Loca   []string `json:"locations"`
	Dat    string   `json:"dates"`
}
type Location struct {
	Loc []Index `json:"index"`
}
type Localisation struct {
	Id        []int
	Local     []string
	Artists   []Artist
	LocatFilt []string
}
type LocationInfo struct {
	Loc        string
	Id         int
	Artistimg  []string
	Artistname []string
	Date       [][]string
	GlobLoc    []string
	locat      []string
	Artists    []Artist
	LocatFilt  []string
	Boul       bool
	Artboul    bool
}

var NewLocalisation Localisation

func GetLocationData(api_loc string) Location {

	data_api, err := http.Get(api_loc)
	if err != nil {
		fmt.Println("Erreur de recuperation des donnees", err)
	}
	defer data_api.Body.Close()

	scanner, er := ioutil.ReadAll(data_api.Body)
	if er != nil {
		fmt.Println("Erreur lors de la lecture", er)
	}

	var ArtistLoc Location
	erreur := json.Unmarshal([]byte(scanner), &ArtistLoc)
	if erreur != nil {
		fmt.Println("Erreur lors du decodage JSON", erreur)
	}
	return ArtistLoc
}

func HandleLocation(w http.ResponseWriter, r *http.Request) {

	location := GetLocationData("https://groupietrackers.herokuapp.com/api/locations")

	tabloc := TabLoc(location)

	NewLocalisation = Localisation{
		Local:     tabloc,
		Artists:   GetArtistData("https://groupietrackers.herokuapp.com/api/artists"),
		LocatFilt: tabloc,
	}

	temp := template.Must(template.ParseFiles("templates/location.html", "templates/navbar.html"))
	err := temp.Execute(w, NewLocalisation)
	if err != nil {
		fmt.Println("Erreur lors de l'execution du template", err)
	}
}

func HandleLocationDetail(w http.ResponseWriter, r *http.Request) {

	one_location := r.URL.Query().Get("Loc")

	loc := NewLocalisation.Local
	if CheckURL(loc, one_location) {
		Error404Handler(w, r)
		return
	} else {
		MyLocationInfo := ResultLocation(one_location)
		temp := template.Must(template.ParseFiles("templates/location_detail.html", "templates/navbar.html"))
		err := temp.Execute(w, MyLocationInfo)
		if err != nil {
			fmt.Println("Erreur lors de l'execution du template", err)
		}
	}

}

func TabLoc(loc Location) []string {

	tabloc := []string{}

	for _, v := range loc.Loc {
		for i := 0; i < len(v.Loca); i++ {
			if !NoRepeatLoc(tabloc, v.Loca[i]) {
				tabloc = append(tabloc, v.Loca[i])
			}
		}
	}
	for i := 0; i < len(tabloc); i++ {
		for j := i + 1; j < len(tabloc); j++ {
			if tabloc[i] > tabloc[j] {
				swap := tabloc[i]
				tabloc[i] = tabloc[j]
				tabloc[j] = swap
			}
		}
	}
	return tabloc
}

func NoRepeatLoc(tab []string, str string) bool {

	for _, v := range tab {
		if v == str {
			return true
		}
	}
	return false
}

func CheckURL(tab []string, str string) bool {

	count := 0
	for _, v := range tab {
		if str != v {
			count++
		}
	}
	if count == len(tab) {
		return true
	}

	return false
}

func ResultLocation(loc string) LocationInfo {

	relation := GetRelationData("https://groupietrackers.herokuapp.com/api/relation")

	var (
		datlocation [][]string
		tab         []string
		IdLoc       []int
	)

	for _, v := range relation.Relat {
		if v.IRdatloc[loc] != nil {
			tab = v.IRdatloc[loc]
			datlocation = append(datlocation, tab)
			IdLoc = append(IdLoc, v.IRid)
		}
	}
	artists := GetArtistData("https://groupietrackers.herokuapp.com/api/artists")
	ArtistLoc := []Artist{}

	for _, id := range IdLoc {
		ArtistLoc = append(ArtistLoc, artists[id-1])
	}

	stockname := []string{}
	stockimg := []string{}

	for _, v := range ArtistLoc {
		stockname = append(stockname, v.Aname)
		stockimg = append(stockimg, v.AImg)
	}
	MyLocation := []string{}
	for _, v := range NewLocalisation.Local {
		if v != loc {
			MyLocation = append(MyLocation, v)
		}
	}
	NewLocationInfo := LocationInfo{
		Loc:        loc,
		Artistimg:  stockimg,
		Artistname: stockname,
		Date:       datlocation,
		GlobLoc:    MyLocation,
		locat:      NewLocalisation.Local,
		Artists:    GetArtistData("https://groupietrackers.herokuapp.com/api/artists"),
		LocatFilt:  TabLoc(GetLocationData("https://groupietrackers.herokuapp.com/api/locations")),
	}
	return NewLocationInfo
}
