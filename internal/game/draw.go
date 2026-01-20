package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {
	// Рисуем видео фон
	var videoFrame *ebiten.Image
	if g.videoPlayer != nil {
		videoFrame = g.videoPlayer.CurrentFrame()
	}

	if videoFrame != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(
			float64(ScreenWidth)/float64(videoFrame.Bounds().Dx()),
			float64(ScreenHeight)/float64(videoFrame.Bounds().Dy()),
		)
		op.ColorM.Scale(0.4, 0.4, 0.4, 1.0)
		screen.DrawImage(videoFrame, op)
	} else {
		ebitenutil.DrawRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 255})
	}

	if g.state == MenuState {
		title := "Aethelgard"
		titleX := 60
		titleY := 120

		for i := 5; i > 0; i-- {
			shadowAlpha := uint8(30 * i)
			text.Draw(screen, title, g.titleFont, titleX+i, titleY+i, color.RGBA{0, 0, 0, shadowAlpha})
		}

		text.Draw(screen, title, g.titleFont, titleX, titleY, color.RGBA{230, 220, 200, 255})

		subtitle := "Realms Unbound"
		subtitleY := titleY + 40
		text.Draw(screen, subtitle, g.menuFont, titleX+10, subtitleY, color.RGBA{180, 170, 150, 200})

		titleBounds := text.BoundString(g.titleFont, title)
		titleWidth := titleBounds.Max.X - titleBounds.Min.X
		ebitenutil.DrawRect(screen, float64(titleX), float64(subtitleY+10), float64(titleWidth), 2, color.RGBA{180, 170, 150, 100})

		menuX := 80
		startY := 320

		for i, item := range g.menuItems {
			itemY := startY + i*60
			isSelected := i == g.selectedIndex

			itemText := g.getText(item.label)
			textBounds := text.BoundString(g.menuFont, itemText)
			textWidth := textBounds.Max.X - textBounds.Min.X

			var textColor color.RGBA
			if isSelected {
				glowValue := uint8(220 + 35*g.glowIntensity)

				lineY := float64(itemY + 8)
				lineWidth := float64(textWidth + 10)

				for j := 0; j < 3; j++ {
					glowAlpha := uint8(float64(60-j*15) * g.glowIntensity)
					ebitenutil.DrawRect(screen, float64(menuX-5-j), lineY+float64(j), lineWidth+float64(j*2), 1, color.RGBA{180, 140, 255, glowAlpha})
				}

				ebitenutil.DrawRect(screen, float64(menuX-5), lineY, lineWidth, 2, color.RGBA{200, 160, 255, uint8(200 * g.glowIntensity)})

				dotX := float64(menuX - 25)
				dotY := float64(itemY - 8)
				g.drawGlowingDot(screen, dotX, dotY, g.glowIntensity)

				textColor = color.RGBA{glowValue, glowValue - 20, 255, 255}
			} else {
				textColor = color.RGBA{150, 140, 130, 200}
			}

			text.Draw(screen, itemText, g.menuFont, menuX+2, itemY+2, color.RGBA{0, 0, 0, 100})
			text.Draw(screen, itemText, g.menuFont, menuX, itemY, textColor)
		}

		g.drawBottomDecoration(screen)
	}

	if g.state == SettingsState {
		// Затемнение фона
		ebitenutil.DrawRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 220})

		// Заголовок "НАСТРОЙКИ"
		settingsTitle := g.getText("Settings")
		titleBounds := text.BoundString(g.titleFont, settingsTitle)
		titleWidth := titleBounds.Max.X - titleBounds.Min.X
		titleX := ScreenWidth/2 - titleWidth/2
		titleY := 80

		// Тень — чёрная + фиолетовая
		for i := 3; i > 0; i-- {
			text.Draw(screen, settingsTitle, g.titleFont, titleX+i, titleY+i, color.RGBA{0, 0, 0, uint8(30 * i)})
			text.Draw(screen, settingsTitle, g.titleFont, titleX+i/2, titleY+i/2, color.RGBA{100, 80, 160, uint8(20 * i)})
		}
		// Основной текст
		text.Draw(screen, settingsTitle, g.titleFont, titleX, titleY, color.RGBA{230, 220, 200, 255})
		// Декоративная линия
		lineY := float64(titleY + 20)
		ebitenutil.DrawRect(screen, float64(ScreenWidth/2-titleWidth/2), lineY, float64(titleWidth), 2, color.RGBA{180, 170, 150, 100})

		// === СЕКЦИЯ ЯЗЫКА ===
		languageLabel := g.getText("Language")
		labelBounds := text.BoundString(g.menuFont, languageLabel)
		labelWidth := labelBounds.Max.X - labelBounds.Min.X
		labelX := ScreenWidth/2 - labelWidth/2
		labelY := 200

		text.Draw(screen, languageLabel, g.menuFont, labelX+2, labelY+2, color.RGBA{0, 0, 0, 150})
		text.Draw(screen, languageLabel, g.menuFont, labelX, labelY, color.RGBA{200, 190, 180, 255})

		mouseX, mouseY := ebiten.CursorPosition()

		// --- КНОПКА "РУССКИЙ" ---
		russianText := "Русский"
		russianTextBounds := text.BoundString(g.menuFont, russianText)
		russianTextWidth := float64(russianTextBounds.Max.X - russianTextBounds.Min.X)
		russianTextHeight := float64(russianTextBounds.Max.Y - russianTextBounds.Min.Y)

		buttonPaddingX := 16.0
		buttonPaddingY := 12.0
		buttonWidth := russianTextWidth + buttonPaddingX*2
		buttonHeight := russianTextHeight + buttonPaddingY*2

		russianButtonX := float64(ScreenWidth/2 - int(buttonWidth) - 10) // 10 — отступ между кнопками
		russianButtonY := float64(230)

		russianButtonHover := mouseX >= int(russianButtonX) && mouseX <= int(russianButtonX+buttonWidth) &&
			mouseY >= int(russianButtonY) && mouseY <= int(russianButtonY+buttonHeight)
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

		// Рисуем фон кнопки
		ebitenutil.DrawRect(screen, russianButtonX, russianButtonY, buttonWidth, buttonHeight, russianBgColor)

		// Рисуем рамку
		borderSize := 2.0
		ebitenutil.DrawRect(screen, russianButtonX, russianButtonY, buttonWidth, borderSize, russianBorderColor)
		ebitenutil.DrawRect(screen, russianButtonX, russianButtonY+buttonHeight-borderSize, buttonWidth, borderSize, russianBorderColor)
		ebitenutil.DrawRect(screen, russianButtonX, russianButtonY, borderSize, buttonHeight, russianBorderColor)
		ebitenutil.DrawRect(screen, russianButtonX+buttonWidth-borderSize, russianButtonY, borderSize, buttonHeight, russianBorderColor)

		// Центрируем текст внутри кнопки
		russianTextX := int(russianButtonX + buttonPaddingX)
		russianTextY := int(russianButtonY + buttonPaddingY + russianTextHeight/2)

		// Тень текста
		text.Draw(screen, russianText, g.menuFont, russianTextX+2, russianTextY+2, color.RGBA{0, 0, 0, 180})
		text.Draw(screen, russianText, g.menuFont, russianTextX, russianTextY, russianTextColor)

		if russianSelected {
			dotX := float64(russianButtonX + 12)
			dotY := float64(russianButtonY + buttonHeight/2)
			g.drawGlowingDot(screen, dotX, dotY, g.glowIntensity)
		}

		// --- КНОПКА "ENGLISH" ---
		englishText := "English"
		englishTextBounds := text.BoundString(g.menuFont, englishText)
		englishTextWidth := float64(englishTextBounds.Max.X - englishTextBounds.Min.X)
		englishTextHeight := float64(englishTextBounds.Max.Y - englishTextBounds.Min.Y)

		englishButtonWidth := englishTextWidth + buttonPaddingX*2
		englishButtonHeight := englishTextHeight + buttonPaddingY*2

		englishButtonX := float64(ScreenWidth/2 + 10) // 10 — отступ между кнопками
		englishButtonY := float64(230)

		englishButtonHover := mouseX >= int(englishButtonX) && mouseX <= int(englishButtonX+englishButtonWidth) &&
			mouseY >= int(englishButtonY) && mouseY <= int(englishButtonY+englishButtonHeight)
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

		ebitenutil.DrawRect(screen, englishButtonX, englishButtonY, englishButtonWidth, englishButtonHeight, englishBgColor)

		ebitenutil.DrawRect(screen, englishButtonX, englishButtonY, englishButtonWidth, borderSize, englishBorderColor)
		ebitenutil.DrawRect(screen, englishButtonX, englishButtonY+englishButtonHeight-borderSize, englishButtonWidth, borderSize, englishBorderColor)
		ebitenutil.DrawRect(screen, englishButtonX, englishButtonY, borderSize, englishButtonHeight, englishBorderColor)
		ebitenutil.DrawRect(screen, englishButtonX+englishButtonWidth-borderSize, englishButtonY, borderSize, englishButtonHeight, englishBorderColor)

		englishTextX := int(englishButtonX + buttonPaddingX)
		englishTextY := int(englishButtonY + buttonPaddingY + englishTextHeight/2)

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

		// Слайдер громкости
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

		// --- КНОПКА "НАЗАД" ---
		backText := g.getText("Back")
		backTextBounds := text.BoundString(g.menuFont, backText)
		backTextWidth := float64(backTextBounds.Max.X - backTextBounds.Min.X)
		backTextHeight := float64(backTextBounds.Max.Y - backTextBounds.Min.Y)

		backButtonWidth := backTextWidth + 32.0   // 16 с каждой стороны
		backButtonHeight := backTextHeight + 24.0 // 12 сверху и снизу

		backButtonX := float64(ScreenWidth/2 - int(backButtonWidth)/2)
		backButtonY := float64(480)

		backButtonHover := mouseX >= int(backButtonX) && mouseX <= int(backButtonX+backButtonWidth) &&
			mouseY >= int(backButtonY) && mouseY <= int(backButtonY+backButtonHeight)

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

		ebitenutil.DrawRect(screen, backButtonX, backButtonY, backButtonWidth, backButtonHeight, backBgColor)

		ebitenutil.DrawRect(screen, backButtonX, backButtonY, backButtonWidth, borderSize, backBorderColor)
		ebitenutil.DrawRect(screen, backButtonX, backButtonY+backButtonHeight-borderSize, backButtonWidth, borderSize, backBorderColor)
		ebitenutil.DrawRect(screen, backButtonX, backButtonY, borderSize, backButtonHeight, backBorderColor)
		ebitenutil.DrawRect(screen, backButtonX+backButtonWidth-borderSize, backButtonY, borderSize, backButtonHeight, backBorderColor)

		backTextX := int(backButtonX + 16)
		backTextY := int(backButtonY + 12 + backTextHeight/2)

		text.Draw(screen, backText, g.menuFont, backTextX+2, backTextY+2, color.RGBA{0, 0, 0, 180})
		text.Draw(screen, backText, g.menuFont, backTextX, backTextY, backTextColor)
	}

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
	for i := 0; i < 4; i++ {
		size := float64(8 - i*2)
		alpha := uint8(50 * intensity * float64(4-i) / 4.0)
		offset := size / 2
		ebitenutil.DrawRect(screen, x-offset, y-offset, size, size, color.RGBA{200, 160, 255, alpha})
	}

	coreAlpha := uint8(220 + 35*intensity)
	ebitenutil.DrawRect(screen, x-1, y-1, 2, 2, color.RGBA{240, 220, 255, coreAlpha})
}

func (g *Game) drawBottomDecoration(screen *ebiten.Image) {
	y := float64(ScreenHeight - 40)
	ebitenutil.DrawRect(screen, 60, y, ScreenWidth-120, 1, color.RGBA{180, 170, 150, 80})

	versionText := "v0.1.0"
	text.Draw(screen, versionText, g.menuFont, ScreenWidth-150, int(y)+30, color.RGBA{120, 110, 100, 150})
}
