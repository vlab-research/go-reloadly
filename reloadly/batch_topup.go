package reloadly

type TopupDetails struct {
	Number string `csv:"number" json:"number"`
	Amount float64 `csv:"amount" json:"amount"`
	Country string `csv:"country" json:"country"`
	Tolerance float64 `csv:"tolerance,omitempty" json:"tolerance,omitempty"`
	Operator string `csv:"operator,omitempty" json:"operator,omitempty"`
}

type BatchTopupResponse struct {
	*TopupResponse
	ErrorMessage string `csv:"errrorMessage" json:"errorMessage,omitempty"`
	ErrorCode string `csv:"errorCode" json:"errorCode,omitempty"`
}

func (r *BatchTopupResponse) SetError(err error) *BatchTopupResponse {
	r.ErrorMessage = err.Error()

	if e, ok := err.(APIError); ok {
		r.ErrorCode = e.ErrorCode
	}
	if e, ok := err.(ReloadlyError); ok {
		r.ErrorCode = e.ErrorCode
	}

	return r
}

type TopupWorker Service

func workErrorResponse(err error, d *TopupDetails) *BatchTopupResponse {
	tr := &TopupResponse{}
	tr.OperatorName = d.Operator
	tr.RecipientPhone =  d.Number
	tr.CountryCode = d.Country
	tr.RequestedAmount = d.Amount

    r := &BatchTopupResponse{tr, "", ""}
	r.SetError(err)
	return r
}

func (t *TopupWorker) Work(i interface{}) interface{} {
	svc := Service(*t)

	d := i.(*TopupDetails)

	var err error
	var res *TopupResponse

	if d.Operator != "" {
		res, err = svc.Topups().FindOperator(d.Country, d.Operator).SuggestedAmount(d.Tolerance).AutoFallback().Topup(d.Number, d.Amount)
	} else {
		res, err = svc.Topups().AutoDetect(d.Country).SuggestedAmount(d.Tolerance).Topup(d.Number, d.Amount)
	}

	if err != nil {
		return workErrorResponse(err, d)
	}

	return &BatchTopupResponse{res, "", ""}
}
