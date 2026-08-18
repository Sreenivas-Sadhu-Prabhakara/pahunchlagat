package backend

import (
	"math"
	"testing"
)

func almost(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

func TestLanded_CompositionAddsGST(t *testing.T) {
	// goods 10000 @18% GST, freight 500, coolie 200, transport 300, 5% wastage, 100 units.
	in := Input{GoodsValue: 10000, GSTPct: 18, Freight: 500, LoadingCoolie: 200,
		LocalTransport: 300, WastagePct: 5, Units: 100, Regular: false}
	r := LandedCostPerUnit(in)
	if !almost(r.GSTInCost, 1800) {
		t.Fatalf("gstInCost=%v want 1800", r.GSTInCost)
	}
	// landed = 10000+1800+500+200+300 = 12800 ; usable = 95 ; perUnit = 134.7368...
	if !almost(r.LandedTotal, 12800) || !almost(r.UsableUnits, 95) {
		t.Fatalf("landed=%v usable=%v", r.LandedTotal, r.UsableUnits)
	}
	if !almost(r.LandedCostPerUnit, 12800.0/95.0) {
		t.Fatalf("perUnit=%v", r.LandedCostPerUnit)
	}
}

func TestLanded_RegularExcludesGST(t *testing.T) {
	in := Input{GoodsValue: 10000, GSTPct: 18, Freight: 500, LoadingCoolie: 200,
		LocalTransport: 300, WastagePct: 5, Units: 100, Regular: true}
	r := LandedCostPerUnit(in)
	if r.GSTInCost != 0 {
		t.Fatalf("regular shop should recover GST, got %v in cost", r.GSTInCost)
	}
	if !almost(r.LandedTotal, 11000) {
		t.Fatalf("landed=%v want 11000", r.LandedTotal)
	}
}

func TestValidate(t *testing.T) {
	if err := (Input{GoodsValue: 100, Units: 10}).Validate(); err != nil {
		t.Fatalf("valid rejected: %v", err)
	}
	for i, bad := range []Input{
		{Units: 0}, {Units: 10, GSTPct: 120}, {Units: 10, WastagePct: 100}, {Units: 10, Freight: -1},
	} {
		if err := bad.Validate(); err == nil {
			t.Fatalf("bad %d accepted", i)
		}
	}
}
