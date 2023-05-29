package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameOX     = 0
	frameOY     = 0
	frameWidth  = 37
	frameHeight = 48
	frameCount  = 4
)

var (
	bombermanImage *ebiten.Image
	rockImage      *ebiten.Image
)

type Game struct {
	keys  []ebiten.Key
	count int
	xPos  int
	yPos  int
}

func init() {
	var err error

	bombermanImage, _, err = ebitenutil.NewImageFromFile("bomberman.png")
	if err != nil {
		log.Fatal(err)
	}

	rockImage, _, err = ebitenutil.NewImageFromFile("rock.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	g.count++
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x5e, 0x33, 0x19, 0x01})
	g.drawRocks(screen)

	g.handleKeyPress()
	g.drawCharacter(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("x: %v, y: %v", g.xPos, g.yPos))
}

func (g *Game) drawRocks(screen *ebiten.Image) {
	for x := 1; x < 7; x += 2 {
		for y := 0; y < 5; y += 2 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*32)+screenWidth/4, float64(y*32)+screenHeight/4)
			op.GeoM.Scale(0.8, 0.8)
			screen.DrawImage(rockImage, op)
		}
	}
}

func (g *Game) drawCharacter(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(float64(g.xPos)+screenWidth/2, float64(g.yPos)+screenHeight/2)
	op.GeoM.Scale(0.5, 0.5)
	i := (g.count / 7) % frameCount // Update runs every millisecond, 8 is delay between switching frames
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(bombermanImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) handleKeyPress() {
	if len(g.keys) < 1 {
		return
	}
	// Selecting the first key, prevents diagonal movement
	pressedKey := g.keys[0]
	if pressedKey == ebiten.KeyArrowUp && g.yPos > 0 {
		g.yPos -= 4
	}
	if pressedKey == ebiten.KeyArrowDown && g.yPos < screenHeight - 40 {
		g.yPos += 4
	}
	if pressedKey == ebiten.KeyArrowLeft && g.xPos > 0 {
		g.xPos -= 4
	}
	if pressedKey == ebiten.KeyArrowRight && g.xPos < screenWidth - 20 && g.yPos > 130 && g.yPos < 170 {
		g.yPos = 150
		g.xPos += 4
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Bomberman")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
