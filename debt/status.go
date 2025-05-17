package debt

import (
	"encoding/json"
	"fmt"
)

type status int

const (
	Pending status = iota
	Paid
	Canceled
)

var statusString = map[status]string{
	Pending:  "pending",
	Paid:     "paid",
	Canceled: "canceled",
}
var statusValue = map[string]status{
	"pending":  Pending,
	"paid":     Paid,
	"canceled": Canceled,
}

func (s status) String() string {
	if int(s) >= 0 && int(s) < len(statusString) {
		return statusString[s]
	}

	return "unknown"
}

func (s *status) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	if val, ok := statusValue[str]; ok {
		*s = val
		return nil
	}

	return fmt.Errorf("invalid status: %s", str)
}
