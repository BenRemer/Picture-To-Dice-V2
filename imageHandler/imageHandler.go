package imageHandler

import (
	"fmt"
	"gonum.org/v1/gonum/stat"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
)

func CreateDiceFromGrayImage(img image.Image, diceImageSize int) (image.Image, error) { // todo remove saving of gray image and do that pixel by pixel at runtime
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	dice := createDiceArray(diceImageSize)
	xCount, yCount := 0, 0

	diceImg := image.NewGray(image.Rectangle{Min: image.Point{}, Max: image.Point{X: width * diceImageSize, Y: height * diceImageSize}})
	for y := 0; y < height; y++ {
		xCount = 0
		for x := 0; x < width; x++ {
			pixelValue := rgbaGetR(img.At(x, y))
			var dieNumber int
			if pixelValue > 212 {
				dieNumber = 0
			} else if pixelValue > 170 {
				dieNumber = 1
			} else if pixelValue > 127 {
				dieNumber = 2
			} else if pixelValue > 85 {
				dieNumber = 3
			} else if pixelValue > 42 {
				dieNumber = 4
			} else {
				dieNumber = 5
			}
			die := dice[dieNumber]

			for j := 0; j < diceImageSize; j++ {
				for i := 0; i < diceImageSize; i++ {
					diceImg.Set(xCount+x+i, yCount+y+j, die.At(i, j))
				}
			}
			xCount += diceImageSize
		}
		yCount += diceImageSize
	}
	return diceImg, nil
}

func createDiceArray(size int) [6]image.Image {
	var array [6]image.Image
	for i := 0; i < 6; i++ {
		die, err := GetGrayImageFromPath(fmt.Sprintf("images/dice/%vdice%v.png", size, i+1))
		if err != nil {
			fmt.Println("Error creating dice array")
			panic(err)
		}
		array[i] = die
	}
	return array
}

func PixelateImage(img image.Image, blurSize int) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X/blurSize, bounds.Max.Y/blurSize // only smaller for now todo add larger images

	pixImg := image.NewGray(image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}})
	newX := 0
	newY := 0
	for y := 0; y < bounds.Max.Y; y += blurSize {
		for x := 0; x < bounds.Max.X; x += blurSize {
			ave := blockMean(img, x, y, blurSize)
			pixImg.Set(newX, newY, color.Gray{Y: ave})
			newX += 1
		}
		newX = 0
		newY += 1
	}
	return pixImg, nil
}

func blockMean(img image.Image, startingX, staringY, size int) uint8 {
	yLow := staringY
	yHigh := staringY + size
	xLow := startingX
	xHigh := startingX + size
	var block []float64
	for y := yLow; y < yHigh; y++ {
		for x := xLow; x < xHigh; x++ {
			value := rgbaGetR(img.At(x, y))
			block = append(block, float64(value))
		}
	}
	m := stat.Mean(block, nil)
	return uint8(m)
}

func rgbaGetR(rgba color.Color) int {
	r, _, _, _ := rgba.RGBA()
	r = r / 257
	return int(r)
}

func SaveImage(image image.Image, name string) {
	path := fmt.Sprintf("output/%s.png", name)
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, image); err != nil {
		log.Fatal(err)
	}
}

func GetGrayImageFromPath(filePath string) (image.Image, error) {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error reading file:", filePath)
		return nil, err
	}
	defer file.Close()

	grayImg, err := imageToGreyscale(file)

	if err != nil {
		fmt.Println("Error converting file to gray")
		return nil, err
	}

	return grayImg, err
}

func imageToGreyscale(file io.Reader) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	grayImg := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}
	return grayImg, nil
}
