package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/nu7hatch/gouuid"
)

type GameState struct {
	players []Player
	player_to_move_ind int 
	score int 
	winners_name string 
}

type Player struct {
	name string 
}

var GAMES = make(map[string] *GameState)

// Put Request
func start_game(w http.ResponseWriter, r *http.Request) {
	log.Println("Started")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	uuid, err := uuid.NewV4()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	players_info := map[string]string {}

	err = json.Unmarshal(body, &players_info)
	if err != nil {
		panic(err)
	}

	player := Player{name: players_info["name"]}
	
	game_state := &GameState{score:21}
	game_state.player_to_move_ind = 0
	(*game_state).players = append((*game_state).players, player)
	uuid_string := uuid.String()

	GAMES[uuid_string] = game_state

	log.Printf("Player: %s, started game: %s", player.name, uuid_string)
	fmt.Fprint(w, uuid_string)
} 

// Get 
func join_game (w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	players_info := map[string]string {}

	err = json.Unmarshal(body, &players_info)
	if err != nil {
		panic(err)
	}

	uuid_string := players_info["uuid"]
	game_state, ok := GAMES[uuid_string]
	if !ok {
		log.Println("uuid not found")
		fmt.Fprint(w, "no game with give uuid found")
		return
	}

	if players_info["name"] == (*game_state).players[0].name {
		log.Printf("Player attempeted to join the game using the same name as the creator: %s\n", players_info["name"])
		fmt.Fprintf(w, "Name: %s is already used\n", players_info["name"])
		return 
	}

	if len((*game_state).players) == 2 {
		log.Printf("Player attempted to join a game with already two players")
		fmt.Fprintf(w, "Game uuid: %s already has two players\n", uuid_string)
		return 
	}

	player := Player{name: players_info["name"]}
	(*game_state).players = append((*game_state).players, player)

	log.Printf("Player: %s joined game: %s\n", player.name, uuid_string)
}

// Get 
func check_status (w http.ResponseWriter, r *http.Request) {
	
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	uuid_string := r.Form["uuid"][0]
	
	game_state, ok := GAMES[uuid_string]
	if ok {
		ind := (*game_state).player_to_move_ind
		fmt.Fprint(w, (*game_state).players[ind].name)
		fmt.Fprint(w, (*game_state).score)
	} else {
		fmt.Fprint(w, "No game with that uuid found")
	}
}

// Put
func make_move (w http.ResponseWriter, r *http.Request) {
	// TODO
	// Don't make a move when the game already
	// has a winner
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	move_info := map[string]string {}

	err = json.Unmarshal(body, &move_info)
	if err != nil {
		panic(err)
	}

	// TODO: handle the case where not all fields are provided

	player_name := move_info["name"]
	game_state := GAMES[move_info["uuid"]]
	move_str := move_info["move"]

	move_int, err := strconv.Atoi(move_str)
	if err != nil {
		log.Printf("Player: %s attempted invalid move: %s\n", player_name, move_str)
		fmt.Fprintf(w, "Can't convert: %s to int", move_str)
	}

	if move_int < 1 || move_int > 3 {
		log.Printf("Player: %s attempted invalid move: %s\n", player_name, move_str)
		fmt.Fprint(w, "A valid move is either 1, 2 or 3")
		return 
	}

	ind := (*game_state).player_to_move_ind
	if (*game_state).players[ind].name != player_name {
		log.Printf("Player: %s attempted move out of turn\n", player_name)
		fmt.Fprint(w, "This is not your move")
		return 
	}

	log.Printf("Move: Game: %s, Player: %s, Move %d\n", move_info["uuid"], player_name, move_int)

	(*game_state).score -= move_int 

	log.Printf("Game Score: %d\n", (*game_state).score)

	var other_player_ind int 

	if (*game_state).player_to_move_ind == 0 {
		other_player_ind = 1
	} else {
		other_player_ind = 0
	}

	if (*game_state).score <= 0 {
		(*game_state).winners_name = (*game_state).players[other_player_ind].name
		log.Printf("Player: %s has won the game\n", (*game_state).winners_name)
		return
	}

	(*game_state).player_to_move_ind = other_player_ind
}

func main() {
	// messenger := Messenger{}
	// messenger.message_map = map[string]string {}
	
	http.HandleFunc("/start", start_game)
	http.HandleFunc("/join", join_game)
	http.HandleFunc("/check_status", check_status)
	http.HandleFunc("/move", make_move)

	http.ListenAndServe("0.0.0.0:8090", nil)
}
