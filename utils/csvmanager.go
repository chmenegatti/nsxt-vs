package utils

import (
	"encoding/csv"
	"os"

	"go.uber.org/zap"
)

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
