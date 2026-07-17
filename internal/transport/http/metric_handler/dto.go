package metric_handler

import domainmetrics "laboratory-internet-ai-test/internal/domain/metrics"

type ErrorDTO struct {
	Error interface{} `json:"error"`
}

func newErrorDTO(err error) ErrorDTO {
	return ErrorDTO{Error: err.Error()}
}

//-----------------------------------

type metricDTO struct {
	DateKey           string `json:"date"`
	CriticalStatus    string `json:"critical"`
	OKStatus          string `json:"ok"`
	BadReqStatus      string `json:"bad_request"`
	ManyRequestStatus string `json:"many_request"`
	AnotherStatus     string `json:"another"`
}

func newMetricDTO(metric domainmetrics.Metric) metricDTO {
	return metricDTO{
		DateKey:           metric.DateKey,
		CriticalStatus:    metric.CriticalStatus,
		OKStatus:          metric.OKStatus,
		BadReqStatus:      metric.BadReqStatus,
		ManyRequestStatus: metric.ManyRequestStatus,
		AnotherStatus:     metric.AnotherStatus,
	}
}
