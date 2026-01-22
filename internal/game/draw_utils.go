package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// drawGlowingDot рисует светящуюся точку
func (g *Game) drawGlowingDot(screen *ebiten.Image, x, y, intensity float64) {
	for i := 0; i < 4; i++ {
		size := float64(8 - i*2)
		alpha := uint8(50 * intensity * float64(4-i) / 4.0)
		offset := size / 2
		ebitenutil.DrawRect(screen, x-offset, y-offset, size, size, color.RGBA{200, 160, 255, alpha})
	}

	coreAlpha := uint8(220 + 35*intensity)
	ebitenutil.DrawRect(screen, x-1, y-1, 2, 2, color.RGBA{240, 220, 255, coreAlpha})
}

// drawBottomDecoration рисует декоративную линию и версию внизу экрана
func (g *Game) drawBottomDecoration(screen *ebiten.Image) {
	y := float64(ScreenHeight - 40)
	ebitenutil.DrawRect(screen, 60, y, ScreenWidth-120, 1, color.RGBA{180, 170, 150, 80})

	versionText := "v0.1.0"
	text.Draw(screen, versionText, g.menuFont, ScreenWidth-150, int(y)+30, color.RGBA{120, 110, 100, 150})
}
