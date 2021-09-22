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
	"strconv"
)

type sessionData struct {
    Uid string
}

type shareData struct {
    Imgd string
}

type file struct {
    Name string
	Base64 string
	Uid string
}

var activeUid []string
var dateUid  []time.Time

func generateToken() string {
	//Insecure uid generator?
    b := make([]byte, 10)
	rand.Read(b)
    return hex.EncodeToString(b)
}

func fileCleanerWorker(){
	log.Println("Cleaning old files...")
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	
	for{
		deleted := 0
		files, _ := ioutil.ReadDir("uploads/")
		
		if(len(files) > 500){
			// this is trending rolf
			panic("Waga na wa Megumin! 💥")
		}
		
		for _, f := range files {
			file, _ := os.Stat("uploads/"+f.Name())
			
			//Check modification time (Also creation)
			if(now.Sub(file.ModTime()).Hours() > 12){
				 err := os.Remove("uploads/"+f.Name()) 
				  if err != nil {
					log.Println("File "+f.Name()+" cant be deleted!")
				  }	else {
					  deleted++
				  }
			}
			
		}
		log.Println("Deleted "+ strconv.Itoa(deleted) +" files!")
		time.Sleep(24 * time.Hour)
	}
}

func expireUid() {
	
	if(len(activeUid) > 300){
		// this is trending rolf
	}
	
	for i := 0; i < len(activeUid); i++ {
		//Check date whith array date correlation 
		if(math.Trunc(time.Now().Sub(dateUid[i]).Minutes()) > 2){
			activeUid=append(activeUid[:i], activeUid[i+1:]...)
			dateUid=append(dateUid[:i], dateUid[i+1:]...)
		}
	}
	
}


///////////////////////////// HandleFuncs /////////////////////////////

func homePage(w http.ResponseWriter, r *http.Request) {
	
	if r.URL.Path != "/" {
        log.Println("404")
        return
    }
	
	uid := generateToken()
	authStruct := sessionData{
        Uid:       uid,
    }
    parsedTemplate, _ := template.ParseFiles("views/index.html")
    parsedTemplate.Execute(w, authStruct)
	
	//Clean old UID
	expireUid()
	
	//Generate UID
	activeUid = append(activeUid, uid)
	dateUid = append(dateUid, time.Now())
}

func imageViewer(w http.ResponseWriter, r *http.Request) {
	
	id := r.URL.Query().Get("file")
	content, _ := ioutil.ReadFile("uploads/"+id+".data")
	rawdata:=string(content)
	imgStruct := shareData{
        Imgd:       rawdata,
    }
    parsedTemplate, _ := template.ParseFiles("views/share.html")
    parsedTemplate.Execute(w,imgStruct)
	
	
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
	  go func() {
			datafile, _ := os.Create("uploads/"+tempFile.Uid+".data")
			defer datafile.Close()
			datafile.WriteString(tempFile.Base64)
			datafile.Close()
	  }()
	}
}



///////////////////////////// Main /////////////////////////////
func main() {
	go fileCleanerWorker() //Launch worker for clean files every day
	
	http.HandleFunc("/", homePage)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/upload", uploader)
	http.HandleFunc("/share", imageViewer)
	
	log.Println("Server running...")
	log.Fatal(http.ListenAndServe(":9090", nil))	
}