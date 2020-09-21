package reloadly

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithNumberAsStringJson(t *testing.T) {
	input := `{"number":"+34987","amount":100,"country":"ES"}`

	var job TopupJob
	err := json.Unmarshal([]byte(input), &job)
	assert.Nil(t, err)

	assert.Equal(t, "+34987", job.Number)
	assert.Equal(t, "ES", job.Country)

	b, err := json.Marshal(job)
	assert.Nil(t, err)
	assert.Equal(t, input, string(b))
}

func TestWithNumberAsNumberJson(t *testing.T) {
	input := `{"number":34987987,"amount":100,"country":"ES"}`
	output := `{"number":"34987987","amount":100,"country":"ES"}`

	job := new(TopupJob)
	err := json.Unmarshal([]byte(input), job)
	assert.Nil(t, err)

	b, err := json.Marshal(job)
	assert.Nil(t, err)
	assert.Equal(t, output, string(b))
}

func TestWithNumberWithIDString(t *testing.T) {
	input := `{"number":34987987,"amount":100,"country":"ES","id":"foo"}`
	output := `{"number":"34987987","amount":100,"country":"ES","id":"foo"}`

	job := new(TopupJob)
	err := json.Unmarshal([]byte(input), job)
	assert.Nil(t, err)

	b, err := json.Marshal(job)
	assert.Nil(t, err)
	assert.Equal(t, output, string(b))
}

func TestWithNumberWithIDFloat(t *testing.T) {
	input := `{"number":34987987,"amount":100,"country":"ES","id":500.5}`
	output := `{"number":"34987987","amount":100,"country":"ES","id":"500.5"}`

	job := new(TopupJob)
	err := json.Unmarshal([]byte(input), job)
	assert.Nil(t, err)

	b, err := json.Marshal(job)
	assert.Nil(t, err)
	assert.Equal(t, output, string(b))
}
