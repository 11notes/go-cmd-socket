package main

import (
	"os"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"encoding/json"
	"github.com/gorilla/mux"
)

const Socket = "/run/cmd.sock"

type Payload struct {
	Bin string `json:"bin,omitempty"`
	Args []string `json:"args,omitempty"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", post).Methods("POST")
	r.HandleFunc("/", get).Methods("GET")
 
	srv := &http.Server{
	 Handler: r,
	}

	err := os.Remove(Socket) 
	unix, err := net.Listen("unix", Socket)
	if err != nil {
		panic(err)
	}
	srv.Serve(unix)
}

func get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "post json with bin and args (array)")
}
 
func post(w http.ResponseWriter, r *http.Request) {
	var p Payload
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {  
		http.Error(w, err.Error(), http.StatusBadRequest)  
		return
	} 

	cmd := exec.Command(p.Bin, p.Args...)
	data, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)  
		return
	}

	fmt.Fprintf(w, string(data))
}