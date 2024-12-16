package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

var Lettreutilise []string

type tes struct {
	Title          string
	Tabvide        []rune
	Vie            int
	Lettreutilise  []string
	Input          string
	Lettrealeatoir []rune
	Motaleatoire   string
	Lettreafficher []string
	Fini           bool
	Finie          bool
}

var myStruct tes

func accueilHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("hangman-web-corentin323/accueil.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {

	if myStruct.Vie == 0 {
		myStruct = tes{
			Vie: 10,
		}
		myStruct.Motaleatoirdansletxt()
		myStruct.Sansdoublon()
		myStruct.Lettrealeatoire()
	}

	if r.Method == "POST" && !myStruct.Fini {
		r.ParseForm()
		myStruct.Input = r.FormValue("lettre")
		myStruct.Complete()
	}

	tmpl, err := template.ParseFiles("hangman-web-corentin323/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, myStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	http.HandleFunc("/", accueilHandler)

	http.HandleFunc("/jeu", Handler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Serveur démarré sur http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur :", err)
	}
}

func (t *tes) Motaleatoirdansletxt() {
	hase, ee := os.ReadFile("words.txt")
	if ee != nil {
		return
	}
	rest := strings.Split(string(hase), "\n")
	random := rand.Intn(len(rest))
	mot := rest[random]
	t.Motaleatoire = mot
}

func (t *tes) Sansdoublon() {
	var tabmot = []rune(t.Motaleatoire)
	booly := false
	for j := range tabmot {
		booly = false
		for i := range t.Tabvide {
			if t.Tabvide[i] == tabmot[j] {
				booly = true
			}
		}
		if !booly {
			t.Tabvide = append(t.Tabvide, tabmot[j])
		}
	}
}

func (t *tes) Lettrealeatoire() {
	var nom = len(t.Motaleatoire)/2 - 1
	var melangedelettre = []rune(t.Tabvide)

	rand.Shuffle(len(melangedelettre), func(i, j int) { melangedelettre[i], melangedelettre[j] = melangedelettre[j], melangedelettre[i] })
	for i := 0; i <= nom; i++ {
		t.Lettrealeatoir = append(t.Lettrealeatoir, melangedelettre[i])
	}

	runemot := []rune(t.Motaleatoire)
	strlettre := t.Lettrealeatoir
	runelettre := []rune(strlettre)

	Same := true
	for i := range runemot {
		Same = false
		for j := range runelettre {
			if runemot[i] == runelettre[j] {
				Same = true
				break
			}
		}
		if !Same {
			t.Lettreafficher = append(t.Lettreafficher, "_")
		}
		if Same {
			t.Lettreafficher = append(t.Lettreafficher, string(runemot[i]))
		}
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (t *tes) Complete() {
	atrouver := []rune(t.Motaleatoire)
	entree := []rune(strings.TrimSpace(t.Input))

	if contains(t.Lettreutilise, string(entree[0])) {
		return
	}

	Same := false
	for j := range atrouver {
		if entree[0] == atrouver[j] {
			t.Lettreafficher[j] = string(entree[0])
			Same = true
		}
	}

	if !Same {
		t.Vie--
	} else {
	}

	t.Lettreutilise = append(t.Lettreutilise, string(entree[0]))

	if t.Vie <= 0 {
		t.Fini = true
	}

	if strings.Join(t.Lettreafficher, "") == t.Motaleatoire {
		t.Finie = true
		t.Vie = 10
		t.Lettreutilise = nil
		t.Lettreafficher = nil
		t.Motaleatoire = ""
		t.Tabvide = nil
		t.Motaleatoirdansletxt()
		t.Sansdoublon()
		t.Lettrealeatoire()
		t.Fini = false
		t.Finie = false
	}
}
