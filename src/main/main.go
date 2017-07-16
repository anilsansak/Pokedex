package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"reflect"

	"github.com/gorilla/mux"
)

//Type struct
type Type struct {
	// Name of the type
	Name string `json:"name"`
	// The effective basedata, damage multiplize 2x
	EffectiveAgainst []string `json:"effectiveAgainst"`
	// The weak basedata that against, damage multiplize 0.5x
	WeakAgainst []string `json:"weakAgainst"`
}

//Pokemon struct
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
	by       func(p1, p2 *Pokemon) bool
}

//Function to handle /list request. Contains information about instructions.
func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/list url:", r.URL)
	fmt.Fprintln(w, "To list Pokemons by their type, use /list/< insert type name here >. e.g: /list/Fire  ")
	fmt.Fprintln(w, "To sort Pokemons by their Height, use /list/<insert type here>/sortByHeight e.g: /list/Fire/sortByHeight ")
	fmt.Fprintln(w, "To sort Pokemons by their Weight, use /list/<insert type here>/sortByWeight e.g: /list/Fire/sortByWeight")
	fmt.Fprintln(w, "To sort Pokemons by their Base Attacks, use /list/<insert type here>/sortByBaseAttack e.g: /list/Fire/sortByBaseAttack")
	fmt.Fprintln(w, "To sort Pokemons by their Base Defenses, use /list/<insert type here>/sortByBaseDefense e.g: /list/Fire/sortByBaseDefense")
	fmt.Fprintln(w, "To sort Pokemons by their Base Stamina, use /list/<insert type here>/sortByBaseStamina e.g: /list/Fire/sortByBaseStamina")
	fmt.Fprintln(w, "To sort Pokemons by their Capture Rates, use /list/<insert type here>/sortByCaptureRate e.g: /list/Fire/sortByCaptureRate")
	fmt.Fprintln(w, "To sort Pokemons by their Flee Rates, use /list/<insert type here>/sortByFleeRate e.g: /list/Fire/sortByFleeRate")
	fmt.Fprintln(w, "To sort Pokemons by their Buddy Distances Needed, use /list/<insert type here>/sortByBuddyDistanceNeeded e.g: /list/Fire/sortByBuddyDistanceNeeded")
}

//Function to handle single Pokemon type request.
func returnSingleType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["type"]
	found := false
	b := readData()
	for _, tip := range b.Types {
		if tip.Name == key {
			printType(tip, w, r)
			found = true
		}
	}
	//Error handling.
	if found == false {
		fmt.Fprintln(w, "Please check your input and do not forget to start with an uppercase. e.g: /types/Fire")
	}

}

//Function to handle single Pokemon request.
func returnSinglePokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	found := false
	b := readData()
	for _, pokemon := range b.Pokemons {
		if pokemon.Name == key {
			printPokemon(pokemon, w, r)
			found = true
		}
	}
	//error handling
	if found == false {
		fmt.Fprintln(w, "Please check your input and do not forget to start with an uppercase. e.g: /pokemons/Pikachu")
	}

}

//Function to handle single Pokemon Move request.
func returnSingleMove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	found := false
	b := readData()
	for _, move := range b.Moves {
		if move.Name == key {
			printMove(move, w, r)
			found = true
		}
	}
	//Error handling.
	if found == false {
		fmt.Fprintln(w, "Please check your input and do not forget to start with an uppercase. e.g: /moves/Hyber Beam")
	}

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

