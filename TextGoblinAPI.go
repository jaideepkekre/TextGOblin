package textgoblin

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	

	"github.com/gorilla/mux"
)

func RunServer() {
	router := mux.NewRouter()
	router.HandleFunc("/classify", processData).Methods("POST")
	log.Fatal(http.ListenAndServe(":8001", router))
}

type message struct {
	Query              string
	StandardCategories []string
}

type transaction struct {
	Request  message
	Response string
	ID       int
}

func processData(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	println("Request Recieved!")
	req := message{}
	json.Unmarshal([]byte(body), &req)
	tr := transaction{}
	tr.Request = req
	tr.Response = tr.getMatches(req.Query, req.StandardCategories)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)	
	println("Response Created!")
	w.Write([]byte(tr.Response))

}

//GetMatches returns the scores for each standard category
func (tr transaction) getMatches(inputItem string, standardItems []string) string {

	mp := Match{}
	return mp.MatchProcessor(inputItem, standardItems)

}
