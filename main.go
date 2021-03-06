package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
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

// returns game score, helpful for AI
func evaluation(squares *[][]string) int {
	points := 0
	for i := 0; i < len(*squares); i++ {
		for j := 0; j < len((*squares)[i]); j++ {
			if (*squares)[i][j] == "blue-background" {
				points += 1
			} else if (*squares)[i][j] == "red-background" {
				points -= 1
			}
		}
	}
	return points
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

func getLegalMoves(board Board) [][]int {

	moveCount := 0
	legalMoves := make([][]int, 200)
	for i := range legalMoves {
		legalMoves[i] = make([]int, 2)
		legalMoves[i][0] = -1
		legalMoves[i][1] = -1
	}
	for i := 0; i < len(board.Lines); i++ {
		for j := 0; j < len(board.Lines[i]); j++ {
			if board.Lines[i][j] == "white" {
				legalMoves[moveCount][0] = i
				legalMoves[moveCount][1] = j
				moveCount += 1
			}
		}
	}

	return legalMoves
}

// original legal moves array has length 1000 with lots of placeholder values
// this function trims off the placeholders
func trim_moves(moves_untrimmed [][]int) [][]int {

	// count number of real moves
	num_real_moves := 0
	for i := range moves_untrimmed {
		if moves_untrimmed[i][0] > -1 {
			num_real_moves += 1
		}
	}

	// copy real moves into new move array
	moves_trimmed := make([][]int, num_real_moves)
	for i := 0; i < num_real_moves; i++ {
		moves_trimmed[i] = moves_untrimmed[i]
	}

	return moves_trimmed

}

// directly copies string structs for lines and board
// otherwise new object copies by reference
func deepCopy(board Board) Board {
	new_board := board
	new_board.Lines = make([][]string, len(board.Lines))
	for i := range board.Lines {
		new_board.Lines[i] = make([]string, len(board.Lines[i]))
		copy(new_board.Lines[i], board.Lines[i])
	}
	new_board.Squares = make([][]string, len(board.Squares))
	for i := range board.Squares {
		new_board.Squares[i] = make([]string, len(board.Squares[i]))
		copy(new_board.Squares[i], board.Squares[i])
	}
	return new_board
}

// returns min integer value of array
func return_min(scores []int) int {
	min_score := 1000
	for i := range scores {
		if scores[i] < min_score {
			min_score = scores[i]
		}
	}
	return min_score
}

// returns max integer value of array
func return_max(scores []int) int {
	max_score := -1000
	for i := range scores {
		if scores[i] > max_score {
			max_score = scores[i]
		}
	}
	return max_score
}

// randomly select move with lowest expected score
// server AI wants low score
func return_random_best_move(legalMoves [][]int, scores []int, min_score int) []int {
	index := rand.Intn(len(legalMoves))
	for scores[index] != min_score {
		index = rand.Intn(len(legalMoves))
	}
	return legalMoves[index]
}

// creates array of size length and fills each index with base value 
func init_scores(base_value int, length int) []int {
	scores := make([]int, length)
	for i := 0; i < length; i++ {
		scores[i] = base_value
	}
	return scores
}

// returns player color string based on integer value
// player value is -1 for server, 1 for player
func get_color(player int) string {
	if player == 1 {
		return "blue"
	}
	if player == -1 {
		return "red"
	}

	panic("get_color() not working properly!")
}

// recursive function to return the value of a given move using minmax and a given depth 
// player value is -1 for server, 1 for player
func scoring_worker(board Board, move_index int, parent_scores *[]int, depth int, player int, parent_thread *sync.WaitGroup) {

	// lil debugging and thread chaos
	defer parent_thread.Done()
	var current_thread sync.WaitGroup
	//fmt.Println("depth: ", depth, " worker_ID: ", move_index)

	// get legal moves and # legal moves
	legalMoves := getLegalMoves(board)
	legalMoves = trim_moves(legalMoves)
	numLegalMoves := len(legalMoves)

	////////////////////////////////////////
	// base cases
	////////////////////////////////////////
	// 1. search depth is reached
	// 2. there are no legal moves
	if ( depth == 0 || numLegalMoves == 0 ){
		// next line almost definitely going to break
		(*parent_scores)[move_index] = evaluation(&(board).Squares)
		return
	}

	////////////////////////////////////////
	// recursive case
	////////////////////////////////////////
	// search to another depth and minmax scores
	myTurn := true
	//gameOver := ""
	//valid := true
	scores := init_scores(0, numLegalMoves)

	// spin up thread to calculate each legal move
	for index, move := range legalMoves {
		// prepare current thread for child thread
		current_thread.Add(1)

		// make move
		new_board := deepCopy(board)
		color := get_color(player)
		_, _, myTurn = moveHandler(move, &new_board, color)

		// create child thread
		if myTurn == true {
			go scoring_worker(new_board, index, &scores, depth - 1, player, &current_thread)
		} else {
			go scoring_worker(new_board, index, &scores, depth - 1, player * -1, &current_thread)
		}

	}
	current_thread.Wait()
	//fmt.Println(scores)

	// minmax and return value
	if player == 1 {
		(*parent_scores)[move_index] = return_max(scores)
	} else if player == -1{
		(*parent_scores)[move_index] = return_min(scores)
	} else {
		panic("player integer in scoring_worker wrong")
	}

	return
}

// simple AI min max
func makeMove(board Board) []int {

	// init variables
	depth := 3
	best_move := make([]int, 2)
	legalMoves := getLegalMoves(board)
	legalMoves = trim_moves(legalMoves)
	numLegalMoves := len(legalMoves)
	myTurn := true
	//gameOver := ""
	//valid := true
	scores := init_scores(0, numLegalMoves)

	// start waitgroup
	var wg sync.WaitGroup

	// search all moves
	for index, move := range legalMoves {
		// prepare current thread for child thread
		wg.Add(1)

		// dev update
		fmt.Println("searching: ", index)

		// make move
		new_board := deepCopy(board)
		_, _, myTurn = moveHandler(move, &new_board, "red")

		// create child thread
		if myTurn == true {
			go scoring_worker(new_board, index, &scores, depth - 1, -1, &wg)
		} else {
			go scoring_worker(new_board, index, &scores, depth - 1, 1, &wg)
		}
	}
	wg.Wait()
	//fmt.Println(scores)

	// select move
	min_score := return_min(scores)
	best_move = return_random_best_move(legalMoves, scores, min_score)

	return best_move
}

//This is a function which gets a move and plays it onto the board
func opponentsTurn(board *Board) string {
	myTurn := true
	gameOver := ""
	valid := true
	for !valid || (myTurn && gameOver == "") {
		move := makeMove(*board)
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

func errorCheckBoard(board *Board) {
	firstRowLen := len((*board).Lines[0])
	// for all of the rows go through and make sure that every other row is one more than the previous. Using odd and even indexes if something goes wrong the defer function should be ran
	for i := 0; i < len((*board).Lines); i++ {
		if (i%2 == 1) && (len((*board).Lines[i]) != firstRowLen+1) {
			fmt.Println("Board dimensions are off")
			panic("Board dimensions are off")
		} else if i%2 == 0 && len((*board).Lines[i]) != firstRowLen {
			fmt.Println("Board dimensions are off")
			panic("Board dimensions are off")

		}

	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/updateTurn", func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				gameOver := "error"
				var board Board
				var responseObject ServerResponse = ServerResponse{board, gameOver}
				response, _ := json.Marshal(&responseObject)
				fmt.Fprintf(w, string(response))
			}
		}()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal((err))
		}
		var structuredBody ClientRequest

		json.Unmarshal([]byte(body), &structuredBody)

		// Once have access to the board check its dimensions if something goes wrong the defer function will run
		errorCheckBoard(&structuredBody.Game)

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

// send back empty board if game corrupted. set game over to somethig send back empty board h=
