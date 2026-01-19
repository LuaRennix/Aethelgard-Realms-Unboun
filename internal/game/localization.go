package game

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
		case "Volume":
			return "Громкость"
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
		case "Volume":
			return "Volume"
		}
	}
	return key
}
