package game

import (
	"image"
	"image/color"
	"log"
)

var (
	snakeHeadColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	snakeBodyColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	foodColor      = color.RGBA{R: 0, G: 255, B: 0, A: 255}
)

const bytesPerPixel = 3 // RGB: 3 bytes per pixel

func StartFrameRenderer(gameStateCh chan *Snake, pixelCh chan []byte) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic: %v", x)
		}
	}()

	for {
		gameState, ok := <-gameStateCh
		if !ok {
			break
		}

		img := image.NewRGBA(image.Rect(0, 0, FRAME_WIDTH, FRAME_HEIGHT))
		matrix := gameState.GetMatrix()

		for y := 0; y < len(matrix); y++ {
			for x := 0; x < len(matrix[0]); x++ {
				rectMin := image.Point{X: x * CHUNK_SIZE, Y: y * CHUNK_SIZE}
				rectMax := image.Point{X: rectMin.X + CHUNK_SIZE, Y: rectMin.Y + CHUNK_SIZE}
				switch matrix[y][x] {
				case HEAD:
					drawRectangle(img, rectMin, rectMax, snakeHeadColor)
				case BODY:
					drawRectangle(img, rectMin, rectMax, snakeBodyColor)
				case FOOD:
					drawRectangle(img, rectMin, rectMax, foodColor)
				}
			}
		}

		pixelCh <- convertRGBAtoRGB(img)
	}
}

func drawRectangle(img *image.RGBA, min, max image.Point, col color.RGBA) {
	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			img.Set(x, y, col)
		}
	}
}

func convertRGBAtoRGB(img *image.RGBA) []byte {
	width, height := img.Rect.Dx(), img.Rect.Dy()
	rawRGBData := make([]byte, bytesPerPixel*width*height)

	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.RGBAAt(x, y)
			rawRGBData[idx] = pixel.R
			rawRGBData[idx+1] = pixel.G
			rawRGBData[idx+2] = pixel.B
			idx += bytesPerPixel
		}
	}

	return rawRGBData
}
