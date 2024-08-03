package tinyface

import (
	"crypto/rand"
	"encoding/hex"
	goFace "github.com/Kagami/go-face"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
)

/*
LoadImage Load an image from file
*/
func (r *Recognizer) LoadImage(Path string) (image.Image, error) {

	existingImageFile, err := os.Open(Path)
	if err != nil {
		return nil, err
	}
	defer existingImageFile.Close()

	imageData, _, err := image.Decode(existingImageFile)
	if err != nil {
		return nil, err
	}

	return imageData, nil

}

/*
SaveToJpeg Save an image to jpeg file
*/
func (r *Recognizer) SaveToJpeg(Path string, img image.Image) error {

	imgWriter, err := os.Create(Path)
	if err != nil {
		return err
	}
	defer imgWriter.Close()

	err = jpeg.Encode(imgWriter, img, nil)

	if err != nil {
		return err
	}

	return nil

}

/*
GrayScale Convert an image to grayscale
*/
func (r *Recognizer) GrayScale(imgSrc image.Image) image.Image {

	return imaging.Grayscale(imgSrc)

}

/*
createTempGrayFile create a temporary image in grayscale
*/
func (r *Recognizer) createTempGrayFile(Path, Id string) (string, error) {

	name := r.tempFileName(Id, ".jpeg")

	img, err := r.LoadImage(Path)

	if err != nil {
		return "", err
	}

	img = r.GrayScale(img)
	err = r.SaveToJpeg(name, img)

	if err != nil {
		return "", err
	}

	return name, nil

}

// tempFileName generates a temporary filename
func (r *Recognizer) tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

/*
DrawFaces draws the faces identified in the original image
*/
func (r *Recognizer) DrawFaces(Path string, F []Face) (image.Image, error) {

	img, err := r.LoadImage(Path)

	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 24})

	dc := gg.NewContextForImage(img)
	dc.SetFontFace(face)

	for _, f := range F {

		dc.SetRGB255(0, 0, 255)

		x := float64(f.Rectangle.Min.X)
		y := float64(f.Rectangle.Min.Y)
		w := float64(f.Rectangle.Dx())
		h := float64(f.Rectangle.Dy())

		dc.DrawString(f.Id, x, y+h+20)

		dc.DrawRectangle(x, y, w, h)
		dc.SetLineWidth(4.0)
		dc.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 0, G: 0, B: 255, A: 255}))
		dc.Stroke()

	}

	img = dc.Image()

	return img, nil

}

/*
DrawFaces2 draws the faces in the original image
*/
func (r *Recognizer) DrawFaces2(Path string, F []goFace.Face) (image.Image, error) {

	aux := make([]Face, 0)

	for _, f := range F {

		auxFace := Face{}
		auxFace.Rectangle = f.Rectangle
		auxFace.Descriptor = f.Descriptor

		aux = append(aux, auxFace)

	}

	return r.DrawFaces(Path, aux)

}
