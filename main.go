package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"image"
	"image/color"
	"math/rand"
	"time"
)

type Cell struct {
	row, column        int
	alive, willBeAlive bool
}

type Row struct {
	cells []Cell
}

type Board struct {
	rows          []Row
	width, height int
}

func NewBoard(width, height int) *Board {
	m := Board{
		rows:   make([]Row, height),
		width:  width,
		height: height,
	}
	for rowIndex := 0; rowIndex < height; rowIndex++ {
		m.rows[rowIndex] = Row{
			cells: make([]Cell, width),
		}
		for colIndex := 0; colIndex < width; colIndex++ {
			m.rows[rowIndex].cells[colIndex].row = rowIndex
			m.rows[rowIndex].cells[colIndex].column = colIndex
			m.rows[rowIndex].cells[colIndex].alive = rand.Float32() < .5
			m.rows[rowIndex].cells[colIndex].willBeAlive = false
		}
	}
	return &m
}


func newGameView() *GameView {
	gv := &GameView{
		isVisible: true,
		board:     *NewBoard(30, 30),
	}
	return gv
}


func Print(m *Board) {
	msg := "\n";
	for colIndex := 0; colIndex < m.width; colIndex++ {
		msg += "--"
	}
	for rowIndex := 0; rowIndex < m.height; rowIndex++ {
		msg += "\n"
		for colIndex := 0; colIndex < m.width; colIndex++ {
			cell := m.rows[rowIndex].cells[colIndex];
			//msg += "|"
			if (cell.alive) {
				msg += "{}"
				//msg += "\u2588\u2588"
			} else {
				msg += "  "
			}
		}
		msg += "|"
	}
	fmt.Println("\x0c", msg)
}

func (c *Cell) Cycle(m Board) {
	numberLiveNeighbors := c.CountLiveNeighbors(m)
	if (c.alive) {
		if (numberLiveNeighbors > 3 || numberLiveNeighbors < 2) {
			c.willBeAlive = false
		}
	} else if (numberLiveNeighbors == 3) {
		c.willBeAlive = true
	}
}

func (c *Cell) Commit() {
	c.alive = c.willBeAlive
}

func (m *Board) Step() {
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

func (c *Cell) CountLiveNeighbors(m Board) int {
	count := 0
	sameRow := c.row
	upRow := c.row - 1
	downRow := c.row + 1
	sameCol := c.column
	leftCol := c.column - 1
	rightCol := c.column + 1
	if (upRow < 0) {
		upRow = m.height - 1
	}
	if (downRow >= m.height) {
		downRow = 0
	}
	if (leftCol < 0) {
		leftCol = m.width - 1
	}
	if (rightCol >= m.width) {
		rightCol = 0
	}
	//
	if (m.rows[upRow].cells[leftCol].alive) {
		count++
	}
	if (m.rows[upRow].cells[sameCol].alive) {
		count++
	}
	if (m.rows[upRow].cells[rightCol].alive) {
		count++
	}
	//
	if (m.rows[sameRow].cells[leftCol].alive) {
		count++
	}
	if (m.rows[sameRow].cells[rightCol].alive) {
		count++
	}
	//
	if (m.rows[downRow].cells[leftCol].alive) {
		count++
	}
	if (m.rows[downRow].cells[sameCol].alive) {
		count++
	}
	if (m.rows[downRow].cells[rightCol].alive) {
		count++
	}
	return count
}

//~~
type GameView struct {
	aliveColor color.Color
	deadColor  color.Color
	layoutSize fyne.Size
	position   fyne.Position
	isVisible  bool
	//
	board Board
	//$$
	render   *canvas.Raster
	objects  []fyne.CanvasObject
	imgCache *image.RGBA
}

func (g *GameView) Layout(size fyne.Size) {
	fmt.Println("layout", size.Width, size.Height)
	g.render.Resize(size)
}

func (g *GameView) Refresh() {
	canvas.Refresh(g.render)
}

func (g *GameView) ApplyTheme() {
	g.aliveColor = color.White
	g.deadColor = color.Black
}

func (g *GameView) BackgroundColor() color.Color {
	return color.Gray{}
}

func (g *GameView) Objects() []fyne.CanvasObject {
	return g.objects
}

func (g *GameView) Destroy() {
}

func (g *GameView) draw(w, h int) image.Image {
	img := g.imgCache
	if img == nil || img.Bounds().Size().X != w || img.Bounds().Size().Y != h {
		fmt.Println("Creating img", w, h)
		img = image.NewRGBA(image.Rect(0, 0, w, h))
		g.imgCache = img
	}
	//

	cellWidth := w / g.board.width
	cellHeight := h / g.board.height

	for rowIndex := 0; rowIndex < g.board.height; rowIndex++ {
		for colIndex := 0; colIndex < g.board.width; colIndex++ {
			x := colIndex * cellWidth
			y := rowIndex * cellHeight

			if (g.board.rows[rowIndex].cells[colIndex].alive) {
				drawRect(img, x, y, cellWidth, cellHeight, color.White)
			} else {
				drawRect(img, x, y, cellWidth, cellHeight, color.Black)
			}
		}
	}
	return img
}

func drawRect(img *image.RGBA, x int, y int, w int, h int, color color.Color) {
	for X := x; X < x+w; X++ {
		for Y := y; Y < y+h; Y++ {
			img.Set(X, Y, color)
		}
	}
}

func (g *GameView) animate() {
	go func() {
		tick := time.NewTicker(time.Second / 6)

		for {
			select {
			case <-tick.C:
				//if g.paused {
				//	continue
				//}

				//Print(*m)
				g.board.Step()
				widget.Refresh(g)
			}
		}
	}()
}

func (g *GameView) CreateRenderer() fyne.WidgetRenderer {
	renderer := g

	render := canvas.NewRaster(renderer.draw)
	renderer.render = render
	renderer.objects = []fyne.CanvasObject{render}
	renderer.ApplyTheme()

	return renderer
}

func (g *GameView) Hide() {
	g.isVisible = false
}

func (g *GameView) MinSize() fyne.Size {
	return fyne.NewSize(200, 200)
}

func (g *GameView) Move(pos fyne.Position) {
	g.position = pos
	widget.Renderer(g).Layout(g.layoutSize)
}

func (g *GameView) Position() fyne.Position {
	return g.position
}

func (g *GameView) Resize(size fyne.Size) {
	fmt.Println("resize", size.Width, size.Height)
	g.layoutSize = size
	widget.Renderer(g).Layout(size)
}

func (g *GameView) Show() {
	g.isVisible = true
}

func (g *GameView) Size() fyne.Size {
	return g.layoutSize
}

func (g *GameView) Visible() bool {
	return g.isVisible
}
func (g *GameView) typedRune(r rune) {
	if r == ' ' {
		//g.toggleRun()
	}
}

//~

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("start", rand.Float32())
	//

	app := app.New()
	gv := newGameView()
	window := app.NewWindow("go gol")
	window.SetContent(gv)
	window.Canvas().SetOnTypedRune(gv.typedRune)
	gv.animate()

	window.ShowAndRun()

}
