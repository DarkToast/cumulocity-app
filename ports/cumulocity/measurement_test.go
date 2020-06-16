package cumulocity

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestMeasurement(t *testing.T) {
	timeVal := time.Now()
	var measurement = Measurement{
		Source:  "9636292",
		Time:    timeVal,
		Type:    "TestMeasurement",
		Metrics: []Metric{Temperature{value: 23.5, unit: C}, Humidity{value: 64.55}},
	}

	t.Run("Json() returns the correct value", func(t *testing.T) {
		expected := `{
            "time" : "%s",
            "type" : "TestMeasurement",
            "source" : { "id": "9636292" },
            "c8y_TemperatureMeasurement":{"T":{"value":23.50,"unit":"C"}},
            "c8y_HumidityMeasurement":{"h":{"value":64.55,"unit":"%%RH"}}
        }`
		expected = fmt.Sprintf(expected, timeVal.Format(time.RFC3339))
		expected = strings.ReplaceAll(expected, "\n", "")
		expected = strings.ReplaceAll(expected, " ", "")
		given := measurement.Json()

		if given != expected {
			t.Errorf("Source() - \nexpected: \n%s \n\ngiven: \n%s", expected, given)
		}
	})
}

func TestTemperature(t *testing.T) {
	var temp = Temperature{
		value: 23.5,
		unit:  C,
	}

	t.Run("MetricJson() returns the correct value", func(t *testing.T) {
		expected := `"c8y_TemperatureMeasurement":{"T":{"value":23.50,"unit":"C"}}`
		given := temp.MetricJson()

		if given != expected {
			t.Errorf("Source() - expected: \n%s \n\n given: \n%s", expected, given)
		}
	})
}

func TestHumidity(t *testing.T) {
	var temp = Humidity{
		value: 64.55,
	}

	t.Run("MetricJson() returns the correct value", func(t *testing.T) {
		expected := `"c8y_HumidityMeasurement":{"h":{"value":64.55,"unit":"%RH"}}`
		given := temp.MetricJson()

		if given != expected {
			t.Errorf("Source() - expected: \n%s \n\n given: \n%s", expected, given)
		}
	})
}
