package tinyface

// SamplesLoader is the interface for loading a dataset.
type SamplesLoader interface {
	LoadSamples(dataSetPath ...string) ([]*SampleBaseData, error)
}

// SamplesSaver is the interface for storing a dataset.
type SamplesSaver interface {
	SaveSamples(dataSet []*SampleBaseData, dataSetPath ...string) error
}
