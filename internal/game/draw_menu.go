package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// DrawMenu отрисовывает главное меню
func (g *Game) DrawMenu(screen *ebiten.Image) {
	title := "Aethelgard"
	titleX := 60
	titleY := 120

	for i := 5; i > 0; i-- {
		shadowAlpha := uint8(30 * i)
		text.Draw(screen, title, g.titleFont, titleX+i, titleY+i, color.RGBA{0, 0, 0, shadowAlpha})
	}

	text.Draw(screen, title, g.titleFont, titleX, titleY, color.RGBA{230, 220, 200, 255})

	subtitle := "Realms Unbound"
	subtitleY := titleY + 40
	text.Draw(screen, subtitle, g.menuFont, titleX+10, subtitleY, color.RGBA{180, 170, 150, 200})

	titleBounds := text.BoundString(g.titleFont, title)
	titleWidth := titleBounds.Max.X - titleBounds.Min.X
	ebitenutil.DrawRect(screen, float64(titleX), float64(subtitleY+10), float64(titleWidth), 2, color.RGBA{180, 170, 150, 100})

	menuX := 80
	startY := 320

	for i, item := range g.menuItems {
		itemY := startY + i*60
		isSelected := i == g.selectedIndex

		itemText := g.getText(item.label)
		textBounds := text.BoundString(g.menuFont, itemText)
		textWidth := textBounds.Max.X - textBounds.Min.X

		var textColor color.RGBA
		if isSelected {
			glowValue := uint8(220 + 35*g.glowIntensity)

			lineY := float64(itemY + 8)
			lineWidth := float64(textWidth + 10)

			for j := 0; j < 3; j++ {
				glowAlpha := uint8(float64(60-j*15) * g.glowIntensity)
				ebitenutil.DrawRect(screen, float64(menuX-5-j), lineY+float64(j), lineWidth+float64(j*2), 1, color.RGBA{180, 140, 255, glowAlpha})
			}

			ebitenutil.DrawRect(screen, float64(menuX-5), lineY, lineWidth, 2, color.RGBA{200, 160, 255, uint8(200 * g.glowIntensity)})

			dotX := float64(menuX - 25)
			dotY := float64(itemY - 8)
			g.drawGlowingDot(screen, dotX, dotY, g.glowIntensity)

			textColor = color.RGBA{glowValue, glowValue - 20, 255, 255}
		} else {
			textColor = color.RGBA{150, 140, 130, 200}
		}

		text.Draw(screen, itemText, g.menuFont, menuX+2, itemY+2, color.RGBA{0, 0, 0, 100})
		text.Draw(screen, itemText, g.menuFont, menuX, itemY, textColor)
	}

	g.drawBottomDecoration(screen)
}
