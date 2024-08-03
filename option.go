package tinyface

const (
	defaultTolerance = float32(0.6)
)

type Option func(recognizer *Recognizer)

// WithModelPath sets the path to the model directory.
func WithModelPath(modelPath string) Option {
	return func(recognizer *Recognizer) {
		recognizer.modelPath = modelPath
	}
}

// WithTolerance sets the tolerance for the recognition.
func WithTolerance(tolerance float32) Option {
	return func(recognizer *Recognizer) {
		recognizer.tolerance = tolerance
	}
}

// UseCNN  set whether to use CNN or not.
func UseCNN(useCNN bool) Option {
	return func(recognizer *Recognizer) {
		recognizer.useCNN = useCNN
	}
}

// UseGray set whether to use grayscale or not.
func UseGray(useGray bool) Option {
	return func(recognizer *Recognizer) {
		recognizer.useGray = useGray
	}
}

// WithSamplesLoader  sets the samples loader.
func WithSamplesLoader(loader SamplesLoader) Option {
	return func(recognizer *Recognizer) {
		recognizer.samplesLoader = loader
	}
}

// WithSamplesSaver  sets the  dataSet store method
func WithSamplesSaver(saver SamplesSaver) Option {
	return func(recognizer *Recognizer) {
		recognizer.samplesSaver = saver
	}
}

func defaultOts() []Option {
	return []Option{
		WithTolerance(defaultTolerance),
		UseCNN(false),
		UseGray(true),
	}
}
