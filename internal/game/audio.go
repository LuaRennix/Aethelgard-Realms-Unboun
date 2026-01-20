package game

import (
	"bytes"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

// loadAndPlayBackgroundMusic загружает и воспроизводит фоновую музыку
func (g *Game) loadAndPlayBackgroundMusic() error {
	log.Println("Attempting to load background music...")

	// Читаем весь файл в память
	audioData, err := os.ReadFile("assets/main_menu_sound.mp3")
	if err != nil {
		log.Printf("Failed to read audio file: %v", err)
		return err
	}
	log.Printf("Audio file loaded, size: %d bytes", len(audioData))

	// Декодируем MP3 из байтов
	decodedStream, err := mp3.DecodeWithoutResampling(bytes.NewReader(audioData))
	if err != nil {
		log.Printf("Failed to decode MP3: %v", err)
		return err
	}
	log.Printf("MP3 decoded successfully, length: %d", decodedStream.Length())

	// Создаем бесконечный поток для зацикливания
	infiniteLoop := audio.NewInfiniteLoop(decodedStream, decodedStream.Length())
	log.Println("Infinite loop created")

	// Создаем плеер
	player, err := g.audioContext.NewPlayer(infiniteLoop)
	if err != nil {
		log.Printf("Failed to create player: %v", err)
		return err
	}
	log.Println("Player created successfully")

	// Устанавливаем громкость на мастер-громкость
	player.SetVolume(g.masterVolume)
	log.Printf("Volume set to %.2f", g.masterVolume)

	g.bgMusic = player

	// Запускаем воспроизведение
	g.bgMusic.Play()
	log.Printf("Music playback started, IsPlaying: %v", g.bgMusic.IsPlaying())

	return nil
}

// updateMusicState управляет воспроизведением музыки в зависимости от текущего состояния игры
func (g *Game) updateMusicState() {
	if g.bgMusic == nil {
		return
	}

	var targetVolume float64

	if g.state == MenuState {
		// В главном меню - полная громкость
		targetVolume = g.masterVolume

		if !g.bgMusic.IsPlaying() {
			g.bgMusic.Rewind()
			g.bgMusic.Play()
			log.Println("Music resumed")
		}
	} else if g.state == SettingsState || g.state == GameState {
		// В настройках и игре - приглушаем до 20%
		targetVolume = g.masterVolume * 0.2

		if !g.bgMusic.IsPlaying() {
			g.bgMusic.Rewind()
			g.bgMusic.Play()
			log.Println("Music resumed (quiet)")
		}
	} else {
		// В других состояниях - останавливаем
		if g.bgMusic.IsPlaying() {
			g.bgMusic.Pause()
			log.Println("Music paused")
		}
		return
	}

	// Применяем громкость
	g.bgMusic.SetVolume(targetVolume)
}
