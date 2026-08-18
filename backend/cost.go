package backend

import "fmt"

// Input holds a consignment's landed-cost inputs.
type Input struct {
	GoodsValue    float64 `json:"goodsValue"`    // pre-GST invoice value of the goods
	GSTPct        float64 `json:"gstPct"`        // GST rate on the goods
	Freight       float64 `json:"freight"`       // transport in
	LoadingCoolie float64 `json:"loadingCoolie"` // loading/unloading labour
	LocalTransport float64 `json:"localTransport"` // last-mile cartage
	WastagePct    float64 `json:"wastagePct"`    // damage/shrink in transit
	Units         float64 `json:"units"`         // usable units before wastage
	Regular       bool    `json:"regular"`       // true = GST-registered (ITC); false = composition/unregistered
}

// Result is the landed-cost breakdown.
type Result struct {
	GSTInCost        float64 `json:"gstInCost"`        // GST added to cost (0 for regular/ITC)
	LandedTotal      float64 `json:"landedTotal"`      // total landed cost of the consignment
	UsableUnits      float64 `json:"usableUnits"`      // units net of wastage
	LandedCostPerUnit float64 `json:"landedCostPerUnit"`
}

// Validate reports whether the Input is well formed.
func (in Input) Validate() error {
	if in.GoodsValue < 0 || in.Freight < 0 || in.LoadingCoolie < 0 || in.LocalTransport < 0 {
		return fmt.Errorf("costs cannot be negative")
	}
	if in.GSTPct < 0 || in.GSTPct > 100 {
		return fmt.Errorf("GST %% must be between 0 and 100")
	}
	if in.WastagePct < 0 || in.WastagePct >= 100 {
		return fmt.Errorf("wastage %% must be between 0 and 100")
	}
	if in.Units <= 0 {
		return fmt.Errorf("units must be positive")
	}
	return nil
}

// LandedCostPerUnit computes the true delivered cost per usable unit. The GST
// treatment is decisive: a regular (ITC) shop recovers GST so it is NOT a cost;
// a composition/unregistered shop cannot, so GST is ADDED to the cost base.
func LandedCostPerUnit(in Input) Result {
	gstInCost := 0.0
	if !in.Regular {
		gstInCost = in.GoodsValue * in.GSTPct / 100
	}
	landed := in.GoodsValue + gstInCost + in.Freight + in.LoadingCoolie + in.LocalTransport
	usable := in.Units * (1 - in.WastagePct/100)
	perUnit := 0.0
	if usable > 0 {
		perUnit = landed / usable
	}
	return Result{
		GSTInCost:         gstInCost,
		LandedTotal:       landed,
		UsableUnits:       usable,
		LandedCostPerUnit: perUnit,
	}
}
