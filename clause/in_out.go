package clause

import (
	"fmt"
	"strings"
)

type In struct {
	InOutBase
}

const InName = "IN"

func (in In) Name() string {
	return InName
}

func (in In) MergeIn(clause *Clause) {
	exist, ok := clause.Expression.(In)
	if !ok {
		clause.Expression = in
		return
	}
	// edge type merge
	exist.EdgeTypeList = append(exist.EdgeTypeList, in.EdgeTypeList...)
	clause.Expression = exist
}

func (in In) Build(nGQL Builder) error {
	return in.build(nGQL, "IN")
}

type Out struct {
	InOutBase
}

const OutName = "Out"

func (out Out) Name() string {
	return OutName
}

func (out Out) MergeIn(clause *Clause) {
	exist, ok := clause.Expression.(Out)
	if !ok {
		clause.Expression = out
		return
	}
	// edge type merge
	exist.EdgeTypeList = append(exist.EdgeTypeList, out.EdgeTypeList...)
	clause.Expression = exist
}

func (out Out) Build(nGQL Builder) error {
	return out.build(nGQL, "OUT")
}

type Both struct {
	InOutBase
}

const BothName = "Both"

func (both Both) Name() string {
	return BothName
}

func (both Both) MergeIn(clause *Clause) {
	exist, ok := clause.Expression.(Both)
	if !ok {
		clause.Expression = both
		return
	}
	// edge type merge
	exist.EdgeTypeList = append(exist.EdgeTypeList, both.EdgeTypeList...)
	clause.Expression = exist
}

func (both Both) Build(nGQL Builder) error {
	return both.build(nGQL, "BOTH")
}

type InOutBase struct {
	EdgeTypeList []string
}

func (in InOutBase) build(nGQL Builder, clause string) error {
	edgeTypeList := make([]string, 0, len(in.EdgeTypeList))
	for _, edgeType := range in.EdgeTypeList {
		if edgeType != "" {
			edgeTypeList = append(edgeTypeList, edgeType)
		}
	}
	if len(edgeTypeList) == 0 {
		return fmt.Errorf("norm: %w, edge type list is empty in 'in', 'out', 'both' clause", ErrInvalidClauseParams)
	}
	nGQL.WriteString(clause)
	nGQL.WriteByte(' ')
	nGQL.WriteString(strings.Join(edgeTypeList, ", "))
	return nil
}
