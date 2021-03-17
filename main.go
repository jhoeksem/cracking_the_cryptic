package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Board struct {
	Lines   [][]string
	Squares [][]string
}

type ClientRequest struct {
	Game Board
	Move []int
}

type ServerResponse struct {
	Game     Board
	GameOver string
}

// Checks if move is valid and if it is it will update the line structure, use with go
func addLine(move []int, lines *[][]string, color string, result chan bool) {
	if (*lines)[move[0]][move[1]] == "white" {
		(*lines)[move[0]][move[1]] = color
		result <- true
	} else {
		result <- false
	}
}

// Checks if a square is filled in by the move and returns the valid in the channel, use with go
// There are two possible squares you could fill in, so use side = (0,1) to select between them
func checkSquare(move []int, lines *[][]string, squares *[][]string, side int, result chan bool) {
	var dimensions []int
	dimensions = append(dimensions, len((*lines)))
	dimensions = append(dimensions, 1+len((*lines)[0]))

	if move[0]%2 == 0 {
		if side == 0 {
			if move[0]-2 >= 0 && (*lines)[move[0]-2][move[1]] != "white" {
				if (*lines)[move[0]-1][move[1]] != "white" && (*lines)[move[0]-1][move[1]+1] != "white" {
					result <- true
					return
				}
			}
		} else if side == 1 {
			if move[0]+2 < dimensions[0] && (*lines)[move[0]+2][move[1]] != "white" {
				if (*lines)[move[0]+1][move[1]] != "white" && (*lines)[move[0]+1][move[1]+1] != "white" {
					result <- true
					return
				}
			}
		} else {
			panic("Invalid side entered into checkSquare")
		}
	} else {
		if side == 0 {
			if move[1]-1 >= 0 && (*lines)[move[0]][move[1]-1] != "white" {
				if (*lines)[move[0]+1][move[1]-1] != "white" && (*lines)[move[0]-1][move[1]-1] != "white" {
					result <- true
					return
				}
			}
		} else if side == 1 {
			if move[1]+1 < dimensions[1] && (*lines)[move[0]][move[1]+1] != "white" {
				if (*lines)[move[0]+1][move[1]] != "white" && (*lines)[move[0]-1][move[1]] != "white" {
					result <- true
					return
				}
			}
		} else {
			panic("Invalid side entered into checkSquare")
		}

	}
	result <- false
}

// This function actually updates the squares with the appropiate color and return result in the channel when finished
func addSquare(move []int, squares *[][]string, side int, color string, result chan bool) {
	if move[0]%2 == 0 {
		if side == 0 {
			squareColumn := move[0]/2 - 1
			(*squares)[squareColumn][move[1]] = color + "-background"
		} else if side == 1 {
			squareColumn := move[0] / 2
			(*squares)[squareColumn][move[1]] = color + "-background"

		} else {
			panic("Invalid side entered into addSquare")
		}
	} else {
		if side == 0 {
			squareColumn := move[0] / 2
			(*squares)[squareColumn][move[1]-1] = color + "-background"

		} else if side == 1 {
			squareColumn := move[0] / 2
			(*squares)[squareColumn][move[1]] = color + "-background"

		} else {
			panic("Invalid side entered addSquare")
		}
	}
	result <- true
}

// Checks if there are any white squares left and if not computes the winner
func gameOver(squares *[][]string) string {
	points := 0
	for i := 0; i < len(*squares); i++ {
		for j := 0; j < len((*squares)[i]); j++ {
			if (*squares)[i][j] == "white-background" {
				return ""
			} else if (*squares)[i][j] == "blue-background" {
				points += 1
			} else if (*squares)[i][j] == "red-background" {
				points -= 1
			}
		}
	}
	if points == 0 {
		return "Game Over: Tie Game!"
	} else if points > 0 {
		return "Game Over: You Won!"
	} else {
		return "Game Over: You Lost!"
	}
}

// Tries a move and if it is valid it will add the move to the board
// Move is relative to the lines on the board not dots
// First return is if it was a valid move
// Second return is a gameOver string, empty if it is not gameOver
// Third return is if you filled a square, so you know to go again
func moveHandler(move []int, board *Board, color string) (bool, string, bool) {
	isValid := make(chan bool)
	addFirstSquare := make(chan bool)
	addSecondSquare := make(chan bool)
	go addLine(move, &board.Lines, color, isValid)
	go checkSquare(move, &(*board).Lines, &(*board).Squares, 0, addFirstSquare)
	go checkSquare(move, &(*board).Lines, &(*board).Squares, 1, addSecondSquare)
	valid := <-isValid
	firstSquare := <-addFirstSquare
	secondSquare := <-addSecondSquare
	if valid {
		if firstSquare {
			go addSquare(move, &(*board).Squares, 0, color, addFirstSquare)
		}
		if secondSquare {
			go addSquare(move, &(*board).Squares, 1, color, addSecondSquare)
		}

		if firstSquare {
			<-addFirstSquare
		}
		if secondSquare {
			<-addSecondSquare
		}
	} else {
		return false, "", true
	}
	gameOver := gameOver(&(*board).Squares)
	if gameOver != "" {
		return true, gameOver, true
	}
	if firstSquare || secondSquare {
		return true, "", true
	}
	return true, "", false
}

// Dummy func to choose the first unused line for opponent
// Replace this with min max
func makeMove(lines [][]string) []int {
	result := make([]int, 2)
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines[i]); j++ {
			if lines[i][j] == "white" {
				result[0] = i
				result[1] = j
				return result
			}
		}
	}
	panic("No valid moves should not be here")
}

//This is a function which gets a move and plays it onto the board
func opponentsTurn(board *Board) string {
	myTurn := true
	gameOver := ""
	valid := true
	for !valid || (myTurn && gameOver == "") {
		move := makeMove((*board).Lines)
		valid, gameOver, myTurn = moveHandler(move, board, "red")
	}
	return gameOver
}

// Basic function to handle playing the players move and calling oppenents turn if neccessary
func playPlayersTurn(move []int, board *Board, color string) string {
	valid, gameOver, myTurn := moveHandler(move, board, color)
	if gameOver != "" {
		return gameOver
	} else if valid && !myTurn {
		gameOver = opponentsTurn(board)
	}
	return gameOver
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/updateTurn", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal((err))
		}
		var structuredBody ClientRequest

		json.Unmarshal([]byte(body), &structuredBody)
		gameOver := playPlayersTurn(structuredBody.Move, &structuredBody.Game, "blue")
		var responseObject ServerResponse = ServerResponse{structuredBody.Game, gameOver}
		response, _ := json.Marshal(&responseObject)
		fmt.Fprintf(w, string(response))
	})
	fmt.Println("Running Game on Port 8080.")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
