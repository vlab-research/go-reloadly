package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadBatchCsvLoadsFullCsv(t *testing.T) {
	deets, err := LoadBatchCsv("test/batch-full.csv")
	assert.Nil(t, err)
	assert.Equal(t, float64(100), deets[0].Amount)
	assert.Equal(t, 2.5, deets[1].Amount)
	assert.Equal(t, "foo", deets[0].Number)
	assert.Equal(t, "bar", deets[1].Number)

	assert.Equal(t, "Bardafone", deets[0].Operator)
	assert.Equal(t, "Foodafone India", deets[1].Operator)
}


func TestLoadBatchCsvLoadsCsvWithMissingToleranceOrOperators(t *testing.T) {
	deets, err := LoadBatchCsv("test/batch-missing-rows.csv")
	assert.Nil(t, err)
	assert.Equal(t, float64(100), deets[0].Amount)
	assert.Equal(t, 2.5, deets[1].Amount)

	assert.Equal(t, "", deets[0].Operator)
	assert.Equal(t, "Foodafone India", deets[1].Operator)
}


func TestLoadBatchCsvLoadsCsvWithMissingToleranceOrOperatorsInHeader(t *testing.T) {
	deets, err := LoadBatchCsv("test/batch-missing-columns.csv")
	assert.Nil(t, err)

	assert.Equal(t, "foo", deets[0].Number)
	assert.Equal(t, "bar", deets[1].Number)
	assert.Equal(t, float64(100), deets[0].Amount)
	assert.Equal(t, 2.5, deets[1].Amount)
}


func TestLoadBatchCsvErrorsWhenCsvMissingRequiredDetails(t *testing.T) {
	_, err := LoadBatchCsv("test/batch-missing-required.csv")
	assert.NotNil(t, err)
}
