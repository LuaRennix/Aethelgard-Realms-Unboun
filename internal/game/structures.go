package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/image/font"
)

type MenuItem struct {
	label    string
	selected bool
}

type Game struct {
	state         int
	language      int
	background    *ebiten.Image
	menuItems     []MenuItem
	selectedIndex int
	titleFont     font.Face
	menuFont      font.Face
	glowIntensity float64
	glowDirection float64
	keyPressed    bool

	// Audio fields
	audioContext     *audio.Context
	bgMusic          *audio.Player
	masterVolume     float64 // 0.0 - 1.0
	isDraggingVolume bool

	// Video fields
	videoFrames   []*ebiten.Image
	currentFrame  int
	frameCount    int
	frameRate     float64
	lastFrameTime float64
	videoPlaying  bool
}
