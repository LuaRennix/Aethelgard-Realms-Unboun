package game

import "os"

func (g *Game) handleMenuAction(index int) {
	switch g.menuItems[index].label {
	case "New Game":
		g.state = GameState
	case "Settings":
		g.state = SettingsState
	case "Exit":
		os.Exit(0)
	}

	// Обновляем состояние музыки после изменения состояния игры
	g.updateMusicState()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}