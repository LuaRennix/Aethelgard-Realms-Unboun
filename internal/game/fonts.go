package game

import (
	"log"
	"os"
)

var TanaFont []byte
var HudSonicFont []byte

func init() {
	// Заголовочный шрифт (Tana Uncial SP)
	tana, err := os.ReadFile("assets/fonts/TanaUncialSP.ttf")
	if err != nil {
		log.Fatal("Failed to load Tana Uncial SP font: ", err)
	}
	TanaFont = tana

	// Шрифт меню (HUD Sonic X1)
	hud, err := os.ReadFile("assets/fonts/HUD-Sonic-X1.otf")
	if err != nil {
		log.Fatal("Failed to load HUD Sonic X1 font: ", err)
	}
	HudSonicFont = hud
}
