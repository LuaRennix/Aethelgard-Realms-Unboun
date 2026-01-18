package game

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
)

const (
	MenuState = iota
	GameState
	SettingsState
)

const (
	LanguageRussian = iota
	LanguageEnglish
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
}

func NewGame() *Game {
	// Загружаем фон
	img, _, err := ebitenutil.NewImageFromFile("assets/background.png")
	if err != nil {
		log.Fatal("Failed to load background: ", err)
	}

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

	return &Game{
		state:         MenuState,
		language:      LanguageRussian,
		background:    img,
		titleFont:     titleFace,
		menuFont:      menuFace,
		selectedIndex: 0,
		glowIntensity: 0,
		glowDirection: 0.02,
		keyPressed:    false,
		menuItems: []MenuItem{
			{"New Game", false},
			{"Load Game", false},
			{"Settings", false},
			{"Exit", false},
		},
	}
}

func (g *Game) Update() error {
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
		russianButtonX := ScreenWidth/2 - 150
		russianButtonY := 300
		russianButtonWidth := 140
		russianButtonHeight := 50

		if mouseX >= russianButtonX && mouseX <= russianButtonX+russianButtonWidth &&
			mouseY >= russianButtonY && mouseY <= russianButtonY+russianButtonHeight {
			if mouseClicked && !g.keyPressed {
				g.keyPressed = true
				g.language = LanguageRussian
			}
		}

		// Обработка переключателя языка (English)
		englishButtonX := ScreenWidth/2 + 10
		englishButtonY := 300
		englishButtonWidth := 140
		englishButtonHeight := 50

		if mouseX >= englishButtonX && mouseX <= englishButtonX+englishButtonWidth &&
			mouseY >= englishButtonY && mouseY <= englishButtonY+englishButtonHeight {
			if mouseClicked && !g.keyPressed {
				g.keyPressed = true
				g.language = LanguageEnglish
			}
		}

		// Обработка кнопки назад
		backButtonX := ScreenWidth/2 - 100
		backButtonY := 450
		backButtonWidth := 200
		backButtonHeight := 50

		if mouseX >= backButtonX && mouseX <= backButtonX+backButtonWidth &&
			mouseY >= backButtonY && mouseY <= backButtonY+backButtonHeight {
			if mouseClicked && !g.keyPressed {
				g.keyPressed = true
				g.state = MenuState
			}
		}

		// Сброс флага при отпускании мыши
		if !mouseClicked && !escPressed {
			g.keyPressed = false
		}
	}

	return nil
}

func (g *Game) getText(key string) string {
	switch g.language {
	case LanguageRussian:
		switch key {
		case "New Game":
			return "Новая игра"
		case "Load Game":
			return "Загрузить игру"
		case "Settings":
			return "Настройки"
		case "Exit":
			return "Выход"
		case "Language":
			return "Язык"
		case "Back":
			return "Назад"
		case "Game Started":
			return "Игра началась!"
		case "Press ESC":
			return "Нажмите ESC, чтобы вернуться в меню"
		}
	case LanguageEnglish:
		switch key {
		case "New Game":
			return "New Game"
		case "Load Game":
			return "Load Game"
		case "Settings":
			return "Settings"
		case "Exit":
			return "Exit"
		case "Language":
			return "Language"
		case "Back":
			return "Back"
		case "Game Started":
			return "Game Started!"
		case "Press ESC":
			return "Press ESC to return to menu"
		}
	}
	return key
}

