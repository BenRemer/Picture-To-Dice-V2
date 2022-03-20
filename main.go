package main

import (
	"gopkg.in/yaml.v2"
	"home/imageHandler"
	"io/ioutil"
)

type conf struct {
	ImagePath string `yaml:"ImagePath"`
	ImageName string `yaml:"ImageName"`
	BlurSize  int    `yaml:"BlurSize"`
	DiceSize  int    `yaml:"DiceSize"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic(err)
	}
	return c
}

func main() {
	var c conf
	c.getConf()

	blurSize := c.BlurSize
	path := c.ImagePath
	path = path + c.ImageName
	size := c.DiceSize

	grayImg, err := imageHandler.GetGrayImageFromPath(path)
	if err != nil {
		panic(err)
	}

	imageHandler.SaveImage(grayImg, "gray")

	pixImg, err := imageHandler.PixelateImage(grayImg, blurSize)
	if err != nil {
		panic(err)
	}
	imageHandler.SaveImage(pixImg, "pix")

	diceImg, err := imageHandler.CreateDiceFromGrayImage(pixImg, size)
	if err != nil {
		panic(err)
	}
	imageHandler.SaveImage(diceImg, "dice")
}
