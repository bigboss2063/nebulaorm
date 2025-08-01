package statement

import (
	"fmt"
	"github.com/haysons/norm/clause"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateVertexTags(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().CreateVertexTags(&vm1{}, true)
			},
			want: `CREATE TAG IF NOT EXISTS player(name string, age int);`,
		},
		{
			stmt: func() *Statement {
				return New().CreateVertexTags(&vm2{}, true)
			},
			want: `CREATE TAG IF NOT EXISTS no_property();`,
		},
		{
			stmt: func() *Statement {
				return New().CreateVertexTags(&vm3{}, true)
			},
			want: `CREATE TAG IF NOT EXISTS player_with_default(name string DEFAULT "", age int DEFAULT 20);`,
		},
		{
			stmt: func() *Statement {
				return New().CreateVertexTags(vm4{}, true)
			},
			want: `CREATE TAG IF NOT EXISTS woman(name string, age int, married bool, salary double, create_time timestamp) TTL_DURATION = 100, TTL_COL = "create_time";`,
		},
		{
			stmt: func() *Statement {
				return New().CreateVertexTags(vm5{})
			},
			want: `CREATE TAG woman(name fixed_string DEFAULT "hayson", age int32 DEFAULT 20, create_time datetime DEFAULT datetime(1625469277));`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, ngql)
			}
		})
	}
}

func TestDropTag(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().DropVertexTag("test")
			},
			want: `DROP TAG test;`,
		},
		{
			stmt: func() *Statement {
				return New().DropVertexTag("test", true)
			},
			want: `DROP TAG IF EXISTS test;`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, ngql)
			}
		})
	}
}

func TestAlterTag(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().AlterVertexTag(&vm1{}, clause.AlterOperate{
					AddProps: []string{"name", "age"},
				})
			},
			want: `ALTER TAG player ADD (name string, age int);`,
		},
		{
			stmt: func() *Statement {
				return New().AlterVertexTag(&vm2{}, clause.AlterOperate{
					DropProps: []string{"name", "age"},
				})
			},
			want: `ALTER TAG no_property DROP (name, age);`,
		},
		{
			stmt: func() *Statement {
				return New().AlterVertexTag(&vm3{}, clause.AlterOperate{
					ChangeProps: []string{"name", "age"},
				})
			},
			want: `ALTER TAG player_with_default CHANGE (name string DEFAULT "", age int DEFAULT 20);`,
		},
		{
			stmt: func() *Statement {
				return New().AlterVertexTag(&vm4{}, clause.AlterOperate{
					AddProps:    []string{"name", "age"},
					DropProps:   []string{"salary"},
					ChangeProps: []string{"create_time"},
					UpdateTTL:   true,
				})
			},
			want: `ALTER TAG woman ADD (name string, age int), DROP (salary), CHANGE (create_time timestamp) TTL_DURATION = 100, TTL_COL = "create_time";`,
		},
		{
			stmt: func() *Statement {
				return New().AlterVertexTag(v1{}, clause.AlterOperate{
					AddProps: []string{"p1"},
				})
			},
			want: `ALTER TAG t3 ADD (p1 int);`,
		},
		{
			stmt: func() *Statement {
				return New().AlterVertexTag(v1{}, clause.AlterOperate{
					AddProps: []string{"p2"},
				}, clause.WithTagName("t4"))
			},
			want: `ALTER TAG t4 ADD (p2 string);`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, ngql)
			}
		})
	}
}

func TestCreateEdge(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().CreateEdge(&em1{}, true)
			},
			want: `CREATE EDGE IF NOT EXISTS follow(degree int);`,
		},
		{
			stmt: func() *Statement {
				return New().CreateEdge(em2{}, true)
			},
			want: `CREATE EDGE IF NOT EXISTS no_property();`,
		},
		{
			stmt: func() *Statement {
				return New().CreateEdge(&em3{}, true)
			},
			want: `CREATE EDGE IF NOT EXISTS follow_with_default(degree int DEFAULT 20);`,
		},
		{
			stmt: func() *Statement {
				return New().CreateEdge(em4{}, true)
			},
			want: `CREATE EDGE IF NOT EXISTS e1(p1 string, p2 int, p3 timestamp) TTL_DURATION = 100, TTL_COL = "p2";`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, ngql)
			}
		})
	}
}

