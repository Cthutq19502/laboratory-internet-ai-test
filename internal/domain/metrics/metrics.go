package metrics

type Status string

const (
	CriticalStatus    Status = "Critical"
	OKStatus          Status = "OK"
	BadReqStatus      Status = "Bad request"
	ManyRequestStatus Status = "Many Request"
	AnotherStatus     Status = "Another"
)

type Metric struct {
	DateKey           string
	CriticalStatus    string
	OKStatus          string
	BadReqStatus      string
	ManyRequestStatus string
	AnotherStatus     string
}
