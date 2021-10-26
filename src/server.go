package main

import (
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type sessionData struct {
	Uid string
}

type shareData struct {
	Imgd string
}

type fileData struct {
	Name   string `json:"Name"`
	Base64 string `json:"Base64"`
	Uid    string `json:"Uid"`
}

type Config struct {
	Server struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"server"`
	Files struct {
		MaxFileSize int `json:"max-file-size"`
		CleanTime   int `json:"clean-time"`
	} `json:"files"`
}

var config Config
var activeUid []string

///////////////////////////// Utils /////////////////////////////

//Load configuration file for server
func LoadConfiguration(file string) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal("Configuration File NotFound")
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}

func generateToken() string {
	//My own non-signed questionable JWT token
	//The token save in memory,so no one can make your own unexpired non-signed token
	now := time.Now()
	b := make([]byte, 10)
	rand.Read(b)
	return hex.EncodeToString(b) + "." + strconv.FormatInt(now.Unix(), 10) //Random token + timestamp
}

func fileCleanerWorker() {
	log.Println("Cleaning old files...")

	for { //loop every 24h
		deleted := 0
		files, _ := ioutil.ReadDir("uploads/")
		for _, f := range files {

			timeStirng := strings.Split(f.Name(), ".")[1] //Catch timestamp from file name
			timestamp, _ := strconv.ParseInt(timeStirng, 10, 64)
			timeFile := time.Unix(timestamp, 0) // Parse to golang date

			if math.Trunc(time.Now().Sub(timeFile).Hours()) >= float64(config.Files.CleanTime) { //Check expire
				err := os.Remove("uploads/" + f.Name())
				if err != nil {
					log.Println("File " + f.Name() + " cant be deleted!")
				} else {
					deleted++
				}
			}
		}
		log.Println("Deleted " + strconv.Itoa(deleted) + " files!")
		time.Sleep(24 * time.Hour)
	}
}

func expireUid() {
	for i := 0; i < len(activeUid); i++ {
		timeStirng := strings.Split(activeUid[i], ".")[1] //Catch timestamp from token

		timestamp, _ := strconv.ParseInt(timeStirng, 10, 64)
		timeToken := time.Unix(timestamp, 0) // Parse to golang date

		if math.Trunc(time.Now().Sub(timeToken).Minutes()) >= 2 { //Check expire
			activeUid = append(activeUid[:i], activeUid[i+1:]...)
		}
	}

}

///////////////////////////// HandleFuncs /////////////////////////////

func homePage(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		return
	}

	uid := generateToken()
	authStruct := sessionData{
		Uid: uid,
	}
	parsedTemplate, _ := template.ParseFiles("views/index.html")
	parsedTemplate.Execute(w, authStruct)

	//Clean old UID
	expireUid()

	//Generate UID
	activeUid = append(activeUid, uid)
}

func imageViewer(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("file")
	content, _ := ioutil.ReadFile("uploads/" + id + ".data")
	rawdata := string(content)
	imgStruct := shareData{
		Imgd: rawdata,
	}
	parsedTemplate, _ := template.ParseFiles("views/share.html")
	parsedTemplate.Execute(w, imgStruct)

}

func uploader(w http.ResponseWriter, r *http.Request) {
	// TODO create high load management
	// TODO show alert in web, total capacity
	// TODO Create control toomany request
	var tempFile fileData
	var existuid = false
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&tempFile)

	//Check if exist uid, this "prevent" use api out of web
	for _, auid := range activeUid {
		if auid == tempFile.Uid {
			existuid = true
			break
		}
	}

	//Generate file on server
	if existuid {
		go func() {
			filebytes := 3 * (len(tempFile.Base64) / 4)
			if filebytes < config.Files.MaxFileSize { //Check file size before creation, miss 2 or 1 byte from B64
				datafile, _ := os.Create("uploads/" + tempFile.Uid + ".data")
				defer datafile.Close()
				datafile.WriteString(tempFile.Base64)
			}
		}()
	}
}

///////////////////////////// Main /////////////////////////////
func main() {
	go fileCleanerWorker() //Launch worker for clean files every day

	//PAGES
	http.HandleFunc("/", homePage)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/share", imageViewer)

	//REST
	http.HandleFunc("/upload", uploader)

	LoadConfiguration("config.json")
	log.Println("Server running at " + config.Server.Host + ":" + config.Server.Port + "...")
	log.Fatal(http.ListenAndServe(config.Server.Host+":"+config.Server.Port, nil)) // Server listener
}
