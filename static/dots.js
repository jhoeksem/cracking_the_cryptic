"use strict";

var previousClick = [-1, -1];
var dimensions = [4,4];
var turn = 0;
var turnColor = ["red", "blue"];
var lines =new Array(dimensions[0]*2 -1);
for (var i = 0; i < dimensions[0]*2 -1; i++){
    lines[i] = new Array(dimensions[1] - (1+i)%2);
    lines[i].fill("white");
}

var squares = new Array(dimensions[0] - 1);
for (var i = 0; i < dimensions[0] - 1; i++){
    squares[i] = new Array(dimensions[1] - 1);
    squares[i].fill("white-background");
}

createTable();

function createTable(){
    var turnLabel = document.getElementById("turn-label");
    turnLabel.textContent = turnColor[turn]+"'s Turn";
    var table = document.getElementById("game-board");
    var tableStructure = "";
    tableStructure += "<tbody>";
    for (var i =0; i < 2 * dimensions[0] -1; i++){
        var buttonColumn = Math.floor(i/2);
        tableStructure += "<tr>";
        if (i%2 ==0){
            for(var j =0; j < dimensions[1]; j++){
                var position = buttonColumn+"-"+j
                tableStructure += "<td><input type=\"radio\" id=\""+position+"\" onClick=\"radioButtonHandler("+buttonColumn+", "+j+")\"></td>";
                if(j < dimensions[1] - 1){
                    tableStructure += "<td><hr class=\"line horizontal-line "+lines[i][j]+"\"></td>";
                }
            }
        } else{
            for(var j =0; j < dimensions[1]; j++){
                tableStructure += "<td><hr class=\"line vertical-line "+lines[i][j]+"\"></td>";
                if(j < dimensions[1] - 1){
                    tableStructure += "<td><p class=\"box "+squares[buttonColumn][j]+"\"></td>";
                }
            }
        }
        tableStructure += "</tr>";
    }
    tableStructure += "</tbody>";
    table.innerHTML = tableStructure;
}

async function sendRequest(x, y){
    var move = [x, y];
    var data = {"Move": move, "Game": {"Squares": squares, "Lines": lines}};
    var temp = await fetch("/updateTurn", {
        method: "POST", 
        body: JSON.stringify(data)
      }).then(res => {
          return res;
      });
      var text = await temp.text();
      var object = JSON.parse(text);
      lines = object.Game.Lines;
      squares = object.Game.Squares;
      return object.GameOver;
}

function resetButton(i, j){
    var previousButton = document.getElementById(i+"-"+j);
    if(previousButton !== null){
        previousButton.checked=false;
    }
}

async function radioButtonHandler(i, j){
    var insertedX = -1;
    var insertedY = -1;
    if (previousClick[0] < 0 ){
        previousClick[0] = i;
        previousClick[1] = j;
        return;
    }
    var gameOver = "";
    if (previousClick[0] == i){
        if (j - previousClick[1] === 1){
            if (lines[2*i][previousClick[1]] === "white"){
                insertedX = 2*i;
                insertedY = previousClick[1];
                gameOver = await sendRequest(insertedX, insertedY);
            }
        }
        else if (j - previousClick[1] === -1){
            if (lines[2*i][j] === "white"){
                insertedX = 2*i;
                insertedY = j;
                gameOver = await sendRequest(insertedX, insertedY);
            }
        }
    }
    else if (previousClick[1] == j){
        if (i - previousClick[0] === 1){
            if(lines[1 + 2 * previousClick[0]][j] === "white"){
                insertedX = 1 + 2 * previousClick[0];
                insertedY = j;
                gameOver = await sendRequest(insertedX, insertedY);
            }
        }
        else if (i - previousClick[0] === -1){
            if(lines[1 + 2 * i][j] === "white"){
                insertedX = 1 + 2 * i;
                insertedY = j;
                gameOver = await sendRequest(insertedX, insertedY);
            }

        }
    }
    resetButton(previousClick[0], previousClick[1]);
    resetButton(i, j);
    previousClick[0] = -1;
    previousClick[1] = -1;
    createTable();
    if (gameOver !== ""){
        alert(gameOver);
    }
    return;
    
}