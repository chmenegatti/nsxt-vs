package csv

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVWriter struct{}

func NewCSVWriter() *CSVWriter {
	return &CSVWriter{}
}

func (w *CSVWriter) WriteToFile(data [][3]string, filePath string, headers []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if len(headers) > 0 {
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write headers: %w", err)
		}
	}

	for _, record := range data {
		if err := writer.Write(record[:]); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func CompareCSVFiles(file1Path, file2Path, outputPath string) error {
	records1, err := readCSVFile(file1Path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file1Path, err)
	}

	records2, err := readCSVFile(file2Path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file2Path, err)
	}

	ids1 := make(map[string]bool)
	for _, record := range records1 {
		if len(record) > 0 {
			ids1[record[0]] = true
		}
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", outputPath, err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(records2[0]); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for _, record := range records2[1:] {
		if len(record) > 0 {
			if !ids1[record[0]] {
				if err := writer.Write(record); err != nil {
					return fmt.Errorf("failed to write record: %w", err)
				}
			}
		}
	}

	return nil
}

func readCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}
