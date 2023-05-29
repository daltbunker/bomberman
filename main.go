package main

import (
	"fmt"
  "math"
	// "image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	padding = 30

	xTiles   = 7
	yTiles   = 7
	tileSize = 32

	screenWidth  = xTiles*tileSize + padding * 2
	screenHeight = yTiles*tileSize + padding * 2
)

var (
	bombermanImage *ebiten.Image
	rockImage      *ebiten.Image
)

type Game struct {
	keys       []ebiten.Key
	levelMap   [yTiles][xTiles]int // x-val is inner array
	characterX int
	characterY int
}

func init() {
	var err error

	bombermanImage, _, err = ebitenutil.NewImageFromFile("main-character.png")
	if err != nil {
		log.Fatal(err)
	}

	rockImage, _, err = ebitenutil.NewImageFromFile("rock.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) init() {
  g.characterX = 0
  g.characterY = 0
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x5e, 0x33, 0x19, 0x01})

	g.handleKeyPress()

	for y := 0; y < len(g.levelMap); y++ {
		for x := 0; x < len(g.levelMap[y]); x++ {
			if x%2 != 0 && y%2 != 0 {
				g.drawRock(screen, x, y)
			}
		}
	}

  g.drawCharacter(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("x: %v, y: %v", g.characterX, g.characterY))
}

func (g *Game) drawRock(screen *ebiten.Image, x int, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x)*tileSize+padding, float64(y)*tileSize+padding)
	screen.DrawImage(rockImage, op)
}

func (g *Game) drawCharacter(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.characterX)+padding, float64(g.characterY)+padding)
	screen.DrawImage(bombermanImage, op)
}

func (g *Game) handleKeyPress() {
	if len(g.keys) == 0 {
		return
	}  

  turnBuffer := 10
  bottomBorder := tileSize * yTiles - 32
  rightBorder := tileSize * xTiles - 32

  allowMoveHorizontal := g.characterY % 64 > 64 - turnBuffer || g.characterY % 64 < turnBuffer || g.characterY == bottomBorder 
  allowMoveVertical := g.characterX % 64 > 64 - turnBuffer || g.characterX % 64 < turnBuffer || g.characterX == rightBorder 

	// Selecting the first key prevents diagonal movement
	pressedKey := g.keys[0]
	if pressedKey == ebiten.KeyArrowUp && g.characterY > 0 && allowMoveVertical {
    if g.characterX % 64 != 0 {
      g.characterX = 64 * int(math.Round(float64(g.characterX) / 64))
    }
    g.characterY -= 2
	}
	if pressedKey == ebiten.KeyArrowDown && g.characterY < bottomBorder && allowMoveVertical {
    if g.characterX % 64 != 0 {
      g.characterX = 64 * int(math.Round(float64(g.characterX) / 64))
    }
    g.characterY += 2
	}
	if pressedKey == ebiten.KeyArrowLeft && g.characterX > 0 && allowMoveHorizontal {
    if g.characterY % 64 != 0 {
      g.characterY = 64 * int(math.Round(float64(g.characterY) / 64))
    }
    g.characterX -= 2
	}
	if pressedKey == ebiten.KeyArrowRight && g.characterX < rightBorder && allowMoveHorizontal {
    if g.characterY % 64 != 0 {
      g.characterY = 64 * int(math.Round(float64(g.characterY) / 64))
    }
    g.characterX += 2
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Bomberman")

	g := &Game{}
	g.init()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
