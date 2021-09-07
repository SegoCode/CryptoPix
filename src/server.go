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
	"os"
	"io/ioutil"
	"fmt"
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
	//Insecure uid generator
    b := make([]byte, 10)
	rand.Read(b)
    return hex.EncodeToString(b)
}

func fileCleanerWorker(){
	log.Println("Cleaning old files...")
	for{
		files, _ := ioutil.ReadDir("uploads/")
		
		if(len(files) > 500){
			panic("Waga na wa Megumin! 💥")
		}
		
		for _, f := range files {
			file, _ := os.Stat("uploads/"+f.Name())
			modifiedtime := file.ModTime()
			log.Println("Last modified time : ", modifiedtime)	
		}
		time.Sleep(24 * time.Hour)
	}
}

func expireUid() {
	
	if(len(activeUid) > 300){
		// this is trending rolf
	}
	
	for i := 0; i < len(activeUid); i++ {
		if(math.Trunc(time.Now().Sub(dateUid[i]).Minutes()) > 2){
			activeUid=append(activeUid[:i], activeUid[i+1:]...)
			dateUid=append(dateUid[:i], dateUid[i+1:]...)
		}
	}
	
}


///////////////////////////// HandleFuncs /////////////////////////////

func homePage(w http.ResponseWriter, r *http.Request) {
	uid := generateToken()
	authStruct := sessionData{
        Uid:       uid,
    }
    parsedTemplate, _ := template.ParseFiles("index.html")
    parsedTemplate.Execute(w, authStruct)
	
	//Clean old UID
	expireUid()
	
	//Generate UID
	activeUid = append(activeUid, uid)
	dateUid = append(dateUid, time.Now())
}

func uploader(w http.ResponseWriter, r *http.Request) {
	//TODO Create control toomany request
	//TODO Set time out
	var tempFile file
	var existuid = false
	decoder := json.NewDecoder(r.Body)
    decoder.Decode(&tempFile)
	
	//Check if exist uid 
	for _, auid := range activeUid {
        if(auid == tempFile.Uid){
			existuid=true
            break
        }
    }
	
	//Generate file on server
	if(existuid){
		datafile, _ := os.Create("uploads/"+tempFile.Uid+".file")
		defer datafile.Close()
		datafile.WriteString(tempFile.Base64)
		datafile.Close()
	}
}

func imageViewer(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.URL.Query() {
        fmt.Printf("%s: %s\n", k, v)
    }
	
}



func main() {
	go fileCleanerWorker() //Launch worker for clean files every day
	
	http.HandleFunc("/", homePage)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/upload", uploader)
	http.HandleFunc("/share", imageViewer)
	
	log.Println("Server running...")
	log.Fatal(http.ListenAndServe(":9090", nil))	
}