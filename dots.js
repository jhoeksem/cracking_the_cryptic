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

var points = [0, 0];
var TOTAL_POINTS = (dimensions[0] -1) * (dimensions[1] -1);
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

function resetButton(i, j){
    var previousButton = document.getElementById(i+"-"+j);
    if(previousButton !== null){
        previousButton.checked=false;
    }
}

function gameOverHandler(){
    if (points[0] + points[1] === TOTAL_POINTS){
        if(points[0]> points[1]){
            alert(turnColor[0]+" Won");
        } else if (points[9] < points){
            alert(turnColor[1]+" Won");
        } else{
            alert("Tie Game");
        }
    }
}

function checkSquares(x, y){
    if(lines[x][y] === "white"){
        console.log("Error somehow checking a line that was not selected");
        return;
    }
    var squareFilled = false;
    if (x % 2 ==0){
        if (x - 2 >= 0 && lines[x-2][y] !== "white"){
            if (lines[x-1][y] !== "white" && lines[x-1][y+1]!== "white"){
                var squareColumn = Math.floor(x/2) - 1;
                squares[squareColumn][y] = turnColor[turn]+"-background";
                squareFilled = true;
                points[turn] ++;
            }
        }
        if (x + 2 < 2 * dimensions[0] && lines[x+2][y] !== "white"){
            if (lines[x+1][y] !== "white" && lines[x+1][y+1]!== "white"){
                squares[Math.floor(x/2)][y] = turnColor[turn]+"-background";
                squareFilled = true;
                points[turn] ++;
            }
        }
    } else{
        if (y - 1 >= 0 && lines[x][y-1] !== "white"){
            if (lines[x+1][y-1] !== "white" && lines[x-1][y-1]!== "white"){
                squares[Math.floor(x/2)][y-1] = turnColor[turn]+"-background";
                squareFilled = true;
                points[turn] ++;
            }
        }
        if (y + 1 < dimensions[1] && lines[x][y+1] !== "white"){
            if (lines[x+1][y] !== "white" && lines[x-1][y]!== "white"){
                squares[Math.floor(x/2)][y] = turnColor[turn]+"-background";
                squareFilled = true;
                points[turn] ++;
            }
        }

    }
    if(!squareFilled){
        turn = (1 + turn) % 2;
    }

    gameOverHandler();
}

function radioButtonHandler(i, j){
    var validMove = false;
    var insertedX = -1;
    var insertedY = -1;
    if (previousClick[0] < 0 ){
        previousClick[0] = i;
        previousClick[1] = j;
        return;
    }
    if (previousClick[0] == i){
        if (j - previousClick[1] === 1){
            if (lines[2*i][previousClick[1]] === "white"){
                insertedX = 2*i;
                insertedY = previousClick[1];
                lines[insertedX][insertedY] = turnColor[turn];
            }
        }
        else if (j - previousClick[1] === -1){
            if (lines[2*i][j] === "white"){
                insertedX = 2*i;
                insertedY = j;
                lines[insertedX][insertedY] = turnColor[turn];
            }
        }
    }
    else if (previousClick[1] == j){
        if (i - previousClick[0] === 1){
            if(lines[1 + 2 * previousClick[0]][j] === "white"){
                insertedX = 1 + 2 * previousClick[0];
                insertedY = j;
                lines[insertedX][insertedY] = turnColor[turn];
            }
        }
        else if (i - previousClick[0] === -1){
            if(lines[1 + 2 * i][j] === "white"){
                insertedX = 1 + 2 * i;
                insertedY = j;
                lines[insertedX][insertedY] = turnColor[turn];
            }

        }
    }
    resetButton(previousClick[0], previousClick[1]);
    resetButton(i, j);
    previousClick[0] = -1;
    previousClick[1] = -1;
    if (insertedX >= 0){
        checkSquares(insertedX, insertedY);
    }
    createTable();
    return;
    
}