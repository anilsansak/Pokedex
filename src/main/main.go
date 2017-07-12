package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Type struct {
	// Name of the type
	Name string `json:"name"`
	// The effective basedata, damage multiplize 2x
	EffectiveAgainst []string `json:"effectiveAgainst"`
	// The weak basedata that against, damage multiplize 0.5x
	WeakAgainst []string `json:"weakAgainst"`
}
type arr []string

type Pokemon struct {
	Number         string   `json:"Number"`
	Name           string   `json:"Name"`
	Classification string   `json:"Classification"`
	TypeI          []string `json:"Type I"`
	TypeII         []string `json:"Type II,omitempty"`
	Weaknesses     []string `json:"Weaknesses"`
	FastAttackS    []string `json:"Fast Attack(s)"`
	Weight         string   `json:"Weight"`
	Height         string   `json:"Height"`
	Candy          struct {
		Name     string `json:"Name"`
		FamilyID int    `json:"FamilyID"`
	} `json:"Candy"`
	NextEvolutionRequirements struct {
		Amount int    `json:"Amount"`
		Family int    `json:"Family"`
		Name   string `json:"Name"`
	} `json:"Next Evolution Requirements,omitempty"`
	NextEvolutions []struct {
		Number string `json:"Number"`
		Name   string `json:"Name"`
	} `json:"Next evolution(s),omitempty"`
	PreviousEvolutions []struct {
		Number string `json:"Number"`
		Name   string `json:"Name"`
	} `json:"Previous evolution(s),omitempty"`
	SpecialAttacks      []string `json:"Special Attack(s)"`
	BaseAttack          int      `json:"BaseAttack"`
	BaseDefense         int      `json:"BaseDefense"`
	BaseStamina         int      `json:"BaseStamina"`
	CaptureRate         float64  `json:"CaptureRate"`
	FleeRate            float64  `json:"FleeRate"`
	BuddyDistanceNeeded int      `json:"BuddyDistanceNeeded"`
}

// Move is an attack information. The
type Move struct {
	// The ID of the move
	ID int `json:"id"`
	// Name of the attack
	Name string `json:"name"`
	// Type of attack
	Type string `json:"type"`
	// The damage that enemy will take
	Damage int `json:"damage"`
	// Energy requirement of the attack
	Energy int `json:"energy"`
	// Dps is Damage Per Second
	Dps float64 `json:"dps"`
	// The duration
	Duration int `json:"duration"`
}

// BaseData is a struct for reading data.json
type BaseData struct {
	Types    []Type    `json:"types"`
	Pokemons []Pokemon `json:"pokemons"`
	Moves    []Move    `json:"moves"`

	/*Types []struct {
		// Name of the type
		Name string `json:"name"`
		// The effective basedata, damage multiplize 2x
		EffectiveAgainst []string `json:"effectiveAgainst"`
		// The weak basedata that against, damage multiplize 0.5x
		WeakAgainst []string `json:"weakAgainst"`
	} `json:"types"`
	Pokemons []struct {
		Number         string   `json:"Number"`
		Name           string   `json:"Name"`
		Classification string   `json:"Classification"`
		TypeI          []string `json:"Type I"`
		TypeII         []string `json:"Type II,omitempty"`
		Weaknesses     []string `json:"Weaknesses"`
		FastAttackS    []string `json:"Fast Attack(s)"`
		Weight         string   `json:"Weight"`
		Height         string   `json:"Height"`
		Candy          struct {
			Name     string `json:"Name"`
			FamilyID int    `json:"FamilyID"`
		} `json:"Candy"`
		NextEvolutionRequirements struct {
			Amount int    `json:"Amount"`
			Family int    `json:"Family"`
			Name   string `json:"Name"`
		} `json:"Next Evolution Requirements,omitempty"`
		NextEvolutions []struct {
			Number string `json:"Number"`
			Name   string `json:"Name"`
		} `json:"Next evolution(s),omitempty"`
		PreviousEvolutions []struct {
			Number string `json:"Number"`
			Name   string `json:"Name"`
		} `json:"Previous evolution(s),omitempty"`
		SpecialAttacks      []string `json:"Special Attack(s)"`
		BaseAttack          int      `json:"BaseAttack"`
		BaseDefense         int      `json:"BaseDefense"`
		BaseStamina         int      `json:"BaseStamina"`
		CaptureRate         float64  `json:"CaptureRate"`
		FleeRate            float64  `json:"FleeRate"`
		BuddyDistanceNeeded int      `json:"BuddyDistanceNeeded"`
	} `json:"pokemons"`
	Moves []struct {
		// The ID of the move
		ID int `json:"id"`
		// Name of the attack
		Name string `json:"name"`
		// Type of attack
		Type string `json:"type"`
		// The damage that enemy will take
		Damage int `json:"damage"`
		// Energy requirement of the attack
		Energy int `json:"energy"`
		// Dps is Damage Per Second
		Dps float64 `json:"dps"`
		// The duration
		Duration int `json:"duration"`
	} `json:"moves"`*/
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/list url:", r.URL)
	fmt.Fprint(w, "The List Handler\n")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/get url:", r.URL)
	fmt.Fprint(w, "The Get Handler\n")
}

func typeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/types url:", r.URL)
	fmt.Fprint(w, "All of the pokemon types\n")
	bd := getData()
	//fmt.Fprint(w, bd.Types)
	parseMap(bd, w, r)
}
func pokemonHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/pokemons url:", r.URL)
	fmt.Fprint(w, "All pokemons\n")
	/*bd := getData()
	fmt.Fprint(w, bd.Pokemons)*/
}
func moveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/moves url:", r.URL)
	fmt.Fprint(w, "All moves\n")
	/*bd := getData()
	fmt.Fprint(w, bd.Moves)*/
}
func otherwise(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Pokedex\n")

}

func getData() map[string]interface{} {
	log.Println("getData called")
	content, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Print("Error:", err)
	}
	//var basedata BaseData
	m := make(map[string]interface{})
	//err = json.Unmarshal([]byte(content), &basedata)
	err = json.Unmarshal([]byte(content), &m)
	if err != nil {
		fmt.Print("Error:", err)
	}
	return m
}
func parseMap(aMap map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	for key, val := range aMap {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			fmt.Fprintln(w, key)
			parseMap(val.(map[string]interface{}), w, r)
		case []interface{}:
			fmt.Fprintln(w, key)
			parseArray(val.([]interface{}), w, r)
		default:
			fmt.Fprintln(w, key, ":", concreteVal)

		}
	}
}

func parseArray(anArray []interface{}, w http.ResponseWriter, r *http.Request) {
	for i, val := range anArray {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			fmt.Println("Index:", i)
			parseMap(val.(map[string]interface{}), w, r)
		case []interface{}:
			//fmt.Println("Index:", i)
			parseArray(val.([]interface{}), w, r)
		default:
			fmt.Fprintln(w, "-", concreteVal)
			//fmt.Fprintln(w, "Index", i, ":", concreteVal)

		}
	}
}

func main() {
	//TODO: read data.json to a BaseData

	getData()

	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/types", typeHandler)
	http.HandleFunc("/pokemons", pokemonHandler)
	http.HandleFunc("/moves", moveHandler)
	//TODO: add more
	http.HandleFunc("/", otherwise)
	log.Println("starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
