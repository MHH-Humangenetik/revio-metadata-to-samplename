package cmd

import (
	"encoding/xml"
	"reflect"
	"testing"

	"revio-metadata-to-samplename/internal/parser"
)

func TestGenerateFilenames(t *testing.T) {
	tests := []struct {
		name     string
		xmlData  string
		expected []string
	}{
		{
			name: "Case 1: With barcodes - multiple samples",
			xmlData: `<?xml version="1.0" encoding="utf-8"?>
<PacBioDataModel xmlns="http://pacificbiosciences.com/PacBioDataModel.xsd">
  <ExperimentContainer>
    <Runs>
      <Run>
        <Outputs>
          <SubreadSets>
            <SubreadSet xmlns="http://pacificbiosciences.com/PacBioDatasets.xsd">
              <DataSetMetadata>
                <Collections xmlns="http://pacificbiosciences.com/PacBioCollectionMetadata.xsd">
                  <CollectionMetadata Context="m84297_250919_165138_s3">
                    <WellSample>
                      <BioSamples xmlns="http://pacificbiosciences.com/PacBioSampleInfo.xsd">
                        <BioSample Name="289313PR1_02">
                          <DNABarcodes>
                            <DNABarcode Name="bc2001--bc2001" />
                          </DNABarcodes>
                        </BioSample>
                        <BioSample Name="289314PR1_02">
                          <DNABarcodes>
                            <DNABarcode Name="bc2002--bc2002" />
                          </DNABarcodes>
                        </BioSample>
                      </BioSamples>
                    </WellSample>
                    <Primary>
                      <OutputOptions>
                        <ResultsFolder>r84297_20250919_123822/1_C01/</ResultsFolder>
                      </OutputOptions>
                    </Primary>
                  </CollectionMetadata>
                </Collections>
              </DataSetMetadata>
            </SubreadSet>
          </SubreadSets>
        </Outputs>
      </Run>
    </Runs>
  </ExperimentContainer>
</PacBioDataModel>`,
			expected: []string{
				`289313PR1_02,"r84297_20250919_123822/1_C01/hifi_reads/m84297_250919_165138_s3.hifi_reads.bc2001.bam"`,
				`289314PR1_02,"r84297_20250919_123822/1_C01/hifi_reads/m84297_250919_165138_s3.hifi_reads.bc2002.bam"`,
			},
		},
		{
			name: "Case 2: No barcodes - single sample",
			xmlData: `<?xml version="1.0" encoding="utf-8"?>
<PacBioDataModel xmlns="http://pacificbiosciences.com/PacBioDataModel.xsd">
  <ExperimentContainer>
    <Runs>
      <Run>
        <Outputs>
          <SubreadSets>
            <SubreadSet xmlns="http://pacificbiosciences.com/PacBioDatasets.xsd">
              <DataSetMetadata>
                <Collections xmlns="http://pacificbiosciences.com/PacBioCollectionMetadata.xsd">
                  <CollectionMetadata Context="m84297_250922_090411_s1">
                    <WellSample>
                      <BioSamples xmlns="http://pacificbiosciences.com/PacBioSampleInfo.xsd">
                        <BioSample Name="275351PR1_02">
                        </BioSample>
                      </BioSamples>
                    </WellSample>
                    <Primary>
                      <OutputOptions>
                        <ResultsFolder>r84297_20250922_085610/1_A01/</ResultsFolder>
                      </OutputOptions>
                    </Primary>
                  </CollectionMetadata>
                </Collections>
              </DataSetMetadata>
            </SubreadSet>
          </SubreadSets>
        </Outputs>
      </Run>
    </Runs>
  </ExperimentContainer>
</PacBioDataModel>`,
			expected: []string{
				`275351PR1_02,"r84297_20250922_085610/1_A01/hifi_reads/m84297_250922_090411_s1.hifi_reads.bam"`,
			},
		},
		{
			name: "Case 1: Complex barcode names with multiple dashes",
			xmlData: `<?xml version="1.0" encoding="utf-8"?>
<PacBioDataModel xmlns="http://pacificbiosciences.com/PacBioDataModel.xsd">
  <ExperimentContainer>
    <Runs>
      <Run>
        <Outputs>
          <SubreadSets>
            <SubreadSet xmlns="http://pacificbiosciences.com/PacBioDatasets.xsd">
              <DataSetMetadata>
                <Collections xmlns="http://pacificbiosciences.com/PacBioCollectionMetadata.xsd">
                  <CollectionMetadata Context="m12345_123456_789012_s1">
                    <WellSample>
                      <BioSamples xmlns="http://pacificbiosciences.com/PacBioSampleInfo.xsd">
                        <BioSample Name="TestSample1">
                          <DNABarcodes>
                            <DNABarcode Name="bc1001--bc1001--extra" />
                          </DNABarcodes>
                        </BioSample>
                      </BioSamples>
                    </WellSample>
                    <Primary>
                      <OutputOptions>
                        <ResultsFolder>test/results/folder/</ResultsFolder>
                      </OutputOptions>
                    </Primary>
                  </CollectionMetadata>
                </Collections>
              </DataSetMetadata>
            </SubreadSet>
          </SubreadSets>
        </Outputs>
      </Run>
    </Runs>
  </ExperimentContainer>
</PacBioDataModel>`,
			expected: []string{
				`TestSample1,"test/results/folder/hifi_reads/m12345_123456_789012_s1.hifi_reads.bc1001.bam"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pacBioData parser.PacBioDataModel
			if err := xml.Unmarshal([]byte(tt.xmlData), &pacBioData); err != nil {
				t.Fatalf("Failed to parse test XML: %v", err)
			}

			// Extract the collection metadata for testing
			collection := pacBioData.ExperimentContainer.Runs[0].Outputs.SubreadSets[0].DataSetMetadata.Collections[0]

			result := generateFilenames(collection)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("generateFilenames() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractSampleData(t *testing.T) {
	tests := []struct {
		name             string
		xmlData          string
		includeResults   bool
		includeFileNames bool
		expected         []string
	}{
		{
			name: "Default mode - sample names only",
			xmlData: `<?xml version="1.0" encoding="utf-8"?>
<PacBioDataModel xmlns="http://pacificbiosciences.com/PacBioDataModel.xsd">
  <ExperimentContainer>
    <Runs>
      <Run>
        <Outputs>
          <SubreadSets>
            <SubreadSet xmlns="http://pacificbiosciences.com/PacBioDatasets.xsd">
              <DataSetMetadata>
                <Collections xmlns="http://pacificbiosciences.com/PacBioCollectionMetadata.xsd">
                  <CollectionMetadata Context="m84297_250919_165138_s3">
                    <WellSample>
                      <BioSamples xmlns="http://pacificbiosciences.com/PacBioSampleInfo.xsd">
                        <BioSample Name="Sample1">
                          <DNABarcodes>
                            <DNABarcode Name="bc2001--bc2001" />
                          </DNABarcodes>
                        </BioSample>
                        <BioSample Name="Sample2">
                          <DNABarcodes>
                            <DNABarcode Name="bc2002--bc2002" />
                          </DNABarcodes>
                        </BioSample>
                      </BioSamples>
                    </WellSample>
                    <Primary>
                      <OutputOptions>
                        <ResultsFolder>results/</ResultsFolder>
                      </OutputOptions>
                    </Primary>
                  </CollectionMetadata>
                </Collections>
              </DataSetMetadata>
            </SubreadSet>
          </SubreadSets>
        </Outputs>
      </Run>
    </Runs>
  </ExperimentContainer>
</PacBioDataModel>`,
			includeResults:   false,
			includeFileNames: false,
			expected:         []string{"Sample1,Sample2"},
		},
		{
			name: "Results folder mode",
			xmlData: `<?xml version="1.0" encoding="utf-8"?>
<PacBioDataModel xmlns="http://pacificbiosciences.com/PacBioDataModel.xsd">
  <ExperimentContainer>
    <Runs>
      <Run>
        <Outputs>
          <SubreadSets>
            <SubreadSet xmlns="http://pacificbiosciences.com/PacBioDatasets.xsd">
              <DataSetMetadata>
                <Collections xmlns="http://pacificbiosciences.com/PacBioCollectionMetadata.xsd">
                  <CollectionMetadata Context="m84297_250919_165138_s3">
                    <WellSample>
                      <BioSamples xmlns="http://pacificbiosciences.com/PacBioSampleInfo.xsd">
                        <BioSample Name="Sample1">
                        </BioSample>
                      </BioSamples>
                    </WellSample>
                    <Primary>
                      <OutputOptions>
                        <ResultsFolder>results/folder/</ResultsFolder>
                      </OutputOptions>
                    </Primary>
                  </CollectionMetadata>
                </Collections>
              </DataSetMetadata>
            </SubreadSet>
          </SubreadSets>
        </Outputs>
      </Run>
    </Runs>
  </ExperimentContainer>
</PacBioDataModel>`,
			includeResults:   true,
			includeFileNames: false,
			expected:         []string{"Sample1;results/folder/"},
		},
		{
			name: "File names mode",
			xmlData: `<?xml version="1.0" encoding="utf-8"?>
<PacBioDataModel xmlns="http://pacificbiosciences.com/PacBioDataModel.xsd">
  <ExperimentContainer>
    <Runs>
      <Run>
        <Outputs>
          <SubreadSets>
            <SubreadSet xmlns="http://pacificbiosciences.com/PacBioDatasets.xsd">
              <DataSetMetadata>
                <Collections xmlns="http://pacificbiosciences.com/PacBioCollectionMetadata.xsd">
                  <CollectionMetadata Context="m84297_250919_165138_s3">
                    <WellSample>
                      <BioSamples xmlns="http://pacificbiosciences.com/PacBioSampleInfo.xsd">
                        <BioSample Name="TestSample">
                        </BioSample>
                      </BioSamples>
                    </WellSample>
                    <Primary>
                      <OutputOptions>
                        <ResultsFolder>test/results/</ResultsFolder>
                      </OutputOptions>
                    </Primary>
                  </CollectionMetadata>
                </Collections>
              </DataSetMetadata>
            </SubreadSet>
          </SubreadSets>
        </Outputs>
      </Run>
    </Runs>
  </ExperimentContainer>
</PacBioDataModel>`,
			includeResults:   false,
			includeFileNames: true,
			expected:         []string{`TestSample,"test/results/hifi_reads/m84297_250919_165138_s3.hifi_reads.bam"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pacBioData parser.PacBioDataModel
			if err := xml.Unmarshal([]byte(tt.xmlData), &pacBioData); err != nil {
				t.Fatalf("Failed to parse test XML: %v", err)
			}

			result := extractSampleData(&pacBioData, tt.includeResults, tt.includeFileNames)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("extractSampleData() = %v, want %v", result, tt.expected)
			}
		})
	}
}
