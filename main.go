package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Trip struct {
	ID      int     `json:"id"`
	Destino string  `json:"destino"`
	Preco   float64 `json:"preco"`
	Duracao int     `json:"duracao"`
}

var trips []Trip

func init() {
	rand.Seed(time.Now().UnixNano())
	trips = []Trip{
		{ID: 1, Destino: "Paris", Preco: 4500.50, Duracao: 7},
		{ID: 2, Destino: "Roma", Preco: 3200.00, Duracao: 5},
	}
}

func getTrips(w http.ResponseWriter, r *http.Request) {
	setJSONHeader(w)
	json.NewEncoder(w).Encode(trips)
}

func getTripByDestino(w http.ResponseWriter, r *http.Request) {
	setJSONHeader(w)
	destino := r.URL.Query().Get("destino")
	if destino == "" {
		http.Error(w, "Destino não informado", http.StatusBadRequest)
		return
	}
	var resultados []Trip
	for _, trip := range trips {
		if strings.Contains(strings.ToLower(trip.Destino), strings.ToLower(destino)) {
			resultados = append(resultados, trip)
		}
	}
	json.NewEncoder(w).Encode(resultados)
}

func createTrip(w http.ResponseWriter, r *http.Request) {
	setJSONHeader(w)
	var newTrip Trip
	if err := json.NewDecoder(r.Body).Decode(&newTrip); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	newTrip.ID = gerarIDUnico()
	trips = append(trips, newTrip)
	json.NewEncoder(w).Encode(newTrip)
}

func updateTrip(w http.ResponseWriter, r *http.Request) {
	setJSONHeader(w)
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var updatedTrip Trip
	if err := json.NewDecoder(r.Body).Decode(&updatedTrip); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	for i, trip := range trips {
		if trip.ID == id {
			updatedTrip.ID = id
			trips[i] = updatedTrip
			json.NewEncoder(w).Encode(updatedTrip)
			return
		}
	}
	http.Error(w, "Viagem não encontrada", http.StatusNotFound)
}

func deleteTrip(w http.ResponseWriter, r *http.Request) {
	setJSONHeader(w)
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	for i, trip := range trips {
		if trip.ID == id {
			trips = append(trips[:i], trips[i+1:]...)
			fmt.Fprintf(w, "Viagem com ID %d removida\n", id)
			return
		}
	}
	http.Error(w, "Viagem não encontrada", http.StatusNotFound)
}

func gerarIDUnico() int {
	for {
		id := rand.Intn(1000) + 1
		existe := false
		for _, trip := range trips {
			if trip.ID == id {
				existe = true
				break
			}
		}
		if !existe {
			return id
		}
	}
}

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func serveFrontend(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao carregar página: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/trips", getTrips)
	http.HandleFunc("/trips/search", getTripByDestino)
	http.HandleFunc("/trips/create", createTrip)
	http.HandleFunc("/trips/update", updateTrip)
	http.HandleFunc("/trips/delete", deleteTrip)
	http.HandleFunc("/", serveFrontend)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("rodando no postman...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
