package main

import (
	"fmt"
	"math"
	"syscall/js"
	"time"

	"github.com/sina-am/chess/chess"
)

var (
	WhiteSquire         = "#e9edcc"
	BlackSquire         = "#779954"
	SelectedSquireColor = "#f4f67e"
	BoardSize           = 700
	SquireSize          = BoardSize / 8
)

func getSquireColor(x int, y int) string {
	if y%2 == 0 {
		if x%2 == 0 {
			return BlackSquire
		}
		return WhiteSquire
	}
	if x%2 == 0 {
		return WhiteSquire
	}
	return BlackSquire
}

type Image struct {
	blob js.Value
	meta map[string]string
}

type ChessUI struct {
	game     *chess.ChessEngine
	draw2d   js.Value
	canvas   js.Value
	document js.Value

	lastClicked time.Time
	pickedPiece *chess.Piece
	viewAs      chess.Color

	images         map[*chess.Piece]Image
	pickupHandlers []HandlePickupPiece
	dropHandlers   []HandleDropPiece
}

type HandlePickupPiece func(piece *chess.Piece) error
type HandleDropPiece func(piece *chess.Piece, x, y int) error

func NewChessUI(game *chess.ChessEngine, viewAs chess.Color) *ChessUI {
	document := js.Global().Get("document")

	loadingElement := document.Call("getElementById", "loading")
	loadingElement.Get("classList").Call("add", "d-none")

	gameElement := document.Call("getElementById", "game")
	gameElement.Get("classList").Call("remove", "d-none")

	canvasElement := document.Call("getElementById", "board")
	canvas := document.Call("createElement", "canvas")
	canvas.Set("width", BoardSize)
	canvas.Set("height", BoardSize)
	canvasElement.Call("appendChild", canvas)
	draw2d := canvas.Call("getContext", "2d")
	chessUI := &ChessUI{
		game:        game,
		draw2d:      draw2d,
		canvas:      canvas,
		lastClicked: time.Now(),
		pickedPiece: nil,
		viewAs:      viewAs,
		images:      map[*chess.Piece]Image{},
		document:    document,

		pickupHandlers: []HandlePickupPiece{},
		dropHandlers:   []HandleDropPiece{},
	}

	canvas.Call("addEventListener", "click", js.FuncOf(chessUI.clickHandler))
	return chessUI
}

func (ui *ChessUI) HookPickupHandler(h HandlePickupPiece) {
	ui.pickupHandlers = append(ui.pickupHandlers, h)
}
func (ui *ChessUI) HookDropHandler(h HandleDropPiece) {
	ui.dropHandlers = append(ui.dropHandlers, h)
}

func (ui *ChessUI) convertToBoardCoordination(x, y int) [2]int {
	dy := 7
	if ui.viewAs == chess.Black {
		dy = 0
	}
	return [2]int{x, int(math.Abs(float64(y - dy)))}
}

func (ui *ChessUI) clickHandler(this js.Value, args []js.Value) any {
	if !ui.isClicked() {
		return nil
	}
	event := args[0]
	rect := ui.canvas.Call("getBoundingClientRect")
	x := int(math.Floor((event.Get("clientX").Float() - rect.Get("x").Float()) / float64(SquireSize)))
	y := int(math.Floor((event.Get("clientY").Float() - rect.Get("y").Float()) / float64(SquireSize)))
	if x < 0 || x > 7 || y < 0 || y > 7 {
		return nil
	}

	loc := ui.convertToBoardCoordination(x, y)
	if ui.pickedPiece == nil {
		ui.handlePickupPiece(loc[0], loc[1])
	} else {
		ui.handleDropPiece(loc[0], loc[1])
	}
	fmt.Println("You clicked", x, y)

	return nil
}
func (ui *ChessUI) changeBackground(x int, y int, color string) {
	ui.drawSquire(x, y, color)
	ui.drawPiece(x, y, ui.game.GetBoard()[y][x])
}

func (ui *ChessUI) handlePickupPiece(x, y int) {
	piece := ui.game.GetBoard()[y][x]
	if piece == nil {
		return
	}

	for _, handler := range ui.pickupHandlers {
		if err := handler(piece); err != nil {
			return
		}
	}
	ui.changeBackground(x, y, SelectedSquireColor)
	ui.pickedPiece = piece

}

func (ui *ChessUI) handleDropPiece(x, y int) {
	for _, handler := range ui.dropHandlers {
		if err := handler(ui.pickedPiece, x, y); err != nil {
			return
		}
	}

	ui.changeBackground(ui.pickedPiece.Location.Col, ui.pickedPiece.Location.Row, getSquireColor(ui.pickedPiece.Location.Col, ui.pickedPiece.Location.Row))
	err := ui.game.Play(
		ui.pickedPiece.Color,
		chess.Move{From: ui.pickedPiece.Location, To: chess.Location{Row: y, Col: x}},
	)

	if err != nil {
		fmt.Println(err)
	} else {
		ui.Render()
	}

	if ui.game.GetResult() != chess.NoResult {
		fmt.Println("result", ui.game.GetResult())
	}
	ui.pickedPiece = nil
}
func (ui *ChessUI) isClicked() bool {
	if (time.Now().Sub(ui.lastClicked)) < 10 {
		return false
	}
	ui.lastClicked = time.Now()
	return true
}

func (ui *ChessUI) drawSquire(x int, y int, color string) {
	loc := ui.convertToBoardCoordination(x, y)
	x = loc[0]
	y = loc[1]

	ui.draw2d.Set("fillStyle", color)
	ui.draw2d.Call("fillRect", x*SquireSize, y*SquireSize, SquireSize, SquireSize)
}

func (ui *ChessUI) drawPiece(x int, y int, piece *chess.Piece) {
	loc := ui.convertToBoardCoordination(x, y)
	x = loc[0]
	y = loc[1]

	image, found := ui.images[piece]
	if !found || image.meta["name"] != piece.Type.GetName() {
		// first time loading images or promotion happened
		blob := js.Global().Get("Image").New()
		blob.Set("src", fmt.Sprintf("static/img/pieces/%s-%s.svg", piece.Type.GetName(), piece.Color))
		blob.Set("onload", js.FuncOf(func(this js.Value, args []js.Value) any {
			ui.draw2d.Call("drawImage", this, x*SquireSize, y*SquireSize, SquireSize, SquireSize)
			return nil
		}))
		ui.images[piece] = Image{blob: blob, meta: map[string]string{"name": piece.Type.GetName()}}
		return
	}

	// Just move the image around
	ui.draw2d.Call("drawImage", image.blob, x*SquireSize, y*SquireSize, SquireSize, SquireSize)
}

func (ui *ChessUI) Render() {
	ui.draw2d.Call("reset")
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			ui.drawSquire(x, y, getSquireColor(x, y))
			if ui.game.GetBoard()[y][x] != nil {
				ui.drawPiece(x, y, ui.game.GetBoard()[y][x])
			}
		}
	}
}

func (ui *ChessUI) Finish(result chess.Result) {
	if result.Reason != chess.Stalemate || result.Reason != chess.Draw {
		ui.document.Call("getElementById", "result").Set("innerText", fmt.Sprintf("%s won by %s", result.WinnerColor, result.Reason))
	} else {
		ui.document.Call("getElementById", "result").Set("innerText", fmt.Sprintf("%s", result.Reason))
	}
}
