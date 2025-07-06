package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestGetTrips(t *testing.T) {
	req, _ := http.NewRequest("GET", "/trips", nil)
	res := httptest.NewRecorder()
	getTrips(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Esperado status 200, recebeu %d", res.Code)
	}
}

func TestCreateTrip(t *testing.T) {
	trip := Trip{Destino: "Lisboa", Preco: 2800.50, Duracao: 6}
	body, _ := json.Marshal(trip)
	req, _ := http.NewRequest("POST", "/trips/create", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	createTrip(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Esperado status 200, recebeu %d", res.Code)
	}
	var novaTrip Trip
	json.NewDecoder(res.Body).Decode(&novaTrip)
	if novaTrip.Destino != "Lisboa" {
		t.Errorf("Destino incorreto: %s", novaTrip.Destino)
	}
}

func TestDeleteTrip(t *testing.T) {
	trip := Trip{Destino: "Teste", Preco: 1000, Duracao: 3}
	trip.ID = gerarIDUnico()
	trips = append(trips, trip)

	req, _ := http.NewRequest("DELETE", "/trips/delete?id="+strconv.Itoa(trip.ID), nil)
	res := httptest.NewRecorder()
	deleteTrip(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Esperado status 200, recebeu %d", res.Code)
	}
	if strings.Contains(res.Body.String(), "Viagem com ID") == false {
		t.Errorf("Mensagem de sucesso não encontrada")
	}
}

func TestGetTripByDestino(t *testing.T) {
	req, _ := http.NewRequest("GET", "/trips/search?destino=Paris", nil)
	res := httptest.NewRecorder()
	getTripByDestino(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Esperado status 200, recebeu %d", res.Code)
	}
}
