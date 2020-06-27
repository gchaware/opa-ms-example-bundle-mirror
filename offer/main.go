package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type opaResult struct {
	Allow bool `json:"allow"`
}
type opaResponse struct {
	DecisionID string    `json:"decision_id"`
	Result     opaResult `json: result`
}

type inputJSONData struct {
	Method string `json:"method"`
	API    string `json:"api"`
	Jwt    string `json:"jwt"`
}

type opaRequest struct {
	Input inputJSONData `json:"input"`
}

type offer struct {
	OfferID  string `json:"offerid"`
	Title    string `json:"title"`
	Customer string `json:"customerid"`
	Segment  string `json:"segment"`
	Comments string `json:"notes"`
}

type allOffers = []offer

var offers = allOffers{
	{
		OfferID:  "1000",
		Title:    "New Office in LA",
		Customer: "1",
		Segment:  "LE",
		Comments: "Demo Offer # 1",
	},
}

func createOffer(w http.ResponseWriter, r *http.Request) {
	var allow = authorize(r)
	if allow {
		var newOffer offer
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "Invalid Format")
		}
		json.Unmarshal(reqBody, &newOffer)
		newOffer.OfferID = strconv.Itoa(rangeIn(1000, 10000))
		offers = append(offers, newOffer)
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newOffer)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
func getAllOffers(w http.ResponseWriter, r *http.Request) {
	var allow = authorize(r)
	if allow {
		json.NewEncoder(w).Encode(offers)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func getOffer(w http.ResponseWriter, r *http.Request) {
	var allow = authorize(r)
	if allow {
		offerID := mux.Vars(r)["id"]

		for _, singleOffer := range offers {
			if singleOffer.OfferID == offerID {
				json.NewEncoder(w).Encode(singleOffer)
			}
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func updateOffer(w http.ResponseWriter, r *http.Request) {
	var allow = authorize(r)
	if allow {
		offerID := mux.Vars(r)["id"]
		var updatedOffer offer

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "Invalid Format")
		}
		json.Unmarshal(reqBody, &updatedOffer)

		for i, singleOffer := range offers {
			if singleOffer.OfferID == offerID {
				singleOffer.Title = updatedOffer.Title
				singleOffer.Comments = updatedOffer.Comments
				offers = append(offers[:i], singleOffer)
				json.NewEncoder(w).Encode(singleOffer)
			}
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func deleteOffer(w http.ResponseWriter, r *http.Request) {
	var allow = authorize(r)
	if allow {
		offerID := mux.Vars(r)["id"]
		offerIndex := getIndex(offers, offerID)
		updatedList := make([]offer, 0)
		updatedList = append(updatedList, offers[:offerIndex]...)
		updatedList = append(updatedList, offers[offerIndex+1:]...)
		offers = updatedList
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func getIndex(offers []offer, offerID string) int {
	for i, singleOffer := range offers {
		if singleOffer.OfferID == offerID {
			return i
		}
	}
	return -1
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func authorize(r *http.Request) bool {
	var authHeader = r.Header.Get("Authorization")
	var apiEndPoint = r.URL.Path
	var reqMethod = r.Method
	var jwt = strings.Split(authHeader, "Bearer ")[1]

	varInputJSONData := &inputJSONData{API: apiEndPoint, Method: reqMethod, Jwt: jwt}
	varOpaRequest := &opaRequest{Input: *varInputJSONData}
	jsonValue, _ := json.Marshal(varOpaRequest)
	fmt.Println(string(jsonValue))
	response, err := http.Post("http://localhost:8181/v1/data/httpapi/authz", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("OPA request failed with error %s\n", err)
		return false
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))

		var res = new(opaResponse)
		err = json.Unmarshal(data, &res)
		if err != nil {
			fmt.Println("Error unmarshalling OPA response")
		}
		return res.Result.Allow
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/offer", createOffer).Methods("POST")
	router.HandleFunc("/offers", getAllOffers).Methods("GET")
	router.HandleFunc("/offer/{id}", getOffer).Methods("GET")
	router.HandleFunc("/offer/{id}", updateOffer).Methods("PATCH")
	router.HandleFunc("/offer/{id}", deleteOffer).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
