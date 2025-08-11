package plan

import (
	"encoding/json"
	"fmt"
)

const (
	Free planType = iota
	Starter
	Pro
)

type planType int

func (p planType) String() string {
	if int(p) >= 0 && int(p) < len(PlanTypeString) {
		return PlanTypeString[p]
	}

	return "unknown"
}

func (p *planType) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	if val, ok := PlanTypeValue[str]; ok {
		*p = val
		return nil
	}

	return fmt.Errorf("invalid plan type: %s", str)
}

var PlanTypeString = map[planType]string{
	Free:    "Free",
	Starter: "Starter",
	Pro:     "Pro",
}

var PlanTypeValue = map[string]planType{
	"Free":    Free,
	"Starter": Starter,
	"Pro":     Pro,
}

type Feature struct {
	Name          string
	ResourceLimit uint
}

type Plan struct {
	Type     planType
	Price    float64
	Features []Feature
	Status   string
}
