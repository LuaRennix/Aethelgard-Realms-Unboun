package game

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/image/font"
)

type MenuItem struct {
	label    string
	selected bool
}

type Game struct {
	// Состояние и язык
	state    int
	language int

	// Меню
	menuItems     []MenuItem
	selectedIndex int

	// Шрифты
	titleFont font.Face
	menuFont  font.Face

	// Эффекты
	glowIntensity float64
	glowDirection float64
	keyPressed    bool

	// Аудио
	audioContext     *audio.Context
	bgMusic          *audio.Player
	masterVolume     float64
	isDraggingVolume bool

	// Видео
	videoPlayer *VideoPlayer
}
