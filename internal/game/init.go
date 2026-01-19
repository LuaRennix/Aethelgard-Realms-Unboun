package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func NewGame() *Game {
	// Создаем видеопроигрыватель для видеофона
	videoPlayer := NewDesktopVideoPlayer("assets/background_video.mp4")

	// Загружаем готический шрифт для заголовка
	tt, err := opentype.Parse(AethelgardFont)
	if err != nil {
		log.Fatal("Failed to parse font: ", err)
	}

	// Большой шрифт для заголовка (готический стиль)
	titleFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    72,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("Failed to create title font: ", err)
	}

	// Меньший шрифт для пунктов меню
	menuFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    28,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("Failed to create menu font: ", err)
	}

	// Инициализируем аудио контекст
	audioContext := audio.NewContext(44100)

	game := &Game{
		state:            MenuState,
		language:         LanguageRussian,
		background:       nil, // Не используем статичное изображение
		videoPlayer:      videoPlayer,
		titleFont:        titleFace,
		menuFont:         menuFace,
		selectedIndex:    0,
		glowIntensity:    0,
		glowDirection:    0.02,
		keyPressed:       false,
		audioContext:     audioContext,
		masterVolume:     0.7, // 70% по умолчанию
		isDraggingVolume: false,
		menuItems: []MenuItem{
			{"New Game", false},
			{"Load Game", false},
			{"Settings", false},
			{"Exit", false},
		},
	}

	// Загружаем и запускаем фоновую музыку
	if err := game.loadAndPlayBackgroundMusic(); err != nil {
		log.Printf("Warning: Failed to load background music: %v", err)
		// Продолжаем работу без музыки
	} else {
		log.Println("Background music loaded successfully!")
	}

	return game
}
