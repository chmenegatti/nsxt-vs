package csvapi

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func WriteCSV(data [][3]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"id", "display_name", "service"}); err != nil {
		return err
	}

	for _, record := range data {
		if err := writer.Write(record[:]); err != nil {
			return err
		}
	}
	return nil
}

func CompareCSVFiles(file1Path, file2Path, outputPath string) error {
	records1, err := ReadCSVFile(file1Path)
	if err != nil {
		return err
	}

	records2, err := ReadCSVFile(file2Path)
	if err != nil {
		return err
	}

	ids1 := makeIDMap(records1)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", outputPath, err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if len(records2) > 0 {
		if err := writer.Write(records2[0]); err != nil {
			return fmt.Errorf("error writing header to %s: %v", outputPath, err)
		}
	}

	return writeDiffRecords(writer, records2, ids1)
}

func ReadCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func makeIDMap(records [][]string) map[string]bool {
	ids := make(map[string]bool)
	for i, record := range records {
		if i == 0 || len(record) == 0 {
			continue
		}
		ids[strings.TrimSpace(record[0])] = true
	}
	return ids
}

func writeDiffRecords(writer *csv.Writer, records [][]string, ids1 map[string]bool) error {
	for i, record := range records {
		if i == 0 || len(record) == 0 {
			continue
		}
		id := strings.TrimSpace(record[0])
		if !ids1[id] {
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing record: %v", err)
			}
		}
	}
	return nil
}
