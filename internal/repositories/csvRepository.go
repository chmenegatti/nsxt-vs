package repositories

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/chmenegatti/nsxt-vs/cmd"
	"github.com/chmenegatti/nsxt-vs/config"
	"github.com/chmenegatti/nsxt-vs/internal/models"
	"github.com/chmenegatti/nsxt-vs/utils"
)

type CSVRepository struct {
	filename string
}

func NewCSVRepository(filename string) *CSVRepository {
	return &CSVRepository{filename: filename}
}

func (c *CSVRepository) LoadCSVData(edge string) ([]models.RegisterCSV, error) {
	var (
		filepath string
		err      error
	)

	if filepath, err = utils.GetCSVFilePath(fmt.Sprintf("%s-%s", edge, c.filename)); err != nil {
		return nil, fmt.Errorf("error getting CSV file path: %v", err)
	}
	file, err := os.Open(filepath)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	reader := csv.NewReader(file)
	registers, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var csvRegisters []models.RegisterCSV
	if len(registers) > 0 {
		for _, row := range registers[1:] {
			csvRegisters = append(
				csvRegisters, models.RegisterCSV{
					ID:          row[0],
					DisplayName: row[1],
					Service:     row[2],
					ClientCode:  row[3],
				},
			)
		}
	}

	return csvRegisters, nil
}

func (c *CSVRepository) SaveCSVData(registers []models.RegisterCSV, edge string) error {
	var (
		filepath string
		err      error
	)

	if filepath, err = utils.GetCSVFilePath(fmt.Sprintf("%s-%s", edge, c.filename)); err != nil {
		return fmt.Errorf("error getting CSV file path: %v", err)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var rows [][]string
	rows = append(rows, []string{"id", "display_name", "service", "client_code"})
	for _, register := range registers {
		rows = append(rows, []string{register.ID, register.DisplayName, register.Service, register.ClientCode})
	}

	if err := writer.WriteAll(rows); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func (c *CSVRepository) GetCSVData(edge string) error {
	configuration, err := config.LoadConfig("config.yaml")
	if err != nil {
		return fmt.Errorf("could not load config: %v", err)
	}

	token := configuration.GetToken()

	cmd.Run(configuration, edge, token)
	return nil
}
