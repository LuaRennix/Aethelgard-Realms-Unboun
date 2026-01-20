package game

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// VideoPlayer — стриминговый плеер: в памяти только один кадр
type VideoPlayer struct {
	framePaths   []string
	currentIndex int
	frameCount   int
	frameDelay   int
	frameTimer   int
	fps          int
	currentFrame *ebiten.Image
}

func NewVideoPlayer(videoPath string, targetFPS int) (*VideoPlayer, error) {
	videoDir := "assets/video_frames"

	if _, err := os.Stat(videoDir); os.IsNotExist(err) {
		log.Printf("Video frames directory not found, path: %s", videoDir)
		return nil, fmt.Errorf("video frames directory not found")
	}

	files, err := os.ReadDir(videoDir)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	var framePaths []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			continue
		}
		fullPath := filepath.Join(videoDir, file.Name())
		framePaths = append(framePaths, fullPath)
	}

	if len(framePaths) == 0 {
		return nil, fmt.Errorf("no video frames found in %s", videoDir)
	}

	log.Printf("VideoPlayer: found %d frame files", len(framePaths))

	frameDelay := 1
	if targetFPS > 0 {
		frameDelay = int(ebiten.TPS() / targetFPS)
		if frameDelay <= 0 {
			frameDelay = 1
		}
	}

	firstImg, _, err := ebitenutil.NewImageFromFile(framePaths[0])
	if err != nil {
		return nil, fmt.Errorf("failed to load first frame: %w", err)
	}

	v := &VideoPlayer{
		framePaths:   framePaths,
		currentIndex: 0,
		frameCount:   len(framePaths),
		frameDelay:   frameDelay,
		frameTimer:   0,
		fps:          targetFPS,
		currentFrame: firstImg,
	}

	log.Printf("VideoPlayer initialized: %d frames, targetFPS=%d, frameDelay=%d",
		v.frameCount, v.fps, v.frameDelay)

	return v, nil
}

func (v *VideoPlayer) Update() {
	if v.frameCount == 0 || v.currentFrame == nil {
		return
	}

	v.frameTimer++
	if v.frameTimer < v.frameDelay {
		return
	}
	v.frameTimer = 0

	v.currentIndex++
	if v.currentIndex >= v.frameCount {
		v.currentIndex = 0
	}

	if v.currentFrame != nil {
		v.currentFrame.Dispose()
		v.currentFrame = nil
	}

	path := v.framePaths[v.currentIndex]
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Printf("VideoPlayer: failed to load frame %s: %v", path, err)
		return
	}

	v.currentFrame = img
}

func (v *VideoPlayer) CurrentFrame() *ebiten.Image {
	return v.currentFrame
}

func (v *VideoPlayer) Close() {
	if v.currentFrame != nil {
		v.currentFrame.Dispose()
		v.currentFrame = nil
	}
}

func (v *VideoPlayer) FrameCount() int {
	return v.frameCount
}

func (v *VideoPlayer) FPS() int {
	return v.fps
}
