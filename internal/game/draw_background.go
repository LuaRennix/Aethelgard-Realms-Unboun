package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// DrawBackground отрисовывает фоновое изображение или видео
func (g *Game) DrawBackground(screen *ebiten.Image) {
	var videoFrame *ebiten.Image
	if g.videoPlayer != nil {
		videoFrame = g.videoPlayer.CurrentFrame()
	}

	if videoFrame != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(
			float64(ScreenWidth)/float64(videoFrame.Bounds().Dx()),
			float64(ScreenHeight)/float64(videoFrame.Bounds().Dy()),
		)
		op.ColorM.Scale(0.4, 0.4, 0.4, 1.0)
		screen.DrawImage(videoFrame, op)
	} else {
		ebitenutil.DrawRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 255})
	}
}
