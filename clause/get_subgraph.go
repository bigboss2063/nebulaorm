package clause

import (
	"strconv"
)

type GetSubgraph struct {
	WithProp  bool
	StepCount int
}

const GetSubgraphName = "GET_SUBGRAPH"

func (gs GetSubgraph) Name() string {
	return GetSubgraphName
}

func (gs GetSubgraph) MergeIn(clause *Clause) {
	clause.Expression = gs
}

func (gs GetSubgraph) Build(nGQL Builder) error {
	nGQL.WriteString("GET SUBGRAPH ")
	if gs.WithProp {
		nGQL.WriteString("WITH PROP ")
	}
	if gs.StepCount <= 0 {
		gs.StepCount = 1
	}
	nGQL.WriteString(strconv.Itoa(gs.StepCount))
	nGQL.WriteString(" STEPS")
	return nil
}
