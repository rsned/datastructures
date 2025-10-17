package tree

import (
	"math"
	"strings"
)

// BalanceQuality is a enum for the balance quality of a tree.
// It is an arbitrary set of values ranging from Well Balanced to
// Severely Right/Left Heavy, to Degenerate.
// The change over points are chosen by me to what feels reasonable.
type BalanceQuality int

const (
	BalanceUnknown BalanceQuality = iota
	DegenerateLeft                // Essentially a linked list
	SeverelyLeftHeavy
	ModeratelyLeftHeavy
	SlightlyLeftHeavy
	WellBalanced
	SlightlyRightHeavy
	ModeratelyRightHeavy
	SeverelyRightHeavy
	DegenerateRight // Essentially a linked list
)

// balanceQualityForScore converts a balance score [-1.0, 1.0] to
// a BalanceQuality enum.
func balanceQualityForScore(score float64) BalanceQuality {
	switch {
	case math.Abs(score) <= 0.1:
		return WellBalanced
	case score <= -1.0:
		return DegenerateLeft
	case score <= -0.6:
		return SeverelyLeftHeavy
	case score <= -0.3:
		return ModeratelyLeftHeavy
	case score <= -0.1:
		return SlightlyLeftHeavy
	case score >= 1.0:
		return DegenerateRight
	case score >= 0.6:
		return SeverelyRightHeavy
	case score >= 0.3:
		return ModeratelyRightHeavy
	case score >= 0.1:
		return SlightlyRightHeavy
	}

	return BalanceUnknown
}

// balanceQualityString converts BalanceQuality enum to string
func (bq BalanceQuality) String() string {
	switch bq {
	case BalanceUnknown:
		return treeTypeUnknown // Replace this const with a better named one.
	case DegenerateLeft:
		return "Degenerate Left (Linear)"
	case SeverelyLeftHeavy:
		return "Severely Left Heavy"
	case ModeratelyLeftHeavy:
		return "Moderately Left Heavy"
	case SlightlyLeftHeavy:
		return "Slightly Left Heavy"
	case WellBalanced:
		return "Well Balanced"
	case SlightlyRightHeavy:
		return "Slightly Right Heavy"
	case ModeratelyRightHeavy:
		return "Moderately Right Heavy"
	case SeverelyRightHeavy:
		return "Severely Right Heavy"
	case DegenerateRight:
		return "Degenerate Right (Linear)"
	default:
		return treeTypeUnknown // Replace this const with a better named one.
	}
}

// BalanceQualityGraph returns a string representation of the balance quality
// as an short ASCII art horizontal graph with label.
func BalanceQualityGraph(score float64) string {
	const lineChartLabel = "-1.0         0         1.0"
	const lineChart = "|.........|.........|"
	const labelPadding = "   "
	sb := strings.Builder{}

	sb.WriteString(lineChartLabel + "\n")
	sb.WriteString(labelPadding + lineChart + "\n")
	// Convert the score [-1, 1] to a horizontal offset into the line chart.
	offset := int((score + 1) * 10)
	sb.WriteString(labelPadding)
	sb.WriteString(strings.Repeat(" ", offset))
	sb.WriteString("^\n")
	sb.WriteString(centerString(balanceQualityForScore(score).String(),
		" ", len(lineChart)+2*len(labelPadding)) + "\n")

	return sb.String()
}