//Lists pokemon by their type. Can also sort them by the input.
func listByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["type"]
	key2 := vars["sort"]
	//variables for input checking.
	keyfound := false
	nokey2 := false

	b := readData()

	//Definitons for sorting inputs.
	height := func(p1, p2 *Pokemon) bool {
		return p1.Height > p2.Height
	}
	weight := func(p1, p2 *Pokemon) bool {
		return p1.Weight > p2.Weight
	}
	baseattack := func(p1, p2 *Pokemon) bool {
		return p1.BaseAttack > p2.BaseAttack
	}
	basedefense := func(p1, p2 *Pokemon) bool {
		return p1.BaseDefense > p2.BaseDefense
	}
	basestamina := func(p1, p2 *Pokemon) bool {
		return p1.BaseStamina > p2.BaseStamina
	}
	capturerate := func(p1, p2 *Pokemon) bool {
		return p1.CaptureRate > p2.CaptureRate
	}
	fleerate := func(p1, p2 *Pokemon) bool {
		return p1.FleeRate > p2.FleeRate
	}
	buddydistanceneeded := func(p1, p2 *Pokemon) bool {
		return p1.BuddyDistanceNeeded > p2.BuddyDistanceNeeded
	}

	//Switch-case for sorting inputs.
	switch key2 {
	case "Height":
		By(height).Sort(b.Pokemons)

	case "Weight":
		By(weight).Sort(b.Pokemons)

	case "BaseAttack":
		By(baseattack).Sort(b.Pokemons)

	case "BaseDefense":
		By(basedefense).Sort(b.Pokemons)

	case "BaseStamina":
		By(basestamina).Sort(b.Pokemons)

	case "CaptureRate":
		By(capturerate).Sort(b.Pokemons)

	case "FleeRate":
		By(fleerate).Sort(b.Pokemons)

	case "BuddyDistanceNeeded":
		By(buddydistanceneeded).Sort(b.Pokemons)

	case "":

	default:
		//Error handling.
		fmt.Fprintln(w, " Please check your input. e.g: /list/Fire/sortByBaseAttack ")
		nokey2 = true
	}

	if nokey2 != true {
		for _, pokemon := range b.Pokemons {
			//TODO: check for TypeII
			if pokemon.TypeI[0] == key /*|| pokemon.TypeII[0] == key */ {
				fmt.Fprintln(w)
				keyfound = true
				printPokemon(pokemon, w, r)
			}

		}
		if keyfound == false {
			fmt.Fprintln(w, "Please check your input and do not forget to start with an uppercase. e.g: /list/Water")
		}
	}
}

//Function for main page.Contains info about how to use.
func otherwise(w http.ResponseWriter, r *http.Request) {
	//Prints instructions.
	fmt.Fprintln(w, "Welcome to Pokedex")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "/types for all of the Pokemon types.")
	fmt.Fprintln(w, "/types/< insert type here > to see information about given type. e.g: /types/Water")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "/moves for all of the Pokemon moves.")
	fmt.Fprintln(w, "/moves/< insert move here > to see information about given move. e.g: /moves/Flamethrower")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "/pokemons for all of the Pokemons.")
	fmt.Fprintln(w, "/pokemons/< insert pokemon here > to see information about given pokemon. e.g: /pokemons/Oddish")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "/list/< insert type here > to list Pokemons by the given type. e.g: /list/Fire")
	fmt.Fprintln(w, "/list/< insert type here >/sortBy<insert sorting criteria here > to list Pokemons by the given type and sort them by the given criteria.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "For more detailed information about /list , go to /list")

}

//Function to handle wrong path names.
func errorHandler(w http.ResponseWriter, r *http.Request) {
	//If input is not an accepted one,redirect to home page.
	//Redirecting to the home page.
	http.Redirect(w, r, "http://localhost:8080", 301)
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

//By is used for advanced sorting
type By func(p1, p2 *Pokemon) bool

//Sort is redefined for advanced sorting
func (by By) Sort(p []Pokemon) {
	ps := &BaseData{
		Pokemons: p,
		by:       by,
	}
	sort.Sort(ps)
}

//Len is a part of sort.Interface
func (slice *BaseData) Len() int {
	return len(slice.Pokemons)
}

//Less is a part of sort.Interface
func (slice *BaseData) Less(i, j int) bool {
	return slice.by(&slice.Pokemons[i], &slice.Pokemons[j])
}

//Swap is a part of sort.Interface
func (slice *BaseData) Swap(i, j int) {
	slice.Pokemons[i], slice.Pokemons[j] = slice.Pokemons[j], slice.Pokemons[i]
}

//Function to use the data from JSON file. Returns BaseData.
func readData() BaseData {

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

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/list", listHandler)
	myRouter.HandleFunc("/list/{type}", listByType)
	myRouter.HandleFunc("/list/{type}/sortBy{sort}", listByType)
	myRouter.HandleFunc("/list/{type}/{error}", errorHandler)
	myRouter.HandleFunc("/types", typeHandler)
	myRouter.HandleFunc("/pokemons", pokemonHandler)
	myRouter.HandleFunc("/moves", moveHandler)
	myRouter.HandleFunc("/types/{type}", returnSingleType)
	myRouter.HandleFunc("/pokemons/{name}", returnSinglePokemon)
	myRouter.HandleFunc("/moves/{name}", returnSingleMove)
	myRouter.HandleFunc("/{path}", errorHandler)
	myRouter.HandleFunc("/{err}/{error}", errorHandler)
	myRouter.HandleFunc("/", otherwise)

	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
