package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io"
	"io/ioutil"
)

type Messenger struct {
	message_map map[string]string
}

// Put Request
func (m *Messenger) send_message(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	new_message := map[string]string {}

	err = json.Unmarshal(body, &new_message)
	if err != nil {
		panic(err)
	}

	for k, v := range new_message {
		m.message_map[k] = v
	}
}

// Get Request
func (m *Messenger) get_message(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	message_key := r.Form["message_key"][0]
	
	message, ok := m.message_map[message_key]
	if ok {
		fmt.Fprint(w, message)
	} else {
		fmt.Fprint(w, "No message with that key found")
	}
}


func main() {
	messenger := Messenger{}
	messenger.message_map = map[string]string {}
	
	http.HandleFunc("/send_message", messenger.send_message)
	http.HandleFunc("/get_message", messenger.get_message)
	http.ListenAndServe("0.0.0.0:8090", nil)
}
