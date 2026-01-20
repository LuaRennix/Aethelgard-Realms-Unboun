package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func NewGame() *Game {
	// Видео
	videoPlayer, err := NewVideoPlayer("assets/menu-background-video.mp4", 30)
	if err != nil {
		log.Printf("Failed to init VideoPlayer: %v, falling back to static image", err)

		img, _, err := ebitenutil.NewImageFromFile("assets/background.png")
		if err != nil {
			log.Fatal("Failed to load fallback background: ", err)
		}

		videoPlayer = &VideoPlayer{
			framePaths:   nil,
			currentIndex: 0,
			frameCount:   1,
			frameDelay:   1,
			frameTimer:   0,
			fps:          0,
			currentFrame: img,
		}
	}

	// === Заголовочный шрифт (Tana Uncial SP) ===
	ttTitle, err := opentype.Parse(TanaFont)
	if err != nil {
		log.Fatal("Failed to parse Tana Uncial SP font:", err)
	}

	titleFace, err := opentype.NewFace(ttTitle, &opentype.FaceOptions{
		Size:    72,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("Failed to create title font:", err)
	}

	// === Шрифт меню (HUD Sonic X1) ===
	ttMenu, err := opentype.Parse(HudSonicFont)
	if err != nil {
		log.Fatal("Failed to parse HUD Sonic X1 font:", err)
	}

	menuFace, err := opentype.NewFace(ttMenu, &opentype.FaceOptions{
		Size:    28,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("Failed to create menu font:", err)
	}

	// Аудио
	audioContext := audio.NewContext(44100)

	// === Создаём игру ===
	game := &Game{
		state:         MenuState,
		language:      LanguageRussian,
		videoPlayer:   videoPlayer,
		titleFont:     titleFace, // Tana Uncial SP
		menuFont:      menuFace,  // HUD Sonic X1
		selectedIndex: 0,
		glowIntensity: 0,
		glowDirection: 0.02,
		keyPressed:    false,
		audioContext:  audioContext,
		masterVolume:  0.7,
		menuItems: []MenuItem{
			{"New Game", false},
			{"Load Game", false},
			{"Settings", false},
			{"Exit", false},
		},
	}

	// Музыка
	if err := game.loadAndPlayBackgroundMusic(); err != nil {
		log.Printf("Warning: Failed to load background music: %v", err)
	}

	return game
}
