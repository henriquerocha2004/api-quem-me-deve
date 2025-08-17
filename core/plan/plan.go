package plan

import (
	"encoding/json"
	"fmt"

	"github.com/oklog/ulid/v2"
)

const (
	Free planType = iota
	Pro
)

type FeatureDesc string

const (
	FeatureDebt   FeatureDesc = "debt"
	FeatureClient FeatureDesc = "client"
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
	Free: "Free",
	Pro:  "Pro",
}

var PlanTypeValue = map[string]planType{
	"Free": Free,
	"Pro":  Pro,
}

type Plan struct {
	Id       ulid.ULID
	Type     planType
	Price    float64
	Features []Feature
}

func (p *Plan) AddFeature(featDesc FeatureDesc, limit uint) {
	feature := Feature{
		Name:          featDesc,
		ResourceLimit: limit,
	}
	p.Features = append(p.Features, feature)
}

type Feature struct {
	Name          FeatureDesc
	ResourceLimit uint
}
