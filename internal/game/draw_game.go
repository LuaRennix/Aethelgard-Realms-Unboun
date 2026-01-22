package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// DrawGame отрисовывает игровое состояние
func (g *Game) DrawGame(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 200})

	gameText := g.getText("Game Started")
	bounds := text.BoundString(g.titleFont, gameText)
	textWidth := bounds.Max.X - bounds.Min.X
	textX := ScreenWidth/2 - int(textWidth)/2
	textY := ScreenHeight / 2

	text.Draw(screen, gameText, g.titleFont, textX+3, textY+3, color.RGBA{0, 0, 0, 200})
	text.Draw(screen, gameText, g.titleFont, textX, textY, color.RGBA{220, 200, 180, 255})

	hintText := g.getText("Press ESC")
	hintBounds := text.BoundString(g.menuFont, hintText)
	hintWidth := hintBounds.Max.X - hintBounds.Min.X
	hintX := ScreenWidth/2 - int(hintWidth)/2
	hintY := textY + 80

	text.Draw(screen, hintText, g.menuFont, hintX, hintY, color.RGBA{180, 170, 160, 200})
}