func (g *Game) handleMenuAction(index int) {
	switch g.menuItems[index].label {
	case "New Game":
		g.state = GameState
	case "Settings":
		g.state = SettingsState
	case "Exit":
		os.Exit(0)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Рисуем фон с затемнением для атмосферы
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(ScreenWidth)/float64(g.background.Bounds().Dx()),
		float64(ScreenHeight)/float64(g.background.Bounds().Dy()))
	op.ColorScale.Scale(0.4, 0.4, 0.4, 1.0)
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
		titleY := 120

		// Многослойная тень
		for i := 4; i > 0; i-- {
			shadowAlpha := uint8(30 * i)
			text.Draw(screen, settingsTitle, g.titleFont, titleX+i, titleY+i, color.RGBA{0, 0, 0, shadowAlpha})
		}
		text.Draw(screen, settingsTitle, g.titleFont, titleX, titleY, color.RGBA{230, 220, 200, 255})

		// Декоративная линия
		lineY := float64(titleY + 20)
		ebitenutil.DrawRect(screen, float64(ScreenWidth/2-int(titleWidth)/2), lineY, float64(titleWidth), 2, color.RGBA{180, 170, 150, 100})

		// Подзаголовок "Язык / Language"
		languageLabel := g.getText("Language")
		labelBounds := text.BoundString(g.menuFont, languageLabel)
		labelWidth := labelBounds.Max.X - labelBounds.Min.X
		labelX := ScreenWidth/2 - labelWidth/2
		labelY := 230

		text.Draw(screen, languageLabel, g.menuFont, labelX+2, labelY+2, color.RGBA{0, 0, 0, 150})
		text.Draw(screen, languageLabel, g.menuFont, labelX, labelY, color.RGBA{200, 190, 180, 255})

		mouseX, mouseY := ebiten.CursorPosition()

		// Кнопка "Русский"
		russianButtonX := ScreenWidth/2 - 150
		russianButtonY := 300
		russianButtonWidth := 140
		russianButtonHeight := 50

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

		// Текст кнопки Russian
		russianText := "Русский"
		russianTextBounds := text.BoundString(g.menuFont, russianText)
		russianTextWidth := russianTextBounds.Max.X - russianTextBounds.Min.X
		russianTextHeight := russianTextBounds.Max.Y - russianTextBounds.Min.Y
		russianTextX := russianButtonX + russianButtonWidth/2 - russianTextWidth/2
		russianTextY := russianButtonY + russianButtonHeight/2 + russianTextHeight/3

		text.Draw(screen, russianText, g.menuFont, russianTextX+2, russianTextY+2, color.RGBA{0, 0, 0, 180})
		text.Draw(screen, russianText, g.menuFont, russianTextX, russianTextY, russianTextColor)

		// Индикатор выбора
		if russianSelected {
			dotX := float64(russianButtonX + 15)
			dotY := float64(russianButtonY + russianButtonHeight/2)
			g.drawGlowingDot(screen, dotX, dotY, g.glowIntensity)
		}

		// Кнопка "English"
		englishButtonX := ScreenWidth/2 + 10
		englishButtonY := 300
		englishButtonWidth := 140
		englishButtonHeight := 50

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

		// Рисуем кнопку English
		ebitenutil.DrawRect(screen, float64(englishButtonX), float64(englishButtonY), float64(englishButtonWidth), float64(englishButtonHeight), englishBgColor)
		ebitenutil.DrawRect(screen, float64(englishButtonX), float64(englishButtonY), float64(englishButtonWidth), 2, englishBorderColor)
		ebitenutil.DrawRect(screen, float64(englishButtonX), float64(englishButtonY+englishButtonHeight-2), float64(englishButtonWidth), 2, englishBorderColor)
		ebitenutil.DrawRect(screen, float64(englishButtonX), float64(englishButtonY), 2, float64(englishButtonHeight), englishBorderColor)
		ebitenutil.DrawRect(screen, float64(englishButtonX+englishButtonWidth-2), float64(englishButtonY), 2, float64(englishButtonHeight), englishBorderColor)

		// Текст кнопки English
		englishText := "English"
		englishTextBounds := text.BoundString(g.menuFont, englishText)
		englishTextWidth := englishTextBounds.Max.X - englishTextBounds.Min.X
		englishTextHeight := englishTextBounds.Max.Y - englishTextBounds.Min.Y
		englishTextX := englishButtonX + englishButtonWidth/2 - englishTextWidth/2
		englishTextY := englishButtonY + englishButtonHeight/2 + englishTextHeight/3

		text.Draw(screen, englishText, g.menuFont, englishTextX+2, englishTextY+2, color.RGBA{0, 0, 0, 180})
		text.Draw(screen, englishText, g.menuFont, englishTextX, englishTextY, englishTextColor)

		// Индикатор выбора
		if englishSelected {
			dotX := float64(englishButtonX + 15)
			dotY := float64(englishButtonY + englishButtonHeight/2)
			g.drawGlowingDot(screen, dotX, dotY, g.glowIntensity)
		}

		// Кнопка "Назад / Back"
		backButtonX := ScreenWidth/2 - 100
		backButtonY := 450
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

		// Рисуем кнопку Back
		ebitenutil.DrawRect(screen, float64(backButtonX), float64(backButtonY), float64(backButtonWidth), float64(backButtonHeight), backBgColor)
		ebitenutil.DrawRect(screen, float64(backButtonX), float64(backButtonY), float64(backButtonWidth), 2, backBorderColor)
		ebitenutil.DrawRect(screen, float64(backButtonX), float64(backButtonY+backButtonHeight-2), float64(backButtonWidth), 2, backBorderColor)
		ebitenutil.DrawRect(screen, float64(backButtonX), float64(backButtonY), 2, float64(backButtonHeight), backBorderColor)
		ebitenutil.DrawRect(screen, float64(backButtonX+backButtonWidth-2), float64(backButtonY), 2, float64(backButtonHeight), backBorderColor)

		// Текст кнопки Back
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
