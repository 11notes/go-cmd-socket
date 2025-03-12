package main

import (
	"os"
	"fmt"
	"net"
	"net/http"
	"os/exec"
  "flag"
	"encoding/json"
	"github.com/gorilla/mux"
)

const DEBUG = getEnv("DEBUG", false)

type Payload struct {
	Bin string `json:"bin,omitempty"`
	Args []string `json:"args,omitempty"`
}

func main() {
  sockp := flag.String("s", "/run/cmd/.sock", "path to socket file")
  flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", post).Methods("POST")
	r.HandleFunc("/", get).Methods("GET")
 
	srv := &http.Server{
	 Handler: r,
	}

	err := os.Remove(*sockp) 
	unix, err := net.Listen("unix", *sockp)
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

func getEnv(key, fallback string) string {
  value, exists := os.LookupEnv(key)
  if !exists {
      value = fallback
  }
  return value
}