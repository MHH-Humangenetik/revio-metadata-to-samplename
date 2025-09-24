package parser

// PacBioDataModel represents the top-level XML structure.
type PacBioDataModel struct {
	ExperimentContainer ExperimentContainer `xml:"ExperimentContainer"`
}

// ExperimentContainer contains the runs.
type ExperimentContainer struct {
	Runs []Run `xml:"Runs>Run"`
}

// Run contains the outputs.
type Run struct {
	Outputs Outputs `xml:"Outputs"`
}

// Outputs contains the subread sets.
type Outputs struct {
	SubreadSets []SubreadSet `xml:"SubreadSets>SubreadSet"`
}

// SubreadSet contains the dataset metadata.
type SubreadSet struct {
	DataSetMetadata DataSetMetadata `xml:"DataSetMetadata"`
}

// DataSetMetadata contains the collections.
type DataSetMetadata struct {
	Collections []CollectionMetadata `xml:"Collections>CollectionMetadata"`
}

// CollectionMetadata contains the well sample and other details.
type CollectionMetadata struct {
	Context    string     `xml:"Context,attr"`
	WellSample WellSample `xml:"WellSample"`
	Primary    Primary    `xml:"Primary"`
}

// WellSample contains the bio samples.
type WellSample struct {
	BioSamples []BioSample `xml:"BioSamples>BioSample"`
}

// BioSample contains the name of the sample.
type BioSample struct {
	Name        string       `xml:"Name,attr"`
	DNABarcodes []DNABarcode `xml:"DNABarcodes>DNABarcode"`
}

// DNABarcode contains barcode information.
type DNABarcode struct {
	Name string `xml:"Name,attr"`
}

// Primary contains the output options.
type Primary struct {
	OutputOptions OutputOptions `xml:"OutputOptions"`
}

// OutputOptions contains the results folder.
type OutputOptions struct {
	ResultsFolder string `xml:"ResultsFolder"`
}
