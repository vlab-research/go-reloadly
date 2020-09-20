package reloadly

type TopupWorker Service

type TopupJob struct {
	Number    string  `csv:"number" json:"number" validate:"required"`
	Amount    float64 `csv:"amount" json:"amount" validate:"required"`
	Country   string  `csv:"country" json:"country" validate:"required"`
	Tolerance float64 `csv:"tolerance,omitempty" json:"tolerance,omitempty"`
	Operator  string  `csv:"operator,omitempty" json:"operator,omitempty"`
	ID        string  `csv:"id,omitempty" json:"id,omitempty"`
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

	var err error
	var res *TopupResponse

	if d.Operator != "" {
		res, err = svc.Topups().FindOperator(d.Country, d.Operator).SuggestedAmount(d.Tolerance).AutoFallback().Topup(d.Number, d.Amount)
	} else {
		res, err = svc.Topups().AutoDetect(d.Country).SuggestedAmount(d.Tolerance).Topup(d.Number, d.Amount)
	}

	return res, err
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
