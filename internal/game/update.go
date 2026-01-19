package game

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Update() error {
	// Проверяем статус музыки (для отладки)
	if g.bgMusic != nil && g.state == MenuState {
		if !g.bgMusic.IsPlaying() {
			log.Println("Music stopped unexpectedly, restarting...")
			g.bgMusic.Rewind()
			g.bgMusic.Play()
		}
	}

	// Обработка ESC (исправлено - теперь не вылетает)
	escPressed := ebiten.IsKeyPressed(ebiten.KeyEscape)

	if escPressed && !g.keyPressed {
		g.keyPressed = true

		if g.state == GameState {
			g.state = MenuState
		} else if g.state == SettingsState {
			g.state = MenuState
		} else {
			os.Exit(0)
		}

		// Обновляем состояние музыки при изменении состояния игры
		g.updateMusicState()
	}

	// Обработка меню
	if g.state == MenuState {
		// Анимация свечения
		g.glowIntensity += g.glowDirection
		if g.glowIntensity > 1.0 {
			g.glowIntensity = 1.0
			g.glowDirection = -0.02
		} else if g.glowIntensity < 0.3 {
			g.glowIntensity = 0.3
			g.glowDirection = 0.02
		}

		// Навигация клавишами (исправлено - теперь не "летает")
		upPressed := ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)
		downPressed := ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)
		enterPressed := ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace)

		// Обработка только при новом нажатии (не зажатии)
		if (upPressed || downPressed || enterPressed) && !g.keyPressed {
			g.keyPressed = true

			if upPressed {
				g.selectedIndex--
				if g.selectedIndex < 0 {
					g.selectedIndex = len(g.menuItems) - 1
				}
			}

			if downPressed {
				g.selectedIndex++
				if g.selectedIndex >= len(g.menuItems) {
					g.selectedIndex = 0
				}
			}

			if enterPressed {
				g.handleMenuAction(g.selectedIndex)
			}
		}

		// Сброс флага при отпускании клавиши
		if !upPressed && !downPressed && !enterPressed && !escPressed {
			g.keyPressed = false
		}

		// Обработка мыши
		mouseX, mouseY := ebiten.CursorPosition()
		menuX := 80
		startY := 320

		for i := range g.menuItems {
			itemY := startY + i*60
			bounds := text.BoundString(g.menuFont, g.getText(g.menuItems[i].label))
			itemHeight := bounds.Max.Y - bounds.Min.Y

			if mouseX >= menuX && mouseX <= menuX+300 &&
				mouseY >= itemY-itemHeight && mouseY <= itemY+10 {
				g.selectedIndex = i

				// Клик мышью
				if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
					g.handleMenuAction(i)
				}
			}
		}
	}

	// Обработка меню настроек
	if g.state == SettingsState {
		// Анимация свечения продолжается
		g.glowIntensity += g.glowDirection
		if g.glowIntensity > 1.0 {
			g.glowIntensity = 1.0
			g.glowDirection = -0.02
		} else if g.glowIntensity < 0.3 {
			g.glowIntensity = 0.3
			g.glowDirection = 0.02
		}

		mouseX, mouseY := ebiten.CursorPosition()
		mouseClicked := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

		// Обработка переключателя языка (Russian)
		russianButtonX := ScreenWidth/2 - 130
		russianButtonY := 230
		russianButtonWidth := 120
		russianButtonHeight := 45

		if mouseX >= russianButtonX && mouseX <= russianButtonX+russianButtonWidth &&
			mouseY >= russianButtonY && mouseY <= russianButtonY+russianButtonHeight {
			if mouseClicked && !g.keyPressed {
				g.keyPressed = true
				g.language = LanguageRussian
			}
		}

		// Обработка переключателя языка (English)
		englishButtonX := ScreenWidth/2 + 10
		englishButtonY := 230
		englishButtonWidth := 120
		englishButtonHeight := 45

		if mouseX >= englishButtonX && mouseX <= englishButtonX+englishButtonWidth &&
			mouseY >= englishButtonY && mouseY <= englishButtonY+englishButtonHeight {
			if mouseClicked && !g.keyPressed {
				g.keyPressed = true
				g.language = LanguageEnglish
			}
		}

		// Обработка кнопки назад
		backButtonX := ScreenWidth/2 - 100
		backButtonY := 480
		backButtonWidth := 200
		backButtonHeight := 50

		if mouseX >= backButtonX && mouseX <= backButtonX+backButtonWidth &&
			mouseY >= backButtonY && mouseY <= backButtonY+backButtonHeight {
			if mouseClicked && !g.keyPressed {
				g.keyPressed = true
				g.state = MenuState
			}
		}

		// Обработка слайдера громкости
		volumeSliderX := ScreenWidth/2 - 150
		volumeSliderY := 360
		volumeSliderWidth := 300
		volumeSliderHeight := 18

		// Проверяем клик по слайдеру
		if mouseY >= volumeSliderY && mouseY <= volumeSliderY+volumeSliderHeight &&
			mouseX >= volumeSliderX && mouseX <= volumeSliderX+volumeSliderWidth {
			if mouseClicked {
				g.isDraggingVolume = true
			}
		}

		// Обработка перетаскивания слайдера
		if g.isDraggingVolume {
			if mouseClicked {
				// Вычисляем новую громкость на основе позиции мыши
				relativeX := float64(mouseX - volumeSliderX)
				g.masterVolume = relativeX / float64(volumeSliderWidth)

				// Ограничиваем значение от 0 до 1
				if g.masterVolume < 0 {
					g.masterVolume = 0
				}
				if g.masterVolume > 1 {
					g.masterVolume = 1
				}

				// Применяем громкость
				g.updateMusicState()
			} else {
				g.isDraggingVolume = false
			}
		}

		// Сброс флага при отпускании мыши
		if !mouseClicked && !escPressed {
			g.keyPressed = false
		}
	}

	return nil
}
