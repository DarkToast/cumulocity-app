package cumulocity

import (
	"fmt"
	"time"
)

type Measurement struct {
	Source  Id
	Time    time.Time
	Type    string
	Metrics []Metric
}

func (measurement *Measurement) Json() string {
	json := fmt.Sprintf(`{"time":"%s","type":"%s","source":{"id":"%s"},`,
		measurement.Time.Format(time.RFC3339), measurement.Type, measurement.Source)

	for i, metric := range measurement.Metrics {
		json = json + metric.MetricJson()
		if i < len(measurement.Metrics)-1 {
			json = json + ","
		}
	}

	return json + "}"
}

// ---- METRICS ----

type Metric interface {
	MetricJson() string
}

// ---- Temperature
type TempUnit int

const (
	C TempUnit = iota
	F
)

func (t TempUnit) String() string {
	return [...]string{"C", "F"}[t]
}

type Temperature struct {
	value float64
	unit  TempUnit
}

func (temp Temperature) MetricJson() string {
	template := `"c8y_TemperatureMeasurement":{"T":{"value":%.2f,"unit":"%s"}}`
	return fmt.Sprintf(template, temp.value, temp.unit.String())
}

// ---- Humidity
type Humidity struct {
	value float64
}

func (temp Humidity) MetricJson() string {
	template := `"c8y_HumidityMeasurement":{"h":{"value":%.2f,"unit":"%%RH"}}`
	return fmt.Sprintf(template, temp.value)
}
