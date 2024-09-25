package reloadly

import (
	"fmt"

	"github.com/vlab-research/gotils"
)

type TopupWorker Service

func GimmeString(i interface{}) (interface{}, error) {
	switch i.(type) {

	case nil:
		return nil, nil

	case float64:

		// Default type for json numbers.
		// Check if can be integer
		f := i.(float64)
		if float64(int64(f)) == f {
			return fmt.Sprint(int64(f)), nil
		}
		return fmt.Sprint(f), nil

	default:
		return fmt.Sprint(i), nil
	}
}

func (j *TopupJob) UnmarshalJSON(b []byte) error {
	castFns := gotils.CastMap{
		"Number": GimmeString,
		"ID":     GimmeString,
	}
	obj, err := gotils.MarshalWithCasts(b, TopupJob{}, castFns)
	if err != nil {
		return err
	}

	*j = obj.(TopupJob)
	return nil
}

type TopupJob struct {
	Number           string  `csv:"number" json:"number" validate:"required"`
	Amount           float64 `csv:"amount" json:"amount" validate:"required"`
	Country          string  `csv:"country" json:"country" validate:"required"`
	Tolerance        float64 `csv:"tolerance,omitempty" json:"tolerance,omitempty"`
	Operator         string  `csv:"operator,omitempty" json:"operator,omitempty"`
	ID               string  `csv:"id,omitempty" json:"id,omitempty"`
	CustomIdentifier string  `csv:"custom_identifier,omitempty" json:"custom_identifier,omitempty"`
}

type TopupWorkerResponse struct {
	*TopupResponse
	ErrorMessage string `csv:"errrorMessage" json:"errorMessage,omitempty"`
	ErrorCode    string `csv:"errorCode" json:"errorCode,omitempty"`
}

func (r *TopupWorkerResponse) SetError(err error) *TopupWorkerResponse {
	r.ErrorMessage = err.Error()

	if e, ok := err.(APIError); ok {
		r.ErrorCode = e.ErrorCode
	} else if e, ok := err.(ReloadlyError); ok {
		r.ErrorCode = e.ErrorCode
	}

	return r
}

func workErrorResponse(err error, d *TopupJob) *TopupWorkerResponse {
	tr := &TopupResponse{}
	tr.OperatorName = d.Operator
	tr.RecipientPhone = d.Number
	tr.CountryCode = d.Country
	tr.RequestedAmount = d.Amount

	r := &TopupWorkerResponse{tr, "", ""}
	r.SetError(err)
	return r
}

func (t *TopupWorker) DoJob(d *TopupJob) (*TopupResponse, error) {
	svc := Service(*t)

	s := svc.Topups()

	if d.Operator != "" {
		s = s.FindOperator(d.Country, d.Operator).SuggestedAmount(d.Tolerance).AutoFallback()
	} else {
		s = s.AutoDetect(d.Country).SuggestedAmount(d.Tolerance)
	}

	if d.CustomIdentifier != "" {
		s = s.CustomIdentifier(d.CustomIdentifier)
	}

	return s.Topup(d.Number, d.Amount)
}

func (t *TopupWorker) Do(d *TopupJob) *TopupWorkerResponse {
	res, err := t.DoJob(d)

	if err != nil {
		return workErrorResponse(err, d)
	}

	return &TopupWorkerResponse{res, "", ""}
}

func (t *TopupWorker) Work(i interface{}) interface{} {
	d := i.(*TopupJob)
	return t.Do(d)
}
