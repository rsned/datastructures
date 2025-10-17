package tree

import (
	"fmt"
	"testing"
)

// testBalanceQualityGraph ranges over the whole set of balance scores
// and outputs the graph for viewing visually how the result looks.
func testBalanceQualityGraph(t *testing.T) {
	t.Helper()
	for score := -1.0; score <= 1.01; score += 0.05 {
		t.Run(fmt.Sprintf("score_%f", score), func(t *testing.T) {
			gotGraph := BalanceQualityGraph(score)
			t.Errorf("Graph:\n%s", gotGraph)
		})
	}
}