func TestDropEdge(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().DropEdge("e1")
			},
			want: `DROP EDGE e1;`,
		},
		{
			stmt: func() *Statement {
				return New().DropEdge("e1", true)
			},
			want: `DROP EDGE IF EXISTS e1;`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, ngql)
			}
		})
	}
}

func TestAlterEdge(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().AlterEdge(&e2{}, clause.AlterOperate{
					AddProps: []string{"name", "age"},
				})
			},
			want: `ALTER EDGE e2 ADD (name string, age int);`,
		},
		{
			stmt: func() *Statement {
				return New().AlterEdge(&em2{}, clause.AlterOperate{
					DropProps: []string{"name", "age"},
				})
			},
			want: `ALTER EDGE no_property DROP (name, age);`,
		},
		{
			stmt: func() *Statement {
				return New().AlterEdge(&em3{}, clause.AlterOperate{
					ChangeProps: []string{"degree"},
				})
			},
			want: `ALTER EDGE follow_with_default CHANGE (degree int DEFAULT 20);`,
		},
		{
			stmt: func() *Statement {
				return New().AlterEdge(em4{}, clause.AlterOperate{
					AddProps:    []string{"p1"},
					ChangeProps: []string{"p3"},
					UpdateTTL:   true,
				})
			},
			want: `ALTER EDGE e1 ADD (p1 string), CHANGE (p3 timestamp) TTL_DURATION = 100, TTL_COL = "p2";`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, ngql)
			}
		})
	}
}

type vm1 struct {
	VID  string `norm:"vertex_id"`
	Name string
	Age  int
}

func (t vm1) VertexID() string {
	return t.VID
}

func (t vm1) VertexTagName() string {
	return "player"
}

type vm2 struct {
	VID string `norm:"vertex_id"`
}

func (t vm2) VertexID() string {
	return t.VID
}

func (t vm2) VertexTagName() string {
	return "no_property"
}

type vm3 struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name;type:string;default:''"`
	Age  int    `norm:"prop:age;type:int;default:20"`
}

func (t vm3) VertexID() string {
	return t.VID
}

func (t vm3) VertexTagName() string {
	return "player_with_default"
}

type vm4 struct {
	VID        string `norm:"vertex_id"`
	Name       string
	Age        int
	Married    bool
	Salary     float64
	CreateTime time.Time `norm:"type:timestamp;ttl:100"`
}

func (t *vm4) VertexID() string {
	return t.VID
}

func (t *vm4) VertexTagName() string {
	return "woman"
}

type vm5 struct {
	VID        string    `norm:"vertex_id"`
	Name       string    `norm:"prop:name;type:fixed_string;default:hayson"`
	Age        int       `norm:"prop:age;type:int32;default:20"`
	CreateTime time.Time `norm:"default:datetime(1625469277)"`
}

func (t *vm5) VertexID() string {
	return t.VID
}

func (t *vm5) VertexTagName() string {
	return "woman"
}

type em1 struct {
	SrcID  string `norm:"edge_src_id"`
	DstID  string `norm:"edge_dst_id"`
	Degree int
}

func (e em1) EdgeTypeName() string {
	return "follow"
}

type em2 struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
}

func (e em2) EdgeTypeName() string {
	return "no_property"
}

type em3 struct {
	SrcID  string `norm:"edge_src_id"`
	DstID  string `norm:"edge_dst_id"`
	Degree int    `norm:"default:20"`
}

func (e em3) EdgeTypeName() string {
	return "follow_with_default"
}

type em4 struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	P1    string
	P2    int    `norm:"ttl:100"`
	P3    string `norm:"prop:p3;type:timestamp"`
}

func (e em4) EdgeTypeName() string {
	return "e1"
}
