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
	metricName = metricCSVPath(metricName)
	if err := os.MkdirAll(filepath.Dir(metricName), 0755); err != nil {
		return fmt.Errorf("%s mkdir fail: %w", filepath.Dir(metricName), err)
	}

	if _, err := os.Stat(metricName); os.IsNotExist(err) {
		f, err := os.Create(metricName)
		if err != nil {
			return fmt.Errorf("%s create fail: %w", metricName, err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("%s close fail: %w", metricName, err)
		}
		header := []string{"date", "value"}
		if err := WriteCSV(metricName, header); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("%s stat fail: %w", metricName, err)
	}
	return nil
}

func WriteCSV(metric string, record []string) error {
	metric = metricCSVPath(metric)
	if err := os.MkdirAll(filepath.Dir(metric), 0755); err != nil {
		return fmt.Errorf("%s mkdir fail: %w", filepath.Dir(metric), err)
	}

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
	if err := writer.Error(); err != nil {
		return errors.Wrap(err, "csv flush")
	}

	return nil
}

// date, value  ==> 2024-04-16, , 45.21
func ReadCSV(metric string) error {
	metric = metricCSVPath(metric)
	f, err := os.Open(metric) // os.Open => O_RDONLY  Read Only
	if err != nil {
		return fmt.Errorf("%s open error: %s", metric, err.Error())
	}
	defer f.Close()

	reader := csv.NewReader(f)
	//
	record, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "csv write")
	}
	fmt.Println("record: ", record)

	return nil
}

func metricCSVPath(metric string) string {
	if filepath.Dir(metric) != "." {
		return metric
	}
	return filepath.Join("store", metric)
}
