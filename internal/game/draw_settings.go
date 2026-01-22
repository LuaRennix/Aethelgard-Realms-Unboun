package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// DrawSettings отрисовывает экран настроек
func (g *Game) DrawSettings(screen *ebiten.Image) {
	// Затемнение фона
	ebitenutil.DrawRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 200})

	// Заголовок "Settings"
	settingsTitle := g.getText("Settings")
	titleBounds := text.BoundString(g.titleFont, settingsTitle)
	titleWidth := titleBounds.Max.X - titleBounds.Min.X
	titleX := ScreenWidth/2 - titleWidth/2
	titleY := 100

	// Тень
	text.Draw(screen, settingsTitle, g.titleFont, titleX+2, titleY+2, color.RGBA{0, 0, 0, 100})
	text.Draw(screen, settingsTitle, g.titleFont, titleX, titleY, color.RGBA{230, 220, 200, 255})

	// Декоративная линия
	lineY := float64(titleY + 20)
	ebitenutil.DrawRect(screen, float64(ScreenWidth/2-titleWidth/2), lineY, float64(titleWidth), 2, color.RGBA{180, 170, 150, 100})

	// Получаем позицию курсора
	mouseX, mouseY := ebiten.CursorPosition()

	// === СЕКЦИЯ ЯЗЫКА ===
	languageLabel := g.getText("Language")
	labelBounds := text.BoundString(g.menuFont, languageLabel)
	labelWidth := labelBounds.Max.X - labelBounds.Min.X
	labelX := ScreenWidth/2 - labelWidth/2
	labelY := 220

	text.Draw(screen, languageLabel, g.menuFont, labelX+2, labelY+2, color.RGBA{0, 0, 0, 100})
	text.Draw(screen, languageLabel, g.menuFont, labelX, labelY, color.RGBA{200, 190, 180, 255})

	// Кнопки языка
	russianButtonX := ScreenWidth/2 - 200
	russianButtonY := 260
	russianButtonWidth := 180
	russianButtonHeight := 70

	russianButtonHover := mouseX >= russianButtonX && mouseX <= russianButtonX+russianButtonWidth &&
		mouseY >= russianButtonY && mouseY <= russianButtonY+russianButtonHeight
	russianSelected := g.language == LanguageRussian

	var russianBgColor, russianBorderColor, russianTextColor color.RGBA
	if russianSelected {
		russianBgColor = color.RGBA{100, 80, 150, 255}
		russianBorderColor = color.RGBA{150, 120, 200, 255}
		russianTextColor = color.RGBA{255, 255, 255, 255}
	} else if russianButtonHover {
		russianBgColor = color.RGBA{80, 60, 120, 255}
		russianBorderColor = color.RGBA{150, 120, 200, 255}
		russianTextColor = color.RGBA{255, 255, 255, 255}
	} else {
		russianBgColor = color.RGBA{50, 40, 80, 255}
		russianBorderColor = color.RGBA{100, 80, 140, 200}
		russianTextColor = color.RGBA{200, 200, 200, 255}
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

	text.Draw(screen, russianText, g.menuFont, russianTextX+2, russianTextY+2, color.RGBA{0, 0, 0, 150})
	text.Draw(screen, russianText, g.menuFont, russianTextX, russianTextY, russianTextColor)

	if russianSelected {
		dotX := float64(russianButtonX + 15)
		dotY := float64(russianButtonY + russianButtonHeight/2)
		g.drawGlowingDot(screen, dotX, dotY, g.glowIntensity)
	}

	// Кнопка English
	englishButtonX := ScreenWidth/2 + 20
	englishButtonY := 260
	englishButtonWidth := 180
	englishButtonHeight := 70

	englishButtonHover := mouseX >= englishButtonX && mouseX <= englishButtonX+englishButtonWidth &&
		mouseY >= englishButtonY && mouseY <= englishButtonY+englishButtonHeight
	englishSelected := g.language == LanguageEnglish

	var englishBgColor, englishBorderColor, englishTextColor color.RGBA
	if englishSelected {
		englishBgColor = color.RGBA{100, 80, 150, 255}
		englishBorderColor = color.RGBA{150, 120, 200, 255}
		englishTextColor = color.RGBA{255, 255, 255, 255}
	} else if englishButtonHover {
		englishBgColor = color.RGBA{80, 60, 120, 255}
		englishBorderColor = color.RGBA{150, 120, 200, 255}
		englishTextColor = color.RGBA{255, 255, 255, 255}
	} else {
		englishBgColor = color.RGBA{50, 40, 80, 255}
		englishBorderColor = color.RGBA{100, 80, 140, 200}
		englishTextColor = color.RGBA{200, 200, 200, 255}
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

	text.Draw(screen, englishText, g.menuFont, englishTextX+2, englishTextY+2, color.RGBA{0, 0, 0, 150})
	text.Draw(screen, englishText, g.menuFont, englishTextX, englishTextY, englishTextColor)

	if englishSelected {
		dotX := float64(englishButtonX + 15)
		dotY := float64(englishButtonY + englishButtonHeight/2)
		g.drawGlowingDot(screen, dotX, dotY, g.glowIntensity)
	}

	// === СЕКЦИЯ ГРОМКОСТИ ===
	volumeLabelY := 380
	volumeLabel := g.getText("Volume")
	volumeLabelBounds := text.BoundString(g.menuFont, volumeLabel)
	volumeLabelWidth := volumeLabelBounds.Max.X - volumeLabelBounds.Min.X
	volumeLabelX := ScreenWidth/2 - volumeLabelWidth/2

	text.Draw(screen, volumeLabel, g.menuFont, volumeLabelX+2, volumeLabelY+2, color.RGBA{0, 0, 0, 100})
	text.Draw(screen, volumeLabel, g.menuFont, volumeLabelX, volumeLabelY, color.RGBA{200, 190, 180, 255})

	// Слайдер громкости
	volumeSliderX := ScreenWidth/2 - 150
	volumeSliderY := 420
	volumeSliderWidth := 300
	volumeSliderHeight := 12

	// Фон слайдера
	ebitenutil.DrawRect(screen, float64(volumeSliderX), float64(volumeSliderY), float64(volumeSliderWidth), float64(volumeSliderHeight), color.RGBA{40, 30, 60, 255})

	// Заполненная часть
	filledWidth := float64(volumeSliderWidth) * g.masterVolume
	ebitenutil.DrawRect(screen, float64(volumeSliderX), float64(volumeSliderY), filledWidth, float64(volumeSliderHeight), color.RGBA{100, 80, 150, 255})

	// Обводка
	ebitenutil.DrawRect(screen, float64(volumeSliderX), float64(volumeSliderY), float64(volumeSliderWidth), 2, color.RGBA{120, 100, 160, 200})
	ebitenutil.DrawRect(screen, float64(volumeSliderX), float64(volumeSliderY+volumeSliderHeight-2), float64(volumeSliderWidth), 2, color.RGBA{120, 100, 160, 200})
	ebitenutil.DrawRect(screen, float64(volumeSliderX), float64(volumeSliderY), 2, float64(volumeSliderHeight), color.RGBA{120, 100, 160, 200})
	ebitenutil.DrawRect(screen, float64(volumeSliderX+volumeSliderWidth-2), float64(volumeSliderY), 2, float64(volumeSliderHeight), color.RGBA{120, 100, 160, 200})

	// Ползунок
	knobX := float64(volumeSliderX) + filledWidth
	knobY := float64(volumeSliderY + volumeSliderHeight/2)
	knobSize := 12.0

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
	volumeTextY := volumeSliderY + volumeSliderHeight + 30

	text.Draw(screen, volumeText, g.menuFont, volumeTextX+2, volumeTextY+2, color.RGBA{0, 0, 0, 150})
	text.Draw(screen, volumeText, g.menuFont, volumeTextX, volumeTextY, color.RGBA{220, 200, 255, 255})

	// Кнопка "Назад"
	backButtonX := ScreenWidth/2 - 100
	backButtonY := 520
	backButtonWidth := 200
	backButtonHeight := 50

	backButtonHover := mouseX >= backButtonX && mouseX <= backButtonX+backButtonWidth &&
		mouseY >= backButtonY && mouseY <= backButtonY+backButtonHeight

	var backBgColor, backBorderColor, backTextColor color.RGBA
	if backButtonHover {
		backBgColor = color.RGBA{80, 60, 120, 255}
		backBorderColor = color.RGBA{150, 120, 200, 255}
		backTextColor = color.RGBA{255, 255, 255, 255}
	} else {
		backBgColor = color.RGBA{50, 40, 80, 255}
		backBorderColor = color.RGBA{100, 80, 140, 200}
		backTextColor = color.RGBA{200, 200, 200, 255}
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

	text.Draw(screen, backText, g.menuFont, backTextX+2, backTextY+2, color.RGBA{0, 0, 0, 150})
	text.Draw(screen, backText, g.menuFont, backTextX, backTextY, backTextColor)
}
