package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const userStr = `{
	"Navbar": "Home",
	"Users": [
		{
			"FirstName": "John",
			"LastName": "Doe"
		},
		{
			"FirstName": "Jane",
			"LastName": "Doe"
		}
	],
	"Today": "Sunday"
}`

// A Response struct to map the Entire Response
type Response struct {
	Name    string    `json:"name"`
	Pokemon []Pokemon `json:"pokemon_entries"`
}

// A Pokemon Struct to map every pokemon to.
type Pokemon struct {
	EntryNo int            `json:"entry_number"`
	Species PokemonSpecies `json:"pokemon_species"`
}

// A struct to map our Pokemon's Species which includes it's name
type PokemonSpecies struct {
	Name string `json:"name"`
}

type pagename struct {
	Navbar string
	Name   string
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	faviconPath := filepath.Join("static", "favicon.ico")
	http.ServeFile(w, r, faviconPath)
}

func main() {
	p := pagename{Name: "Ralf", Navbar: "About"}
	response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")
	fmt.Println(p)

	fs := http.FileServer(http.Dir("static/css"))
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.Handle("/static/css/", http.StripPrefix("/static/css/", fs))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var user map[string]interface{}
	errorJson := json.Unmarshal([]byte(userStr), &user)
	if errorJson != nil {
		fmt.Println(errorJson)
	}

	basePath := filepath.Join("templates", "base.html")
	navbarPath := filepath.Join("templates", "navbar.html")
	helloWorldPath := filepath.Join("templates", "hello_world.html")

	t, err := template.ParseFiles(basePath, helloWorldPath, navbarPath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t.Name())

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		err = t.ExecuteTemplate(w, "base", p)
		if err != nil {
			fmt.Println(err)
		}
	})

	userPath := filepath.Join("templates", "user.html")
	t2 := template.Must(template.ParseFiles(basePath, userPath, navbarPath))
	fmt.Println(t2.Name())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = t2.ExecuteTemplate(w, "base", user)
		if err != nil {
			fmt.Println(err)
		}
	})

	pokePath := filepath.Join("templates", "pokemon_api.html")
	t3 := template.Must(template.ParseFiles(basePath, pokePath))
	fmt.Println(t3.Name())

	http.HandleFunc("/API", func(w http.ResponseWriter, r *http.Request) {
		responseData, err := io.ReadAll(response.Body)

		if err != nil {
			log.Fatal(err)
		}

		var responseObject Response
		json.Unmarshal(responseData, &responseObject)
		fmt.Println(responseObject.Name)
		fmt.Println(len(responseObject.Pokemon))

		err = t3.ExecuteTemplate(w, "base", responseObject.Pokemon)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.ListenAndServe(":3000", nil)
}
