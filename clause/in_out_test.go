package clause_test

import (
	"fmt"
	"testing"

	"github.com/haysons/norm/clause"
)

func TestInOut(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.In{InOutBase: clause.InOutBase{EdgeTypeList: []string{"follow"}}}},
			gqlWant: "IN follow",
		},
		{
			clauses: []clause.Interface{clause.In{InOutBase: clause.InOutBase{EdgeTypeList: []string{"follow", "server"}}}},
			gqlWant: "IN follow, server",
		},
		{
			clauses: []clause.Interface{clause.In{InOutBase: clause.InOutBase{EdgeTypeList: []string{"follow"}}}, clause.In{InOutBase: clause.InOutBase{EdgeTypeList: []string{"server"}}}},
			gqlWant: "IN follow, server",
		},
		{
			clauses: []clause.Interface{clause.In{InOutBase: clause.InOutBase{EdgeTypeList: []string{""}}}},
			errWant: clause.ErrInvalidClauseParams,
		},
		{
			clauses: []clause.Interface{clause.Out{InOutBase: clause.InOutBase{EdgeTypeList: []string{"follow"}}}},
			gqlWant: "OUT follow",
		},
		{
			clauses: []clause.Interface{clause.Out{InOutBase: clause.InOutBase{EdgeTypeList: []string{"follow", "server"}}}},
			gqlWant: "OUT follow, server",
		},
		{
			clauses: []clause.Interface{clause.Out{InOutBase: clause.InOutBase{EdgeTypeList: []string{"follow"}}}, clause.Out{InOutBase: clause.InOutBase{EdgeTypeList: []string{"server"}}}},
			gqlWant: "OUT follow, server",
		},
		{
			clauses: []clause.Interface{clause.Out{}},
			errWant: clause.ErrInvalidClauseParams,
		},
		{
			clauses: []clause.Interface{clause.Both{InOutBase: clause.InOutBase{EdgeTypeList: []string{"follow"}}}},
			gqlWant: "BOTH follow",
		},
		{
			clauses: []clause.Interface{clause.Both{InOutBase: clause.InOutBase{EdgeTypeList: []string{"follow", "server"}}}},
			gqlWant: "BOTH follow, server",
		},
		{
			clauses: []clause.Interface{clause.Both{InOutBase: clause.InOutBase{EdgeTypeList: []string{"follow"}}}, clause.Both{InOutBase: clause.InOutBase{EdgeTypeList: []string{"server"}}}},
			gqlWant: "BOTH follow, server",
		},
		{
			clauses: []clause.Interface{clause.Both{}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
