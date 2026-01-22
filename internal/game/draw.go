package game

import "github.com/hajimehoshi/ebiten/v2"

// Draw отрисовывает текущее состояние игры
func (g *Game) Draw(screen *ebiten.Image) {
	// Рисуем фон
	g.DrawBackground(screen)

	// Отрисовываем текущее состояние
	switch g.state {
	case MenuState:
		g.DrawMenu(screen)
	case SettingsState:
		g.DrawSettings(screen)
	case GameState:
		g.DrawGame(screen)
	}
}
