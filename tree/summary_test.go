package tree

import (
	"testing"
)

func TestBalanceQualityString(t *testing.T) {
	tests := []struct {
		quality BalanceQuality
		want    string
	}{
		{WellBalanced, "Well Balanced"},
		{SeverelyLeftHeavy, "Severely Left Heavy"},
		{SeverelyRightHeavy, "Severely Right Heavy"},
		{Degenerate, "Degenerate (Linear)"},
		{BalanceUnknown, "Unknown"},
	}

	for _, test := range tests {
		got := balanceQualityString(test.quality)
		if got != test.want {
			t.Errorf("balanceQualityString(%v) = %s, want %s", test.quality, got, test.want)
		}
	}
}
