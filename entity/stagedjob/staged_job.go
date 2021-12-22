package stagedjob

type JobName string

const (
	JobNameSendReceipt JobName = "send_ride_receipt"
)

func (j JobName) String() string {
	return string(j)
}

type JobArgReceipt struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	UserID   int64  `json:"user_id"`
}
