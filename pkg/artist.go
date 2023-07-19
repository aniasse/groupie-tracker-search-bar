package pkg

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Artist struct {
	Aid        int      `json:"id"`
	AImg       string   `json:"image"`
	Aname      string   `json:"name"`
	Amember    []string `json:"members"`
	Acread     int      `json:"creationDate"`
	Afalbum    string   `json:"firstAlbum"`
	Aloc       string   `json:"locations"`
	Aconcerdat string   `json:"concertDates"`
	Arelat     string   `json:"relations"`
}

type OneArtist struct {
	Id         int
	Img        string
	Name       string
	Member     []string
	CreatDat   int
	FirstAlbum string
	DatLoc     map[string][]string
	ConcertDat []string
	Artists    []Artist
	LocatFilt  []string
}

func GetArtistData(api string) []Artist {

	content, err := http.Get(api)
	if err != nil {
		fmt.Println("Erreur de recuperation des donnees", err)
	}
	defer content.Body.Close()

	scan, er := ioutil.ReadAll(content.Body)
	if er != nil {
		fmt.Println("Erreur lors de la lecture des donnees ", er)
	}
	var Artists []Artist
	erreur := json.Unmarshal([]byte(scan), &Artists)
	if erreur != nil {
		fmt.Println("Erreur lors du decodage", erreur)
	}

	return Artists
}

func HandleArtist(w http.ResponseWriter, r *http.Request) {

	loc := GetLocationData("https://groupietrackers.herokuapp.com/api/locations")
	tabloc := TabLoc(loc)

	Art := GetArtistData("https://groupietrackers.herokuapp.com/api/artists")

	FilterArtist := Filter{
		Artists:   Art,
		LocatFilt: tabloc,
	}
	temp := template.Must(template.ParseFiles("templates/artist.html", "templates/form.html", "templates/navbar.html"))
	err := temp.Execute(w, FilterArtist)
	if err != nil {
		fmt.Println("Erreur lors de l'execution du template", err)
	}

}

func HandleArtistDeatail(w http.ResponseWriter, r *http.Request) {

	artistID := r.URL.Query().Get("Id")
	artistid, _ := strconv.Atoi(artistID)

	artist := Artist{}
	if artistid < 1 || artistid > 52 {
		Error404Handler(w, r)
	} else {

		Artists := GetArtistData("https://groupietrackers.herokuapp.com/api/artists")
		artist = Artists[artistid-1]

		relation := artist.Arelat
		data_rel := Rel_Artist(relation)
		date := artist.Aconcerdat
		data_dat := Dat_Artist(date)

		dat := data_dat.DAdat
		newdat := []string{}
		for _, v := range dat {
			newdat = append(newdat, strings.ReplaceAll(v, "*", ""))
		}
		rel := data_rel.RAdatloc

		newrel := map[string][]string{}
		for key, val := range rel {
			newkey := strings.ReplaceAll(key, "-", "\n")
			newrel[newkey] = val
		}

		MyArtist := []Artist{}

		for _, v := range Artists {
			if v.Aid != artistid {
				MyArtist = append(MyArtist, v)
			}
		}

		NewArtist := OneArtist{
			Img:        artist.AImg,
			Name:       artist.Aname,
			Member:     artist.Amember,
			CreatDat:   artist.Acread,
			FirstAlbum: artist.Afalbum,
			DatLoc:     newrel,
			ConcertDat: newdat,
			Artists:    MyArtist,
			LocatFilt:  TabLoc(GetLocationData("https://groupietrackers.herokuapp.com/api/locations")),
		}

		temp := template.Must(template.ParseFiles("templates/artist_detail.html", "templates/form.html", "templates/navbar.html"))
		temp.Execute(w, NewArtist)
	}
}
