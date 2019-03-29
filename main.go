package main

import (
	"fmt"
	"math/rand"
	"time"
)


type Cell struct {
	row, column int
	alive, willBeAlive bool
}

type Row struct {
	cells []Cell;
}

type Map struct {
	rows []Row
	width, height int
}

func NewMap(width, height int) *Map {
	m := Map{
		rows:make([]Row, height),
		width:width,
		height:height,
	}
	for rowIndex := 0; rowIndex < height; rowIndex++ {
		m.rows[rowIndex] = Row{
			cells:make([]Cell, width),
		}
		for colIndex := 0; colIndex < width; colIndex++ {
			m.rows[rowIndex].cells[colIndex].row = rowIndex
			m.rows[rowIndex].cells[colIndex].column = colIndex
			m.rows[rowIndex].cells[colIndex].alive = rand.Float32()<.5
			m.rows[rowIndex].cells[colIndex].willBeAlive = false
		}
	}

	return &m
}

func Print(m Map) {
	msg := "\n";
	for colIndex := 0; colIndex < m.width; colIndex++ {
		msg += "--"
	}
	for rowIndex := 0; rowIndex < m.height; rowIndex++ {
		msg += "\n"
		for colIndex := 0; colIndex < m.width; colIndex++ {
			cell := m.rows[rowIndex].cells[colIndex];
			//msg += "|"
			if(cell.alive) {
				msg += "{}"
				//msg += "\u2588\u2588"
			} else {
				msg += "  "
			}
		}
		msg += "|"
	}
	fmt.Println("\x0c",msg)
}

func (c *Cell) Cycle(m Map){
	numberLiveNeighbors := c.CountLiveNeighbors(m)
	if(c.alive) {
		if (numberLiveNeighbors > 3 || numberLiveNeighbors <2) {
			c.willBeAlive = false
		}
	} else if (numberLiveNeighbors == 3) {
		c.willBeAlive = true
	}
}

func (c *Cell) Commit(){
	c.alive = c.willBeAlive
}

func (m *Map) Step()  {
	for rowIndex := 0; rowIndex < m.height; rowIndex++ {
		for colIndex := 0; colIndex < m.width; colIndex++ {
			m.rows[rowIndex].cells[colIndex].Cycle(*m)
		}
	}
	for rowIndex := 0; rowIndex < m.height; rowIndex++ {
		for colIndex := 0; colIndex < m.width; colIndex++ {
			m.rows[rowIndex].cells[colIndex].Commit()
		}
	}
}

func (c *Cell) CountLiveNeighbors(m Map) int {
	count := 0
	sameRow := c.row
	upRow := c.row - 1
	downRow := c.row + 1
	sameCol := c.column
	leftCol := c.column - 1
	rightCol := c.column + 1
	if(upRow < 0) {
		upRow = m.height-1
	}
	if(downRow >= m.height) {
		downRow = 0
	}
	if(leftCol < 0){
		leftCol = m.width - 1
	}
	if(rightCol >= m.width){
		rightCol = 0
	}
	//
	if(m.rows[upRow].cells[leftCol].alive){
		count++
	}
	if(m.rows[upRow].cells[sameCol].alive){
		count++
	}
	if(m.rows[upRow].cells[rightCol].alive){
		count++
	}
	//
	if(m.rows[sameRow].cells[leftCol].alive){
		count++
	}
	if(m.rows[sameRow].cells[rightCol].alive){
		count++
	}
	//
	if(m.rows[downRow].cells[leftCol].alive){
		count++
	}
	if(m.rows[downRow].cells[sameCol].alive){
		count++
	}
	if(m.rows[downRow].cells[rightCol].alive){
		count++
	}
	return count
}

//

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("start", rand.Float32())
	for x :=0; x < 1000; x++ {
		m := NewMap(40, 40)
		Print(*m)
		time.Sleep(time.Second * 2)
		for i :=0; i < 500; i++ {
			time.Sleep(time.Second/30)
			Print(*m)
			m.Step()
		}
	}

}