package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

func (g *Game) Draw(screen *ebiten.Image) {
	// Рисуем фон с затемнением для атмосферы
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(ScreenWidth)/float64(g.background.Bounds().Dx()),
		float64(ScreenHeight)/float64(g.background.Bounds().Dy()))

	// Используем ColorM для старых версий Ebiten
	op.ColorM.Scale(0.4, 0.4, 0.4, 1.0)
	screen.DrawImage(g.background, op)

	// Меню
	if g.state == MenuState {
		// Название игры
		title := "Aethelgard"
		titleX := 60
		titleY := 120

		// Многослойная тень
		for i := 5; i > 0; i-- {
			shadowAlpha := uint8(30 * i)
			text.Draw(screen, title, g.titleFont, titleX+i, titleY+i, color.RGBA{0, 0, 0, shadowAlpha})
		}

		text.Draw(screen, title, g.titleFont, titleX, titleY, color.RGBA{230, 220, 200, 255})

		// Подзаголовок
		subtitle := "Realms Unbound"
		subtitleY := titleY + 40
		text.Draw(screen, subtitle, g.menuFont, titleX+10, subtitleY, color.RGBA{180, 170, 150, 200})

		// Декоративная линия
		titleBounds := text.BoundString(g.titleFont, title)
		titleWidth := titleBounds.Max.X - titleBounds.Min.X
		ebitenutil.DrawRect(screen, float64(titleX), float64(subtitleY+10), float64(titleWidth), 2, color.RGBA{180, 170, 150, 100})

		// Пункты меню
		menuX := 80
		startY := 320

		for i, item := range g.menuItems {
			itemY := startY + i*60
			isSelected := i == g.selectedIndex

			// Получаем текст и его размеры
			itemText := g.getText(item.label)
			textBounds := text.BoundString(g.menuFont, itemText)
			textWidth := textBounds.Max.X - textBounds.Min.X

			var textColor color.RGBA
			if isSelected {
				// Элегантное тонкое подчеркивание
				glowValue := uint8(220 + 35*g.glowIntensity)

				// Тонкая светящаяся линия под текстом
				lineY := float64(itemY + 8)
				lineWidth := float64(textWidth + 10)

				// Свечение линии (3 слоя для мягкости)
				for j := 0; j < 3; j++ {
					glowAlpha := uint8(float64(60-j*15) * g.glowIntensity)
					ebitenutil.DrawRect(screen, float64(menuX-5-j), lineY+float64(j), lineWidth+float64(j*2), 1, color.RGBA{180, 140, 255, glowAlpha})
				}

				// Основная яркая линия
				ebitenutil.DrawRect(screen, float64(menuX-5), lineY, lineWidth, 2, color.RGBA{200, 160, 255, uint8(200 * g.glowIntensity)})

				// Маленькая светящаяся точка слева
				dotX := float64(menuX - 25)
				dotY := float64(itemY - 8)
				g.drawGlowingDot(screen, dotX, dotY, g.glowIntensity)

				// Яркий текст
				textColor = color.RGBA{glowValue, glowValue - 20, 255, 255}
			} else {
				// Невыбранные пункты - приглушенные
				textColor = color.RGBA{150, 140, 130, 200}
			}

			// Мягкая тень текста
			text.Draw(screen, itemText, g.menuFont, menuX+2, itemY+2, color.RGBA{0, 0, 0, 100})

			// Основной текст
			text.Draw(screen, itemText, g.menuFont, menuX, itemY, textColor)
		}

		g.drawBottomDecoration(screen)
	}

	// Меню настроек
	if g.state == SettingsState {
		// Затемнение фона
		ebitenutil.DrawRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 220})

		// Заголовок настроек
		settingsTitle := g.getText("Settings")
		titleBounds := text.BoundString(g.titleFont, settingsTitle)
		titleWidth := titleBounds.Max.X - titleBounds.Min.X
		titleX := ScreenWidth/2 - titleWidth/2
		titleY := 100

		// Многослойная тень
		for i := 4; i > 0; i-- {
			shadowAlpha := uint8(30 * i)
			text.Draw(screen, settingsTitle, g.titleFont, titleX+i, titleY+i, color.RGBA{0, 0, 0, shadowAlpha})
		}
		text.Draw(screen, settingsTitle, g.titleFont, titleX, titleY, color.RGBA{230, 220, 200, 255})

		// Декоративная линия
		lineY := float64(titleY + 20)
		ebitenutil.DrawRect(screen, float64(ScreenWidth/2-int(titleWidth)/2), lineY, float64(titleWidth), 2, color.RGBA{180, 170, 150, 100})

		// === СЕКЦИЯ ЯЗЫКА ===
		languageLabel := g.getText("Language")
		labelBounds := text.BoundString(g.menuFont, languageLabel)
		labelWidth := labelBounds.Max.X - labelBounds.Min.X
		labelX := ScreenWidth/2 - labelWidth/2
		labelY := 200

		text.Draw(screen, languageLabel, g.menuFont, labelX+2, labelY+2, color.RGBA{0, 0, 0, 150})
		text.Draw(screen, languageLabel, g.menuFont, labelX, labelY, color.RGBA{200, 190, 180, 255})

		mouseX, mouseY := ebiten.CursorPosition()

		// Кнопки языка (уменьшены и ближе друг к другу)
		russianButtonX := ScreenWidth/2 - 130
		russianButtonY := 230
		russianButtonWidth := 120
		russianButtonHeight := 45

		russianButtonHover := mouseX >= russianButtonX && mouseX <= russianButtonX+russianButtonWidth &&
			mouseY >= russianButtonY && mouseY <= russianButtonY+russianButtonHeight
		russianSelected := g.language == LanguageRussian

		var russianBgColor, russianBorderColor, russianTextColor color.RGBA
		if russianSelected {
			russianBgColor = color.RGBA{80, 60, 120, 255}
			russianBorderColor = color.RGBA{150, 120, 200, 255}
			russianTextColor = color.RGBA{255, 240, 220, 255}
		} else if russianButtonHover {
			russianBgColor = color.RGBA{50, 40, 70, 255}
			russianBorderColor = color.RGBA{120, 100, 160, 200}
			russianTextColor = color.RGBA{220, 210, 200, 255}
		} else {
			russianBgColor = color.RGBA{30, 25, 45, 255}
			russianBorderColor = color.RGBA{100, 80, 140, 150}
			russianTextColor = color.RGBA{180, 170, 160, 200}
		}

		// Рисуем кнопку Russian
		ebitenutil.DrawRect(screen, float64(russianButtonX), float64(russianButtonY), float64(russianButtonWidth), float64(russianButtonHeight), russianBgColor)
		ebitenutil.DrawRect(screen, float64(russianButtonX), float64(russianButtonY), float64(russianButtonWidth), 2, russianBorderColor)
		ebitenutil.DrawRect(screen, float64(russianButtonX), float64(russianButtonY+russianButtonHeight-2), float64(russianButtonWidth), 2, russianBorderColor)
		ebitenutil.DrawRect(screen, float64(russianButtonX), float64(russianButtonY), 2, float64(russianButtonHeight), russianBorderColor)
		ebitenutil.DrawRect(screen, float64(russianButtonX+russianButtonWidth-2), float64(russianButtonY), 2, float64(russianButtonHeight), russianBorderColor)

		russianText := "Русский"
		russianTextBounds := text.BoundString(g.menuFont, russianText)
		russianTextWidth := russianTextBounds.Max.X - russianTextBounds.Min.X
		russianTextHeight := russianTextBounds.Max.Y - russianTextBounds.Min.Y
		russianTextX := russianButtonX + russianButtonWidth/2 - russianTextWidth/2
		russianTextY := russianButtonY + russianButtonHeight/2 + russianTextHeight/3

		text.Draw(screen, russianText, g.menuFont, russianTextX+2, russianTextY+2, color.RGBA{0, 0, 0, 180})
		text.Draw(screen, russianText, g.menuFont, russianTextX, russianTextY, russianTextColor)

		if russianSelected {
			dotX := float64(russianButtonX + 12)
			dotY := float64(russianButtonY + russianButtonHeight/2)
			g.drawGlowingDot(screen, dotX, dotY, g.glowIntensity)
		}

		// Кнопка English
		englishButtonX := ScreenWidth/2 + 10
		englishButtonY := 230
		englishButtonWidth := 120
		englishButtonHeight := 45

		englishButtonHover := mouseX >= englishButtonX && mouseX <= englishButtonX+englishButtonWidth &&
			mouseY >= englishButtonY && mouseY <= englishButtonY+englishButtonHeight
		englishSelected := g.language == LanguageEnglish

		var englishBgColor, englishBorderColor, englishTextColor color.RGBA
		if englishSelected {
			englishBgColor = color.RGBA{80, 60, 120, 255}
			englishBorderColor = color.RGBA{150, 120, 200, 255}
			englishTextColor = color.RGBA{255, 240, 220, 255}
		} else if englishButtonHover {
			englishBgColor = color.RGBA{50, 40, 70, 255}
			englishBorderColor = color.RGBA{120, 100, 160, 200}
			englishTextColor = color.RGBA{220, 210, 200, 255}
		} else {
			englishBgColor = color.RGBA{30, 25, 45, 255}
			englishBorderColor = color.RGBA{100, 80, 140, 150}
			englishTextColor = color.RGBA{180, 170, 160, 200}
		}

		ebitenutil.DrawRect(screen, float64(englishButtonX), float64(englishButtonY), float64(englishButtonWidth), float64(englishButtonHeight), englishBgColor)
		ebitenutil.DrawRect(screen, float64(englishButtonX), float64(englishButtonY), float64(englishButtonWidth), 2, englishBorderColor)
		ebitenutil.DrawRect(screen, float64(englishButtonX), float64(englishButtonY+englishButtonHeight-2), float64(englishButtonWidth), 2, englishBorderColor)
		ebitenutil.DrawRect(screen, float64(englishButtonX), float64(englishButtonY), 2, float64(englishButtonHeight), englishBorderColor)
		ebitenutil.DrawRect(screen, float64(englishButtonX+englishButtonWidth-2), float64(englishButtonY), 2, float64(englishButtonHeight), englishBorderColor)

		englishText := "English"
		englishTextBounds := text.BoundString(g.menuFont, englishText)
		englishTextWidth := englishTextBounds.Max.X - englishTextBounds.Min.X
		englishTextHeight := englishTextBounds.Max.Y - englishTextBounds.Min.Y
		englishTextX := englishButtonX + englishButtonWidth/2 - englishTextWidth/2
		englishTextY := englishButtonY + englishButtonHeight/2 + englishTextHeight/3

		text.Draw(screen, englishText, g.menuFont, englishTextX+2, englishTextY+2, color.RGBA{0, 0, 0, 180})
		text.Draw(screen, englishText, g.menuFont, englishTextX, englishTextY, englishTextColor)

		if englishSelected {
			dotX := float64(englishButtonX + 12)
			dotY := float64(englishButtonY + englishButtonHeight/2)
			g.drawGlowingDot(screen, dotX, dotY, g.glowIntensity)
		}

		// === СЕКЦИЯ ГРОМКОСТИ ===
		volumeLabelY := 330
		volumeLabel := g.getText("Volume")
		volumeLabelBounds := text.BoundString(g.menuFont, volumeLabel)
		volumeLabelWidth := volumeLabelBounds.Max.X - volumeLabelBounds.Min.X
		volumeLabelX := ScreenWidth/2 - volumeLabelWidth/2

		text.Draw(screen, volumeLabel, g.menuFont, volumeLabelX+2, volumeLabelY+2, color.RGBA{0, 0, 0, 150})
		text.Draw(screen, volumeLabel, g.menuFont, volumeLabelX, volumeLabelY, color.RGBA{200, 190, 180, 255})

		// Слайдер громкости (более компактный)
		volumeSliderX := ScreenWidth/2 - 150
		volumeSliderY := 360
		volumeSliderWidth := 300
		volumeSliderHeight := 18

		// Фон слайдера
		ebitenutil.DrawRect(screen, float64(volumeSliderX), float64(volumeSliderY), float64(volumeSliderWidth), float64(volumeSliderHeight), color.RGBA{30, 25, 45, 255})

		// Заполненная часть
		filledWidth := float64(volumeSliderWidth) * g.masterVolume
		ebitenutil.DrawRect(screen, float64(volumeSliderX), float64(volumeSliderY), filledWidth, float64(volumeSliderHeight), color.RGBA{100, 70, 150, 255})

		// Обводка
		ebitenutil.DrawRect(screen, float64(volumeSliderX), float64(volumeSliderY), float64(volumeSliderWidth), 2, color.RGBA{120, 100, 160, 200})
		ebitenutil.DrawRect(screen, float64(volumeSliderX), float64(volumeSliderY+volumeSliderHeight-2), float64(volumeSliderWidth), 2, color.RGBA{120, 100, 160, 200})
		ebitenutil.DrawRect(screen, float64(volumeSliderX), float64(volumeSliderY), 2, float64(volumeSliderHeight), color.RGBA{120, 100, 160, 200})
		ebitenutil.DrawRect(screen, float64(volumeSliderX+volumeSliderWidth-2), float64(volumeSliderY), 2, float64(volumeSliderHeight), color.RGBA{120, 100, 160, 200})

		// Ползунок
		knobX := float64(volumeSliderX) + filledWidth
		knobY := float64(volumeSliderY + volumeSliderHeight/2)
		knobSize := 10.0

		for i := 0; i < 3; i++ {
			size := knobSize + float64(i*3)
			alpha := uint8(50 - i*15)
			offset := size / 2
			ebitenutil.DrawRect(screen, knobX-offset, knobY-offset, size, size, color.RGBA{150, 120, 200, alpha})
		}

		ebitenutil.DrawRect(screen, knobX-knobSize/2, knobY-knobSize/2, knobSize, knobSize, color.RGBA{200, 160, 255, 255})

		// Процент
		volumePercent := int(g.masterVolume * 100)
		volumeText := fmt.Sprintf("%d%%", volumePercent)
		volumeTextBounds := text.BoundString(g.menuFont, volumeText)
		volumeTextWidth := volumeTextBounds.Max.X - volumeTextBounds.Min.X
		volumeTextX := ScreenWidth/2 - volumeTextWidth/2
		volumeTextY := volumeSliderY + volumeSliderHeight + 35

		text.Draw(screen, volumeText, g.menuFont, volumeTextX+2, volumeTextY+2, color.RGBA{0, 0, 0, 180})
		text.Draw(screen, volumeText, g.menuFont, volumeTextX, volumeTextY, color.RGBA{220, 200, 255, 255})

		// Кнопка "Назад" (ниже и по центру)
		backButtonX := ScreenWidth/2 - 100
		backButtonY := 480
		backButtonWidth := 200
		backButtonHeight := 50

		backButtonHover := mouseX >= backButtonX && mouseX <= backButtonX+backButtonWidth &&
			mouseY >= backButtonY && mouseY <= backButtonY+backButtonHeight

		var backBgColor, backBorderColor, backTextColor color.RGBA
		if backButtonHover {
			backBgColor = color.RGBA{60, 50, 90, 255}
			backBorderColor = color.RGBA{120, 100, 180, 255}
			backTextColor = color.RGBA{255, 240, 220, 255}
		} else {
			backBgColor = color.RGBA{40, 30, 60, 255}
			backBorderColor = color.RGBA{100, 80, 160, 200}
			backTextColor = color.RGBA{200, 190, 180, 255}
		}

		ebitenutil.DrawRect(screen, float64(backButtonX), float64(backButtonY), float64(backButtonWidth), float64(backButtonHeight), backBgColor)
		ebitenutil.DrawRect(screen, float64(backButtonX), float64(backButtonY), float64(backButtonWidth), 2, backBorderColor)
		ebitenutil.DrawRect(screen, float64(backButtonX), float64(backButtonY+backButtonHeight-2), float64(backButtonWidth), 2, backBorderColor)
		ebitenutil.DrawRect(screen, float64(backButtonX), float64(backButtonY), 2, float64(backButtonHeight), backBorderColor)
		ebitenutil.DrawRect(screen, float64(backButtonX+backButtonWidth-2), float64(backButtonY), 2, float64(backButtonHeight), backBorderColor)

		backText := g.getText("Back")
		backTextBounds := text.BoundString(g.menuFont, backText)
		backTextWidth := backTextBounds.Max.X - backTextBounds.Min.X
		backTextHeight := backTextBounds.Max.Y - backTextBounds.Min.Y
		backTextX := backButtonX + backButtonWidth/2 - backTextWidth/2
		backTextY := backButtonY + backButtonHeight/2 + backTextHeight/3

		text.Draw(screen, backText, g.menuFont, backTextX+2, backTextY+2, color.RGBA{0, 0, 0, 180})
		text.Draw(screen, backText, g.menuFont, backTextX, backTextY, backTextColor)
	}

	// Игровой процесс
	if g.state == GameState {
		ebitenutil.DrawRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 200})

		gameText := g.getText("Game Started")
		bounds := text.BoundString(g.titleFont, gameText)
		textWidth := bounds.Max.X - bounds.Min.X
		textX := ScreenWidth/2 - int(textWidth)/2
		textY := ScreenHeight / 2

		text.Draw(screen, gameText, g.titleFont, textX+3, textY+3, color.RGBA{0, 0, 0, 200})
		text.Draw(screen, gameText, g.titleFont, textX, textY, color.RGBA{220, 200, 180, 255})

		hintText := g.getText("Press ESC")
		hintBounds := text.BoundString(g.menuFont, hintText)
		hintWidth := hintBounds.Max.X - hintBounds.Min.X
		hintX := ScreenWidth/2 - int(hintWidth)/2
		hintY := textY + 80

		text.Draw(screen, hintText, g.menuFont, hintX, hintY, color.RGBA{180, 170, 160, 200})
	}
}

func (g *Game) drawGlowingDot(screen *ebiten.Image, x, y, intensity float64) {
	// Мягкое свечение (меньше и изящнее)
	for i := 0; i < 4; i++ {
		size := float64(8 - i*2)
		alpha := uint8(50 * intensity * float64(4-i) / 4.0)
		offset := size / 2
		ebitenutil.DrawRect(screen, x-offset, y-offset, size, size, color.RGBA{200, 160, 255, alpha})
	}

	// Яркое ядро (маленькое)
	coreAlpha := uint8(220 + 35*intensity)
	ebitenutil.DrawRect(screen, x-1, y-1, 2, 2, color.RGBA{240, 220, 255, coreAlpha})
}

func (g *Game) drawBottomDecoration(screen *ebiten.Image) {
	y := float64(ScreenHeight - 40)
	ebitenutil.DrawRect(screen, 60, y, ScreenWidth-120, 1, color.RGBA{180, 170, 150, 80})

	versionText := "v0.1.0"
	text.Draw(screen, versionText, g.menuFont, ScreenWidth-150, int(y)+30, color.RGBA{120, 110, 100, 150})
}