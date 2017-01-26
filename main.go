package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"github.com/davecgh/go-spew/spew"
	"github.com/mattbaird/elastigo/lib"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

type ingredientDocument struct {
	Name       string `json:"name"`
	Suggestion inputs `json:"suggestion"`
}

type inputs struct {
	Input []string `json:"input"`
}

func main() {
	path := flag.String("path", "", "path to recipes")

	flag.Parse()

	dat, err := ioutil.ReadFile(*path)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	r := csv.NewReader(bytes.NewReader(dat))
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error reading all lines: %v", err)
	}

	ingredients := make(map[string]string)

	for _, line := range lines {
		ingredients[strings.ToLower(line[1])] = strings.ToLower(line[1])
	}

	connection := elastigo.NewConn()
	connection.Domain = "localhost"

	for _, ingredient := range ingredients {
		document := ingredientDocument{Name: ingredient}

		words := strings.Split(ingredient, " ")
		for i := 0; i < len(words); i++ {
			suggestion := strings.Join(words[i:len(words)], " ")
			document.Suggestion.Input = append(document.Suggestion.Input, suggestion)
		}

		_, err := connection.Index("recipes", "wlsm_ingredient", url.QueryEscape(ingredient), nil, document)
		if err != nil {
			spew.Dump(ingredient)
			spew.Dump(err)
		}
	}
}
