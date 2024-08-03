package tinyface

import (
	"errors"
	"fmt"
	goFace "github.com/Kagami/go-face"
	"image"
	"os"
)

var (
	ModelPathNotAvailableError = errors.New("the model is not available")

	NoFaceOnImageError         = errors.New("no face on the image")
	NotASingleFaceOnImageError = errors.New("not a single face on the image")
)

type Recognizer struct {
	tolerance     float32
	useCNN        bool
	useGray       bool
	samplesLoader SamplesLoader
	samplesSaver  SamplesSaver

	modelPath string

	rec         *goFace.Recognizer
	samplesData []*SampleBaseData
}

// SampleBaseData  descriptor of the human face.
type SampleBaseData struct {
	Id         string
	Descriptor goFace.Descriptor
}

// Face holds coordinates and descriptor of the human face.
type Face struct {
	*SampleBaseData
	Rectangle image.Rectangle
}

func NewRecognizer(opts ...Option) (*Recognizer, error) {
	// apply default options
	if len(opts) == 0 {
		opts = defaultOts()
	}
	recognizer := &Recognizer{}
	for _, opt := range opts {
		opt(recognizer)
	}
	if recognizer.modelPath == "" {
		return nil, ModelPathNotAvailableError
	}
	rec, err := goFace.NewRecognizer(recognizer.modelPath)
	if err != nil {
		return nil, err
	}
	recognizer.rec = rec

	return recognizer, nil
}

func getSamplesFrom(baseDatas []*SampleBaseData) (samples []goFace.Descriptor, cats []int32) {
	samples = make([]goFace.Descriptor, len(baseDatas))
	cats = make([]int32, len(baseDatas))
	for i, f := range baseDatas {
		samples = append(samples, f.Descriptor)
		cats = append(cats, int32(i))
	}
	return samples, cats
}

// LoadSamples  load samples using config loader
func (r *Recognizer) LoadSamples(samplesPath ...string) error {
	if r.samplesLoader != nil {
		sampleBaseData, err := r.samplesLoader.LoadSamples(samplesPath...)
		if err != nil {
			return err
		}
		samples, cats := getSamplesFrom(sampleBaseData)
		r.rec.SetSamples(samples, cats)
	}
	return nil
}

// SaveSamples  save samples using config saver
func (r *Recognizer) SaveSamples(samplesPath ...string) error {
	if r.samplesSaver != nil {
		if err := r.samplesSaver.SaveSamples(r.samplesData, samplesPath...); err != nil {
			return err
		}
	}
	return nil
}

// Close frees resources taken by the Recognizer. Safe to call multiple
// times. Don't use Recognizer after close call.
func (r *Recognizer) Close() {
	r.rec.Close()

}

// AddImageToSamples add a sample image to the samples
func (r *Recognizer) AddImageToSamples(Path string, Id string) error {

	file := Path
	var err error

	if r.useCNN {

		file, err = r.createTempGrayFile(file, Id)

		if err != nil {
			return err
		}
		defer os.Remove(file)
	}

	var faces []goFace.Face

	if r.useCNN {
		faces, err = r.rec.RecognizeFileCNN(file)
	} else {
		faces, err = r.rec.RecognizeFile(file)
	}

	if err != nil {
		return err
	}

	if len(faces) == 0 {
		return NoFaceOnImageError
	}

	if len(faces) > 1 {
		return NotASingleFaceOnImageError
	}

	f := &SampleBaseData{}
	f.Id = Id
	f.Descriptor = faces[0].Descriptor

	r.samplesData = append(r.samplesData, f)

	return nil

}

// RecognizeSingle returns face if it's the only face on the image or nil otherwise.
// Only JPEG format is currently supported.
func (r *Recognizer) RecognizeSingle(Path string) (*goFace.Face, error) {

	file := Path
	var err error
	if r.useGray {
		file, err = r.createTempGrayFile(file, "64ab59ac42d69274f06eadb11348969e")
		if err != nil {
			return nil, err
		}
		defer os.Remove(file)
	}

	var idFace *goFace.Face

	if r.useCNN {
		idFace, err = r.rec.RecognizeSingleFileCNN(file)
	} else {
		idFace, err = r.rec.RecognizeSingleFile(file)
	}

	if err != nil {
		return nil, fmt.Errorf("can't recognize: %v", err)

	}
	if idFace == nil {
		return nil, NotASingleFaceOnImageError
	}

	return idFace, nil

}

// RecognizeMultiples returns all faces found on the provided image, sorted from
// left to right. Empty list is returned if there are no faces, error is
// returned if there was some error while decoding/processing image.
// Only JPEG format is currently supported.
func (r *Recognizer) RecognizeMultiples(Path string) ([]goFace.Face, error) {

	file := Path
	var err error

	if r.useGray {
		file, err = r.createTempGrayFile(file, "64ab59ac42d69274f06eadb11348969e")

		if err != nil {
			return nil, err
		}

		defer os.Remove(file)

	}

	var idFaces []goFace.Face

	if r.useCNN {
		idFaces, err = r.rec.RecognizeFileCNN(file)
	} else {
		idFaces, err = r.rec.RecognizeFile(file)
	}

	if err != nil {
		return nil, fmt.Errorf("can't recognize: %v", err)
	}

	return idFaces, nil

}

// Classify returns all faces identified in the image. Empty list is returned if no match.
func (r *Recognizer) Classify(Path string) ([]*Face, error) {

	face, err := r.RecognizeSingle(Path)

	if err != nil {
		return nil, err
	}

	personID := r.rec.ClassifyThreshold(face.Descriptor, r.tolerance)
	if personID < 0 {
		return nil, fmt.Errorf("can't classify")
	}

	facesRec := make([]*Face, 0)
	aux := &Face{SampleBaseData: r.samplesData[personID], Rectangle: face.Rectangle}
	facesRec = append(facesRec, aux)

	return facesRec, nil

}

// ClassifyMultiples returns all faces identified in the image. Empty list is returned if no match.
func (r *Recognizer) ClassifyMultiples(Path string) ([]Face, error) {

	faces, err := r.RecognizeMultiples(Path)

	if err != nil {
		return nil, fmt.Errorf("can't recognize: %v", err)
	}

	facesRec := make([]Face, 0)

	for _, f := range faces {

		personID := r.rec.ClassifyThreshold(f.Descriptor, r.tolerance)
		if personID < 0 {
			continue
		}

		aux := Face{SampleBaseData: r.samplesData[personID], Rectangle: f.Rectangle}

		facesRec = append(facesRec, aux)

	}

	return facesRec, nil

}
