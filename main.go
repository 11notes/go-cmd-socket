package main

import (
	"log"
	"os"
	"fmt"
	"net"
	"net/http"
	"os/exec"
  "flag"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Command struct {
	Bin string `json:"bin"`
	Arguments []string `json:"arguments"`
	Environment map[string]interface{} `json:"environment"`
}

func main() {
  sockp := flag.String("s", "/run/cmd/cmd.sock", "path to socket file")
  flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", post).Methods("POST")
 
	srv := &http.Server{
	 Handler: r,
	}

	err := os.Remove(*sockp) 
	unix, err := net.Listen("unix", *sockp)
	if err != nil {
		log.Fatalf("could not open UNIX socket %v", err)
	}
	srv.Serve(unix)
}
 
func post(w http.ResponseWriter, r *http.Request) {
	var c Command
	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {  
		http.Error(w, err.Error(), http.StatusBadRequest)  
		return
	} 

	cmd := exec.Command(c.Bin, c.Arguments...)
	if(len(c.Environment) <= 0){
		cmd.Env = os.Environ()
	}else{
		env := append(os.Environ())
		for key, value := range c.Environment {
			env = append(env, fmt.Sprintf("%s=%v", key, value))
		}
		cmd.Env = env
	}
	data, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)  
		return
	}

	fmt.Fprintf(w, string(data))
}