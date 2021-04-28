package app

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

type App struct {
	Router *mux.Router
}

func (a *App) Init() {
	log.Println("Initializing server...")
	a.Router = mux.NewRouter()
	a.setRoutes()
}

func (a *App) setRoutes() {
	a.get("/api/coinData", a.getData)
}

func (a *App) get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	apiKey := os.Getenv("API_KEY")

	client := &http.Client{}
	
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	q := url.Values{}
  	q.Add("start", "1")
  	q.Add("limit", "10")
  	q.Add("convert", "USD")

 	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)
	req.URL.RawQuery = q.Encode()
	
	resp, err := client.Do(req)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
		
	RespondJSON(w, http.StatusOK, string(respBody))
}

func (a *App) Run(host string) {
	log.Println("Server running at", host)
	log.Fatal(http.ListenAndServe(host, a.Router))
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}
