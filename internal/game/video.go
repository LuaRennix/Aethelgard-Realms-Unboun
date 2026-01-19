package game

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// VideoPlayer интерфейс для воспроизведения видео
type VideoPlayer interface {
	Update()
	Draw(screen *ebiten.Image)
	IsPlaying() bool
}

// DesktopVideoPlayer для десктопных версий
type DesktopVideoPlayer struct {
	videoPath    string
	frameDir     string
	frameCount   int
	currentFrame int
	frameImages  []*ebiten.Image
	frameRate    float64
	extracted    bool // флаг, указывающий, были ли кадры уже извлечены
}

func NewDesktopVideoPlayer(videoPath string) *DesktopVideoPlayer {
	player := &DesktopVideoPlayer{
		videoPath: videoPath,
		frameDir:  filepath.Join(filepath.Dir(videoPath), "frames_"+filepath.Base(videoPath)),
		extracted: false,
	}

	// Попробуем получить информацию о видео и извлечь кадры
	err := player.extractFramesIfNeeded()
	if err != nil {
		fmt.Printf("Warning: Could not extract video frames: %v\n", err)
		// Создадим пустой проигрыватель
		return &DesktopVideoPlayer{
			videoPath:    videoPath,
			frameDir:     filepath.Join(filepath.Dir(videoPath), "frames_"+filepath.Base(videoPath)),
			frameCount:   0,
			currentFrame: 0,
			frameImages:  make([]*ebiten.Image, 0),
			frameRate:    30.0, // стандартная частота кадров
			extracted:    false,
		}
	}

	return player
}

// extractFramesIfNeeded извлекает кадры из видео с помощью ffmpeg, если они еще не извлечены
func (vp *DesktopVideoPlayer) extractFramesIfNeeded() error {
	// Проверяем, существуют ли уже извлеченные кадры
	frameFiles, err := ioutil.ReadDir(vp.frameDir)
	if err == nil && len(frameFiles) > 0 {
		// Кадры уже существуют, загружаем их
		vp.loadExistingFrames()
		vp.extracted = true
		return nil
	}

	// Удалить существующую директорию с кадрами (на всякий случай)
	cmd := exec.Command("rm", "-rf", vp.frameDir)
	_ = cmd.Run()

	// Создать директорию для кадров
	cmd = exec.Command("mkdir", "-p", vp.frameDir)
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Извлечь кадры с помощью ffmpeg
	outputPattern := filepath.Join(vp.frameDir, "frame_%05d.png")
	cmd = exec.Command("ffmpeg", "-i", vp.videoPath, "-vf", "fps=30,scale=1280:720", outputPattern)
	err = cmd.Run()
	if err != nil {
		// Попробуем без ограничения по размеру
		cmd = exec.Command("ffmpeg", "-i", vp.videoPath, "-vf", "fps=30", outputPattern)
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	// Подсчитать количество кадров
	frameFiles, err = ioutil.ReadDir(vp.frameDir)
	if err != nil {
		return err
	}

	vp.frameCount = len(frameFiles)
	vp.currentFrame = 0
	vp.frameRate = 30.0

	// Загрузить все кадры в память
	vp.frameImages = make([]*ebiten.Image, vp.frameCount)
	for i := 0; i < vp.frameCount; i++ {
		framePath := filepath.Join(vp.frameDir, fmt.Sprintf("frame_%05d.png", i+1))

		img, _, err := ebitenutil.NewImageFromFile(framePath)
		if err != nil {
			fmt.Printf("Warning: Could not load frame %s: %v\n", framePath, err)
			continue
		}
		vp.frameImages[i] = img
	}

	vp.extracted = true
	return nil
}

// loadExistingFrames загружает уже извлеченные кадры
func (vp *DesktopVideoPlayer) loadExistingFrames() {
	frameFiles, err := ioutil.ReadDir(vp.frameDir)
	if err != nil {
		vp.frameCount = 0
		vp.frameImages = make([]*ebiten.Image, 0)
		return
	}

	vp.frameCount = len(frameFiles)
	vp.currentFrame = 0
	vp.frameRate = 30.0

	// Загрузить все кадры в память
	vp.frameImages = make([]*ebiten.Image, vp.frameCount)
	for i := 0; i < vp.frameCount; i++ {
		framePath := filepath.Join(vp.frameDir, fmt.Sprintf("frame_%05d.png", i+1))

		img, _, err := ebitenutil.NewImageFromFile(framePath)
		if err != nil {
			fmt.Printf("Warning: Could not load frame %s: %v\n", framePath, err)
			continue
		}
		vp.frameImages[i] = img
	}
}

func (vp *DesktopVideoPlayer) Update() {
	// Обновление индекса кадра для анимации
	vp.currentFrame++
	if vp.currentFrame >= vp.frameCount {
		vp.currentFrame = 0 // Зацикливание
	}
}

func (vp *DesktopVideoPlayer) Draw(screen *ebiten.Image) {
	if vp.frameCount > 0 && vp.currentFrame < len(vp.frameImages) && vp.frameImages[vp.currentFrame] != nil {
		// Растянуть кадр на весь экран
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(ScreenWidth)/float64(vp.frameImages[vp.currentFrame].Bounds().Dx()),
			float64(ScreenHeight)/float64(vp.frameImages[vp.currentFrame].Bounds().Dy()))

		// Применить затемнение к видео, чтобы текст был виден
		op.ColorM.Scale(0.6, 0.6, 0.6, 1.0)

		screen.DrawImage(vp.frameImages[vp.currentFrame], op)
	}
}

func (vp *DesktopVideoPlayer) IsPlaying() bool {
	return vp.frameCount > 0 && vp.extracted
}
