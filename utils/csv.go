package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// header : daily, now, usage

const Metric_CPU_CSV = "metric_cpu.csv"
const Metric_NETWORK_CSV = "metric_network.csv"

func InitializeCSV(metricName string) error {
	metricName = filepath.Join("csv", metricName)
	if _, err := os.Stat(metricName); os.IsNotExist(err) {
		f, err := os.Create(metricName)
		if err != nil {
			return fmt.Errorf("%s create fail", metricName)
		}
		header := []string{"date", "value"}
		WriteCSV(metricName, header)
		f.Close()
	}
	return nil
}

func WriteCSV(metric string, record []string) error {
	f, err := os.OpenFile(metric, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("%s open error: %s", metric, err.Error())
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	if err = writer.Write(record); err != nil {
		return errors.Wrap(err, "csv write")
	}
	writer.Flush()

	return nil
}

// date, value  ==> 2024-04-16, , 45.21
func ReadCSV(metric string) error {
	f, err := os.Open(metric) // os.Open => O_RDONLY  Read Only
	if err != nil {
		return fmt.Errorf("%s open error: %s", metric, err.Error())
	}

	reader := csv.NewReader(f)
	//
	record, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "csv write")
	}
	fmt.Println("record: ", record)

	return nil
}
