
package main 


import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"time"
	"strconv"
)


type Note struct{
	Title string `json:"title"`
	Description string `json:"description"`
	CreatedOn time.Time `json:createdon`
}

var noteStore = make(map[string]Note)
var id int = 0

func PostNoteHandler(w http.ResponseWriter,r* http.Request){
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil{
		log.Fatal(err)
	}
	note.CreatedOn = time.Now()
	id++
	k := strconv.Itoa(id)
	noteStore[k] = note
	
	j,err := json.Marshal(note)
	
	if err != nil{
		log.Fatal(err)
	}
	
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}


func GetNoteHandler(w http.ResponseWriter,r* http.Request){
	var notes[] Note
	for _,note := range noteStore{
		notes = append(notes,note)
	}
	
	j,err := json.Marshal(notes)
	if err != nil{
		log.Fatal(err)
	}
	
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func UpdateNoteHandler(w http.ResponseWriter,r* http.Request){
	var err error
	var note Note
	
	vars := mux.Vars(r)
	k := vars["id"]
	
	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil{
		log.Fatal(err)
	}
	
	
	if result, ok := noteStore[k]; ok {
		note.CreatedOn = result.CreatedOn
		//delete existing item and add the updated item
		delete(noteStore, k)
		noteStore[k] = note
		} else {
		log.Printf("Could not find key of Note %s to delete", k)
		}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusNoContent)
	//w.Write({"msg":"Update was successfull",})
	
}

func DeleteNoteHandler(w http.ResponseWriter,r* http.Request){
	vars := mux.Vars(r)
	k := vars["id"]
	
	if _,ok := noteStore[k]; ok{
		delete(noteStore,k)
	} else{
		log.Printf("Could not find key of Note %s to delete", k)
	}
	
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusNoContent)
}

func main(){
	r := mux.NewRouter().StrictSlash(false)
	
	r.HandleFunc("/api/v1/note",PostNoteHandler).Methods("POST")
	r.HandleFunc("/api/v1/notes",GetNoteHandler).Methods("GET")
	r.HandleFunc("/api/v1/note/{id}",UpdateNoteHandler).Methods("PUT")
	r.HandleFunc("api/v1/note/{id}",DeleteNoteHandler).Methods("DELETE")
	
	server := &http.Server{
		Addr: ":8090",
		Handler: r,
		}
	
	log.Println("Listening...")
	server.ListenAndServe()
}
