package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "https://minesweeper-sebastian-manfria.herokuapp.com/minesweeper/api/game/"

//const url = "http://localhost:8080/minesweeper/api/game/"

type GameRequest struct {
	Username string `json:"username"`
	Rows     int    `json:"rows"`
	Columns  int    `json:"columns"`
	Mines    int    `json:"mines"`
}

type CellRequest struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

type NewGameDTO struct {
	Username string `json:"username"`
	GameID   string `json:"game_id"`
}

type GameDTO struct {
	Username    string   `json:"username"`
	GameID      string   `json:"game_id"`
	ElapsedTime string   `json:"elapsed_time"`
	Status      string   `json:"status"`
	Board       BoardDTO `json:"board"`
}

type BoardDTO struct {
	Rows    int       `json:"rows"`
	Columns int       `json:"columns"`
	Mines   int       `json:"mines"`
	Cells   []CellDTO `json:"modified_cells"`
}

type CellDTO struct {
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Value  string `json:"value"`
}

func main() {
	var gameID = ""
	var gameDTO GameDTO

	for {
		fmt.Println("Please select an option: ")
		fmt.Println("1) New Game")
		fmt.Println("2) Resume Game")
		var input int
		fmt.Scanln(&input)

		if input == 1 {
			fmt.Print("Username: ")
			var username string
			fmt.Scanln(&username)

			fmt.Print("Rows: ")
			var rows int
			fmt.Scanln(&rows)

			fmt.Print("Columns: ")
			var columns int
			fmt.Scanln(&columns)

			fmt.Print("Mines: ")
			var mines int
			fmt.Scanln(&mines)

			game := newGame(rows, columns, mines, username)
			gameID = game.GameID
			gameDTO = getGame(gameID)
			printGame(gameDTO)
		}

		if input == 2 {
			fmt.Println("Please enter game id: ")
			var gameId string
			fmt.Scanln(&gameId)
			gameDTO = getGame(gameId)
			gameID = gameDTO.GameID
			printGame(gameDTO)
		}
		if gameID != "" {
			break
		}

	}

	for {
		fmt.Println("Please select an option: ")
		fmt.Println("1) Reveal")
		fmt.Println("2) Flag")
		fmt.Println("3) Exit")
		var option int
		fmt.Scanln(&option)

		if option == 3 {
			break
		}

		fmt.Print("Row: ")
		var row int
		fmt.Scanln(&row)

		fmt.Print("Column: ")
		var column int
		fmt.Scanln(&column)

		if option == 1 {
			gameDTO = reveal(gameID, row, column)
			printGame(gameDTO)
		}
		if option == 2 {
			gameDTO = flag(gameID, row, column)
			printGame(gameDTO)
		}
		if gameDTO.Status == "LOST" {
			break
		}
	}

}

func printGame(gameDTO GameDTO) {
	s2 := toMatrix(gameDTO.Board)
	for i := 0; i < len(s2); i++ {
		fmt.Println(s2[i])
	}
	fmt.Println("id: " + gameDTO.GameID)
	fmt.Println("status: " + gameDTO.Status)
	fmt.Println("Elapsed Time: " + gameDTO.ElapsedTime)
}

func toMatrix(board BoardDTO) [][]string {
	numRows := board.Rows
	numColumns := board.Columns

	grid := make([][]string, numRows)
	for i := 0; i < numRows; i++ {
		grid[i] = make([]string, numColumns)
	}
	for _, c := range board.Cells {
		grid[c.Row][c.Column] = c.Value
	}

	for i := 0; i < numRows; i++ {
		for j := 0; j < numColumns; j++ {
			if grid[i][j] == "" {
				grid[i][j] = " |"
			} else {
				grid[i][j] = grid[i][j] + "|"

			}
		}
	}

	return grid
}

func reveal(gameID string, row int, column int) GameDTO {
	cReq := CellRequest{row, column}

	jsonReq, err := json.Marshal(cReq)

	req, err := http.NewRequest(http.MethodPut, url+gameID+"/reveal", bytes.NewBuffer(jsonReq))

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	// Convert response body to NewGameDTO
	var gameDTO GameDTO
	json.Unmarshal(bodyBytes, &gameDTO)
	return gameDTO

}

func flag(gameID string, row int, column int) GameDTO {
	cReq := CellRequest{row, column}

	jsonReq, err := json.Marshal(cReq)

	req, err := http.NewRequest(http.MethodPut, url+gameID+"/flag", bytes.NewBuffer(jsonReq))

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	// Convert response body to NewGameDTO
	var gameDTO GameDTO
	json.Unmarshal(bodyBytes, &gameDTO)
	return gameDTO

}

func getGame(id string) GameDTO {
	resp, err := http.Get(url + id)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	if err != nil {
		log.Fatalln(err)
	}

	// Convert response body to NewGameDTO
	var gameDTO GameDTO
	json.Unmarshal(bodyBytes, &gameDTO)
	return gameDTO
}

func newGame(rows int, columns int, mines int, username string) NewGameDTO {
	req := GameRequest{username, rows, columns, mines}

	jsonReq, err := json.Marshal(req)

	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	// Convert response body to NewGameDTO
	var newGameDTO NewGameDTO
	json.Unmarshal(bodyBytes, &newGameDTO)
	return newGameDTO
}
