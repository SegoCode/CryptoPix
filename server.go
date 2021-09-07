package main

import (
	"html/template"
	"log"
	"net/http"
	"encoding/hex"
	"math/rand"
	"encoding/json"
	"time"
	"math"
)

type sessionData struct {
    Uid string
}

type file struct {
    Name string
	Base64 string
	Uid string
}

var activeUid []string
var dateUid  []time.Time

func generateToken() string {
    b := make([]byte, 10)
	rand.Read(b)
    return hex.EncodeToString(b)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	uid := generateToken()
	authStruct := sessionData{
        Uid:       uid,
    }
    parsedTemplate, _ := template.ParseFiles("index.html")
    parsedTemplate.Execute(w, authStruct)
	activeUid = append(activeUid, uid)
	dateUid = append(dateUid, time.Now())
	//time.Sleep(60 * time.Second)
	
	printSlice()
}

func printSlice() {
	for i, s := range activeUid {
		log.Println(s)
		log.Println(math.Trunc(time.Now().Sub(dateUid[i]).Minutes()))
	}
	
}


func uploader(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
    var t file
    decoder.Decode(&t)
    log.Println(t.Base64)
	log.Println(t.Name)
	log.Println(t.Uid)
}


func main() {
	http.HandleFunc("/", homePage)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/upload", uploader)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
