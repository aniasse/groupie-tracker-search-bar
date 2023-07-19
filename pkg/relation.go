package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type IndexRelat struct {
	IRid     int                 `json:"id"`
	IRdatloc map[string][]string `json:"datesLocations"`
}

type Relation struct {
	Relat []IndexRelat `json:"index"`
}

func GetRelationData(api_relation string) Relation {

	data_api, err := http.Get(api_relation)
	if err != nil {
		fmt.Println("Erreur lors de la recuperation des donnees", err)
	}
	defer data_api.Body.Close()

	scan, er := ioutil.ReadAll(data_api.Body)
	if er != nil {
		fmt.Println("Erreur lors de la lecture des donnees", er)
	}
	var RelationArtist Relation

	erreur := json.Unmarshal([]byte(scan), &RelationArtist)
	if erreur != nil {
		fmt.Println("Erreur lors du decodage des donnees JSON", erreur)
	}

	return RelationArtist
}
