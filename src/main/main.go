package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"reflect"

	"github.com/gorilla/mux"
)

type Type struct {
	// Name of the type
	Name string `json:"name"`
	// The effective basedata, damage multiplize 2x
	EffectiveAgainst []string `json:"effectiveAgainst"`
	// The weak basedata that against, damage multiplize 0.5x
	WeakAgainst []string `json:"weakAgainst"`
}

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
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/list url:", r.URL)
	fmt.Fprint(w, "The List Handler\n")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/get url:", r.URL)
	fmt.Fprint(w, "The Get Handler\n")
}

//Function to handle single Pokemon type request.
func returnSingleType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["type"]
	b := readData()
	for _, tip := range b.Types {
		if tip.Name == key {
			printType(tip, w, r)
		}
	}

	fmt.Println("Key: " + key)
}

//Function to handle single Pokemon request.
func returnSinglePokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	b := readData()
	for _, pokemon := range b.Pokemons {
		if pokemon.Name == key {
			printPokemon(pokemon, w, r)
		}
	}

	fmt.Println("Key: " + key)
}

//Function to handle single Pokemon Move request.
func returnSingleMove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	b := readData()
	for _, move := range b.Moves {
		if move.Name == key {
			printMove(move, w, r)
		}
	}

	fmt.Println("Key: " + key)
}

//Lists all of Pokemons.
func pokemonHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/pokemons url:", r.URL)
	fmt.Fprint(w, "All pokemons\n")

	b := readData()
	for _, pokemon := range b.Pokemons {
		fmt.Fprintln(w)
		printPokemon(pokemon, w, r)
	}

}

//Lists all of Pokemon Moves.
func moveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/moves url:", r.URL)
	fmt.Fprint(w, "All moves\n")

	b := readData()
	for _, move := range b.Moves {
		fmt.Fprintln(w)
		printMove(move, w, r)
	}

}

//Lists all of the Pokemon Types.
func typeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/types url:", r.URL)
	fmt.Fprint(w, "All of the pokemon types\n")
	b := readData()
	for _, tip := range b.Types {
		fmt.Fprintln(w)
		printType(tip, w, r)
	}
}
func test(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["type"]
	b := readData()
	for _, pokemon := range b.Pokemons {
		if pokemon.TypeI[0] == key /*|| pokemon.TypeII[0] == key */ {
			fmt.Fprintln(w)
			printPokemon(pokemon, w, r)
		}
	}

	fmt.Fprintln(w, key)
}

//Function for main page.Contains info about how to use.
func otherwise(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Pokedex\n")
	fmt.Fprintln(w, " /types for all of the Pokemon types.\n /moves for all of the Pokemon moves.\n /pokemons for all of the Pokemons")
}

//Function to print Pokemon Move.
func printMove(m Move, w http.ResponseWriter, r *http.Request) {
	s := reflect.ValueOf(&m).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Fprintln(w, typeOfT.Field(i).Name, f.Interface())
	}
}

//Function to print Pokemon.
func printPokemon(p Pokemon, w http.ResponseWriter, r *http.Request) {
	s := reflect.ValueOf(&p).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Fprintln(w, typeOfT.Field(i).Name, f.Interface())
	}
}

//Function to print Pokemon Type.
func printType(t Type, w http.ResponseWriter, r *http.Request) {
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Fprintln(w, typeOfT.Field(i).Name, f.Interface())
	}
}

//Function to use the data from JSON file.Returns BaseData.
func readData() BaseData {
	log.Println("getData called")
	//Reads data from data.json
	content, err := ioutil.ReadFile("data.json")
	//error handling
	if err != nil {
		fmt.Print("Error:", err)
	}
	var basedata BaseData
	//decoding JSON data into ByteData.
	err = json.Unmarshal([]byte(content), &basedata)
	//error handling
	if err != nil {
		fmt.Print("Error:", err)
	}
	return basedata
}

func main() {
	//TODO: read data.json to a BaseData
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/list", listHandler)
	myRouter.HandleFunc("/list/{type}", test)
	myRouter.HandleFunc("/get", getHandler)
	myRouter.HandleFunc("/types", typeHandler)
	myRouter.HandleFunc("/pokemons", pokemonHandler)
	myRouter.HandleFunc("/moves", moveHandler)
	myRouter.HandleFunc("/types/{type}", returnSingleType)
	myRouter.HandleFunc("/pokemons/{name}", returnSinglePokemon)
	myRouter.HandleFunc("/moves/{name}", returnSingleMove)
	//TODO: add more
	myRouter.HandleFunc("/", otherwise)
	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
