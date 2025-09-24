package cmd

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"revio-metadata-to-samplename/internal/parser"
)

var (
	inputFile     string
	resultsFolder bool
	fileNames     bool
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converts Revio metadata to a sample name",
	Long: `Converts Revio metadata from an XML file to a sample name.
You can also extract the results folder.`,
	Run: runConvert,
}

func runConvert(cmd *cobra.Command, args []string) {
	data, err := parseMetadataFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to parse metadata file: %v", err)
	}

	output := extractSampleData(data, resultsFolder, fileNames)
	fmt.Println(strings.Join(output, "\n"))
}

func parseMetadataFile(filename string) (*parser.PacBioDataModel, error) {
	xmlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", filename, err)
	}

	var pacBioData parser.PacBioDataModel
	if err := xml.Unmarshal(xmlFile, &pacBioData); err != nil {
		return nil, fmt.Errorf("unmarshalling XML: %w", err)
	}

	return &pacBioData, nil
}

func extractSampleData(data *parser.PacBioDataModel, includeResultsFolder bool, includeFileNames bool) []string {
	var output []string

	for _, run := range data.ExperimentContainer.Runs {
		for _, subreadSet := range run.Outputs.SubreadSets {
			for _, collection := range subreadSet.DataSetMetadata.Collections {
				if includeFileNames {
					filenames := generateFilenames(collection)
					output = append(output, filenames...)
				} else {
					sampleNames := extractSampleNames(collection.WellSample.BioSamples)

					if includeResultsFolder {
						resultsFolder := collection.Primary.OutputOptions.ResultsFolder
						for _, sampleName := range sampleNames {
							output = append(output, fmt.Sprintf("%s;%s", sampleName, resultsFolder))
						}
					} else {
						output = append(output, strings.Join(sampleNames, ","))
					}
				}
			}
		}
	}

	return output
}

func extractSampleNames(bioSamples []parser.BioSample) []string {
	sampleNames := make([]string, 0, len(bioSamples))
	for _, bioSample := range bioSamples {
		sampleNames = append(sampleNames, bioSample.Name)
	}
	return sampleNames
}

func generateFilenames(collection parser.CollectionMetadata) []string {
	var filenames []string
	resultsFolder := collection.Primary.OutputOptions.ResultsFolder
	context := collection.Context

	// Check if any bio samples have barcodes
	hasBarcodes := false
	for _, bioSample := range collection.WellSample.BioSamples {
		if len(bioSample.DNABarcodes) > 0 {
			hasBarcodes = true
			break
		}
	}

	if hasBarcodes {
		// Case 1: Barcodes are used (multiple samples)
		for _, bioSample := range collection.WellSample.BioSamples {
			if len(bioSample.DNABarcodes) > 0 {
				for _, barcode := range bioSample.DNABarcodes {
					// Cut barcode before the first '-'
					barcodeShort := strings.Split(barcode.Name, "-")[0]
					filename := fmt.Sprintf("\"%shifi_reads/%s.hifi_reads.%s.bam\"", resultsFolder, context, barcodeShort)
					output := fmt.Sprintf("%s,%s", bioSample.Name, filename)
					filenames = append(filenames, output)
				}
			}
		}
	} else {
		// Case 2: No barcodes are used (one sample)
		filename := fmt.Sprintf("\"%shifi_reads/%s.hifi_reads.bam\"", resultsFolder, context)
		// Use the first (and only) sample name
		sampleName := collection.WellSample.BioSamples[0].Name
		output := fmt.Sprintf("%s,%s", sampleName, filename)
		filenames = append(filenames, output)
	}

	return filenames
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input XML file (required)")
	convertCmd.MarkFlagRequired("input")
	convertCmd.Flags().BoolVar(&resultsFolder, "resultsFolder", false, "Include results folder in output")
	convertCmd.Flags().BoolVar(&fileNames, "fileNames", false, "Generate filenames based on XML metadata")
}
