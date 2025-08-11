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
	Reversed
)

var statusString = map[status]string{
	Pending:  "pending",
	Paid:     "paid",
	Canceled: "canceled",
	Reversed: "reversed",
}
var StatusValue = map[string]status{
	"pending":  Pending,
	"paid":     Paid,
	"canceled": Canceled,
	"reversed": Reversed,
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

	if val, ok := StatusValue[str]; ok {
		*s = val
		return nil
	}

	return fmt.Errorf("invalid status: %s", str)
}
