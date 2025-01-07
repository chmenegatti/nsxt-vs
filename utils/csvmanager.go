package utils

import (
	"encoding/csv"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func ReadFromCSV(filename string, logger *zap.Logger) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		logger.Error("Failed to open CSV file", zap.Error(err))
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		logger.Error("Failed to read CSV file", zap.Error(err))
		return nil, err
	}

	logger.Info("Data successfully read from CSV file")
	return records, nil
}

func WriteToCSV(filename string, data [][]string, logger *zap.Logger) error {
	file, err := os.Create(filename)
	if err != nil {
		logger.Error("Failed to create CSV file", zap.Error(err))
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			logger.Error("Failed to write record to CSV", zap.Error(err))
			return err
		}
	}

	logger.Info("Data successfully saved to CSV file")
	return nil
}

// crie uma função que compara 2 csv e gera o resultado dessa comparação. O campo a ser compardo é o id do nsxt.csv com o campo nsxt_virtual_server_id do vs.csv devolver somente o que não existe em nsxt.csv
func CompareCSV(nsxtFile string, vsFile string, logger *zap.Logger) (int, error) {
	nsxt, err := os.Open(nsxtFile)
	if err != nil {
		logger.Error("Failed to open NSXT CSV file", zap.Error(err))
		return 0, err
	}
	defer func(nsxt *os.File) {
		err := nsxt.Close()
		if err != nil {

		}
	}(nsxt)

	vs, err := os.Open(vsFile)
	if err != nil {
		logger.Error("Failed to open VS CSV file", zap.Error(err))
		return 0, err
	}
	defer func(vs *os.File) {
		err := vs.Close()
		if err != nil {

		}
	}(vs)

	nsxtReader := csv.NewReader(nsxt)
	vsReader := csv.NewReader(vs)

	nsxtRecords, err := nsxtReader.ReadAll()
	if err != nil {
		logger.Error("Failed to read NSXT CSV file", zap.Error(err))
		return 0, err
	}

	vsRecords, err := vsReader.ReadAll()
	if err != nil {
		logger.Error("Failed to read VS CSV file", zap.Error(err))
		return 0, err
	}

	var diff [][]string

	for _, nsxt := range nsxtRecords {
		found := false
		for _, vs := range vsRecords {
			if nsxt[0] == vs[2] {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, nsxt)
		}

	}

	numDif := len(diff)

	if err := WriteToCSV("diff.csv", diff, logger); err != nil {
		logger.Error("Failed to write diff to CSV", zap.Error(err))
		return 0, err
	}

	logger.Info(fmt.Sprintf("Data successfully saved to CSV file - %d - %s", len(diff), "diff.csv"))
	return numDif, nil
}
