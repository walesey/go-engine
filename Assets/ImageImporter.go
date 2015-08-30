package assets

import (
	"os"
	"log"
	"image"
	"math"
	_ "image/png"
	_ "image/jpeg"

	"github.com/disintegration/imaging"
	"hawx.me/code/img/blur"
)

type CubeMapData struct {
	Right, Left, Top, Bottom, Back, Front image.Image
}

func ImportImage( file string ) image.Image {
	imgFile, err := os.Open(file)
	if err != nil {
	    log.Fatal(err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
	    log.Fatal(err)
	}
	img = imaging.FlipV(img)
	return img
}

func ImportCubemap( file string ) *CubeMapData {
	cubeMap := new(CubeMapData)

	baseImage := ImportImage(file)
	x := baseImage.Bounds().Max.X 
	y := baseImage.Bounds().Max.Y

	cubeMap.Right = imaging.Crop( baseImage, image.Rect(x/2, y/3, 3*x/4, 2*y/3) )
	cubeMap.Left = imaging.Crop( baseImage, image.Rect(0, y/3, x/4, 2*y/3) )
	cubeMap.Top = imaging.Crop( baseImage, image.Rect(x/4, 0, x/2, y/3) )
	cubeMap.Bottom = imaging.Crop( baseImage, image.Rect(x/4, 2*y/3, x/2, y) )
	cubeMap.Back = imaging.Crop( baseImage, image.Rect(3*x/4, y/3, x, 2*y/3) )
	cubeMap.Front = imaging.Crop( baseImage, image.Rect(x/4, y/3, x/2, 2*y/3) )

	return cubeMap
}

func (cm *CubeMapData) Clone() *CubeMapData{
	return &CubeMapData{Right: imaging.Clone(cm.Right), Left: imaging.Clone(cm.Left), Top: imaging.Clone(cm.Top), Bottom: imaging.Clone(cm.Bottom), Back: imaging.Clone(cm.Back), Front: imaging.Clone(cm.Front)}
}

func (cm *CubeMapData) Resize(size int) {
	cm.Front = imaging.Resize(cm.Front, size, size, imaging.Gaussian)
	cm.Back = imaging.Resize(cm.Back, size, size, imaging.Gaussian)
	cm.Bottom = imaging.Resize(cm.Bottom, size, size, imaging.Gaussian)
	cm.Top = imaging.Resize(cm.Top, size, size, imaging.Gaussian)
	cm.Left = imaging.Resize(cm.Left, size, size, imaging.Gaussian)
	cm.Right = imaging.Resize(cm.Right, size, size, imaging.Gaussian)
}

func (cm *CubeMapData) Blur(sigma float64) {
	cmCopy := cm.Clone()
	cm.Front = blurCubeSide( cmCopy.Front, cmCopy.Top, cmCopy.Bottom, cmCopy.Left, cmCopy.Right, sigma )
	cm.Back = blurCubeSide( cmCopy.Back, imaging.Rotate180(cmCopy.Top), imaging.Rotate180(cmCopy.Bottom), cmCopy.Right, cmCopy.Left, sigma )
	cm.Bottom = blurCubeSide( cmCopy.Bottom, cmCopy.Front, imaging.Rotate180(cmCopy.Back), imaging.Rotate90(cmCopy.Left), imaging.Rotate270(cmCopy.Right), sigma )
	cm.Top = blurCubeSide( cmCopy.Top, imaging.Rotate180(cmCopy.Back), cmCopy.Front, imaging.Rotate270(cmCopy.Left), imaging.Rotate90(cmCopy.Right), sigma )
	cm.Left = blurCubeSide( cmCopy.Left, imaging.Rotate90(cmCopy.Top), imaging.Rotate270(cmCopy.Bottom), cmCopy.Back, cmCopy.Front, sigma )
	cm.Right = blurCubeSide( cmCopy.Right, imaging.Rotate270(cmCopy.Top), imaging.Rotate90(cmCopy.Bottom), cmCopy.Front, cmCopy.Back, sigma )
}

func blurCubeSide( img, north, south, east, west image.Image, sigma float64 ) image.Image {
	x := img.Bounds().Max.X 
	y := img.Bounds().Max.Y
	var tempImg image.Image
	tempImg = image.NewNRGBA(image.Rect(0, 0, 3*x, 3*y))
	tempImg = imaging.Paste(tempImg, img, image.Pt(x, y))
	tempImg = imaging.Paste(tempImg, north, image.Pt(x, 0)) //n
	tempImg = imaging.Paste(tempImg, south, image.Pt(x, 2*y)) //s
	tempImg = VerticalBlur(tempImg, sigma)
	tempImg = imaging.Paste(tempImg, east, image.Pt(0, y)) //e
	tempImg = imaging.Paste(tempImg, west, image.Pt(2*x, y)) //w
	tempImg = HorizontalBlur(tempImg, sigma)
	return imaging.Crop(tempImg, image.Rect(x, y, 2*x, 2*y))
}

func HorizontalBlur(in image.Image, sigma float64) image.Image {
	f := func(n int) float64 {
		return math.Exp(-float64(n*n) / (2 * sigma * sigma))
	}
	wide := blur.NewHorizontalKernel(int(sigma*3)*2+1, f).Normalised()
	return blur.Convolve(in, wide, blur.WRAP)
}

func VerticalBlur(in image.Image, sigma float64) image.Image {
	f := func(n int) float64 {
		return math.Exp(-float64(n*n) / (2 * sigma * sigma))
	}
	tall := blur.NewVerticalKernel(int(sigma*3)*2+1, f).Normalised()
	return blur.Convolve(in, tall, blur.WRAP)
}