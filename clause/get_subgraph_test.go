package clause_test

import (
	"fmt"
	"testing"

	"github.com/haysons/norm/clause"
)

func TestGetSubgraph(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.GetSubgraph{}},
			gqlWant: "GET SUBGRAPH 1 STEPS",
		},
		{
			clauses: []clause.Interface{clause.GetSubgraph{WithProp: true, StepCount: -1}},
			gqlWant: "GET SUBGRAPH WITH PROP 1 STEPS",
		},
		{
			clauses: []clause.Interface{clause.GetSubgraph{WithProp: true, StepCount: -1}},
			gqlWant: "GET SUBGRAPH WITH PROP 1 STEPS",
		},
		{
			clauses: []clause.Interface{clause.GetSubgraph{WithProp: false, StepCount: 100}},
			gqlWant: "GET SUBGRAPH 100 STEPS",
		},
		{
			clauses: []clause.Interface{clause.GetSubgraph{WithProp: true, StepCount: 2}},
			gqlWant: "GET SUBGRAPH WITH PROP 2 STEPS",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
