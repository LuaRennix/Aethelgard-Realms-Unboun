package game

import (
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// loadVideoForMenu загружает видеофайл и создает анимацию из кадров
func (g *Game) loadVideoForMenu() error {
	// Проверяем, существует ли видеофайл
	if _, err := os.Stat("assets/menu-background-video.mp4"); err != nil {
		return err
	}

	// Попробуем использовать ffmpeg для извлечения кадров
	// Сначала проверим, установлен ли ffmpeg
	if err := exec.Command("which", "ffmpeg").Run(); err != nil {
		// Если ffmpeg недоступен, просто используем статичный фон
		log.Println("Warning: ffmpeg not found. Using static background instead of video.")
		return nil
	}

	// Создаем временные кадры из видео
	frameDir := "tmp_video_frames"
	if err := exec.Command("mkdir", "-p", frameDir).Run(); err != nil {
		return err
	}

	// Извлекаем кадры из видео (например, 15 кадров в секунду для оптимизации)
	cmd := exec.Command("ffmpeg", "-i", "assets/menu-background-video.mp4",
		"-vf", "fps=15,scale="+strconv.Itoa(ScreenWidth)+":"+strconv.Itoa(ScreenHeight),
		"-start_number", "0",
		filepath.Join(frameDir, "frame_%04d.png"))

	if err := cmd.Run(); err != nil {
		log.Printf("Error extracting video frames: %v", err)
		return err
	}

	// Подсчитываем количество извлеченных кадров
	cmd = exec.Command("sh", "-c", "ls "+frameDir+"/*.png | wc -l")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	frameCount, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return err
	}

	// Загружаем кадры в память
	frames := make([]*ebiten.Image, frameCount)
	for i := 0; i < frameCount; i++ {
		framePath := filepath.Join(frameDir, "frame_"+padNumber(i, 4)+".png")
		img, _, err := ebitenutil.NewImageFromFile(framePath)
		if err != nil {
			log.Printf("Warning: Could not load frame %s: %v", framePath, err)
			continue
		}
		frames[i] = img
	}

	g.videoFrames = frames
	g.frameCount = frameCount
	g.currentFrame = 0
	g.frameRate = 15.0 // 15 FPS
	g.lastFrameTime = 0

	log.Printf("Successfully loaded %d video frames", frameCount)

	return nil
}

// padNumber форматирует число с ведущими нулями
func padNumber(n, width int) string {
	numStr := strconv.Itoa(n)
	for len(numStr) < width {
		numStr = "0" + numStr
	}
	return numStr
}

// updateVideo обновляет текущий кадр видео
func (g *Game) updateVideo() {
	if !g.videoPlaying || len(g.videoFrames) == 0 {
		return
	}

	// Увеличиваем время и проверяем, нужно ли переключить кадр
	// В Ebiten Update вызывается каждый кадр, но нам нужно учитывать время
	g.lastFrameTime++

	// Переключаем кадр каждые несколько тактов в зависимости от частоты кадров
	frameInterval := 60.0 / g.frameRate // Предполагаем, что Ebiten работает на 60 FPS
	if g.lastFrameTime >= frameInterval {
		g.currentFrame = (g.currentFrame + 1) % g.frameCount
		g.lastFrameTime = 0
	}
}

// drawVideo рисует текущий кадр видео
func (g *Game) drawVideo(screen *ebiten.Image) {
	if !g.videoPlaying || len(g.videoFrames) == 0 {
		return
	}

	currentFrameImg := g.videoFrames[g.currentFrame]
	if currentFrameImg == nil {
		return
	}

	// Рисуем видео на фоне
	op := &ebiten.DrawImageOptions{}
	// Масштаб уже применен при извлечении кадров с помощью ffmpeg
	screen.DrawImage(currentFrameImg, op)
}
