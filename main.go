package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var templates = template.Must(template.ParseFiles("templates/view.html", "templates/edit.html"))

var baseUrl = "http://api.weatherapi.com/v1/"

type Local struct {
	Location Location `json:"location"`
	Current Current `json:"current"`
}

type Location struct {
	Name string `json:"name"`
	Country string `json:"country"`
	LocalTime string `json:"localTime"`
}

type Current struct {
	Temp float32 `json:"temp_c"`
	IsDay int8 `json:"is_day"`
	Cloud int8 `json:"cloud"`
}

func getLocal(key string) (*Local, error){
	client := &http.Client{}
	
	req, err := http.NewRequest("GET", baseUrl + "current.json?q=auto%3Aip", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("key", key)

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	
	if err != nil {
		return nil, err
	}

	var localClime Local
	if err := json.Unmarshal(body, &localClime); err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	return &localClime, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, c *Local) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func localClimeHandler(w http.ResponseWriter, r *http.Request, key string) {
	local, err := getLocal(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	renderTemplate(w, "view", local)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := os.Getenv("API_CLIME_KEY")

		fn(w, r, key)
	}
}


func main()  {
	fmt.Println("Working just fine...")

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/view/local", makeHandler(localClimeHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
