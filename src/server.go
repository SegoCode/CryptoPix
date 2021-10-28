package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//Struct for html template
type sessionData struct {
	Uid string
}

//Struct for html template
type shareData struct {
	Imgd string
}

//Struct for delete file
type jwtTimeStamp struct {
	CreationDate int `json:"exp"`
}

//Struct for REST request
type fileData struct {
	Name   string `json:"Name"`
	Base64 string `json:"Base64"`
	Uid    string `json:"Uid"`
}

//Struct for local config
type configData struct {
	Server struct {
		Port string `json:"port"` //Server Port
		Host string `json:"host"` //Server ip
	} `json:"server"`
	Files struct {
		MaxFileSize int    `json:"max-file-size"` //Bytes
		CleanTime   int    `json:"clean-time"`    //hours
		SecretKey   string `json:"secret-key"`    //Secret key for JWT
	} `json:"files"`
}

var config configData

///////////////////////////// Utils /////////////////////////////

//Load configuration file for server
func LoadConfiguration(file string) {
	configFile, err := os.Open(file)

	if err != nil {
		log.Fatal("Configuration File NotFound")
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}

//JWT Token generator
func createJWTToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 5).Unix(), //Max token life
	})

	t, _ := token.SignedString([]byte(config.Files.SecretKey))

	return t
}

//JWT Token validator
func verifyToken(toc string) bool {
	token, err := jwt.Parse(toc, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Err")
		}
		return []byte(config.Files.SecretKey), nil
	})

	if err == nil && token.Valid {
		return true
	}

	return false
}

func fileCleanerWorker(timer *time.Ticker) {
	for range timer.C { //loop every tick
		log.Println("Cleaning old files...")
		deleted := 0
		files, _ := ioutil.ReadDir("uploads/")
		for _, f := range files {

			timeStirng := strings.Split(f.Name(), ".")[1]            //Catch timestamp from file name
			data, err := base64.StdEncoding.DecodeString(timeStirng) //Decoding
			if err == nil {
				tmps := jwtTimeStamp{}                             //Init Struct
				json.Unmarshal(data, &tmps)                        //Mapping json
				timeFile := time.Unix(int64(tmps.CreationDate), 0) //Parse to golang date
				elapsed := time.Since(timeFile)                    //Elapse time

				if math.Trunc(elapsed.Hours()) >= float64(config.Files.CleanTime) { //Check expire
					err := os.Remove("uploads/" + f.Name())
					if err != nil {
						log.Println("File " + f.Name() + " cant be deleted!")
					} else {
						deleted++
					}
				}
			}
		}
		log.Println("Deleted " + strconv.Itoa(deleted) + " files!")
	}
}

///////////////////////////// HandleFuncs /////////////////////////////

func homePage(w http.ResponseWriter, r *http.Request) {

	//Redirect to home page
	if r.URL.Path != "/" {
		log.Println("Try to access: " + r.URL.Path)
		http.ServeFile(w, r, "views/404.html")
	} else {
		//Struct for template
		uid := createJWTToken()
		authStruct := sessionData{
			Uid: uid,
		}
		//Create template
		parsedTemplate, _ := template.ParseFiles("views/index.html")
		//Send template
		parsedTemplate.Execute(w, authStruct)
	}
}

func imageViewer(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("file") //Get from url file id without fragment
	content, err := ioutil.ReadFile("uploads/" + id + ".data")

	if err != nil {
		//File dosent exist
		http.ServeFile(w, r, "views/404.html")
	} else {
		//Struct for template
		rawdata := string(content)
		imgStruct := shareData{
			Imgd: rawdata,
		}
		//Create template
		parsedTemplate, _ := template.ParseFiles("views/share.html")
		//Send template
		parsedTemplate.Execute(w, imgStruct)
	}
}

func uploader(w http.ResponseWriter, r *http.Request) {
	var tempFile fileData //Init Struct
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&tempFile)

	//Generate file on server asynchronous
	if verifyToken(tempFile.Uid) {
		go func() { //Asynchronous creation
			filebytes := 3 * (len(tempFile.Base64) / 4)
			if filebytes < config.Files.MaxFileSize { //Check file size before creation, miss 2 or 1 byte from B64
				datafile, _ := os.Create("uploads/" + tempFile.Uid + ".data")
				defer datafile.Close()
				datafile.WriteString(tempFile.Base64)
			}
		}()
	} else {
		w.WriteHeader(403) //Not valid JWT token
	}
}

///////////////////////////// Main /////////////////////////////
func main() {
	ticker := time.NewTicker(24 * time.Hour) //Set worker time
	defer ticker.Stop()
	go fileCleanerWorker(ticker)     //Launch worker for clean files every day
	LoadConfiguration("config.json") //Load config

	//PAGES
	http.HandleFunc("/", homePage)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/share", imageViewer)

	//REST
	http.HandleFunc("/upload", uploader)

	//SERVER
	log.Println("Server running at " + config.Server.Host + ":" + config.Server.Port + "...")
	log.Fatal(http.ListenAndServe(config.Server.Host+":"+config.Server.Port, nil)) // Server listener
	//log.Fatal(http.ListenAndServeTLS(config.Server.Host+":"+config.Server.Port, "full-cert.crt", "private-key.key", nil))

}
