package utils

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/chmenegatti/nsxt-vs/api"
)

func SortLoadBalancesByIP(results [][3]string) {
	sort.Slice(
		results, func(i, j int) bool {
			ipI := strings.Split(results[i][1], "-")[0]
			ipJ := strings.Split(results[j][1], "-")[0]

			partsI := strings.Split(ipI, ".")
			partsJ := strings.Split(ipJ, ".")

			for k := 0; k < 4; k++ {
				octI, _ := strconv.Atoi(partsI[k])
				octJ, _ := strconv.Atoi(partsJ[k])

				if octI != octJ {
					return octI < octJ
				}
			}
			return false
		},
	)
}

func CompareIPPort(a, b api.VirtualServer) bool {
	partsA := strings.Split(a.DisplayName, "-")
	partsB := strings.Split(b.DisplayName, "-")

	if partsA[0] == partsB[0] {
		return partsA[1] < partsB[1]
	}
	return partsA[0] < partsB[0]
}

func GetCSVFilePath(filename string) (string, error) {
	dir := "csv" // Verifica se a pasta existe
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, filename), nil
}
