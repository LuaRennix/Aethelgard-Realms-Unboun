package main

import (
	"aethelgard/internal/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("Aethelgard: Realms Unbound")

	// Заменяем SetWindowSize на полноэкранный режим
	ebiten.SetFullscreen(true) // Это ключевая строка!

	// Опционально: включаем плавное масштабирование
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
