package norm

import (
	"context"

	"github.com/haysons/norm/clause"
	nebula "github.com/vesoft-inc/nebula-go/v3"
)

type op func(*DB) *DB

type g[T any] struct {
	*chainG[T]
	db  *DB
	ops []op
}

func (g *g[T]) apply(_ context.Context) *DB {
	db := g.db.session()

	for _, op := range g.ops {
		db = op(db)
	}
	return db
}

type chainG[T any] struct {
	execG[T]
}

func (c chainG[T]) with(v op) chainG[T] {
	return chainG[T]{
		execG: execG[T]{g: &g[T]{
			db:  c.g.db,
			ops: append(append([]op(nil), c.g.ops...), v),
		}},
	}
}

func (c chainG[T]) Raw(raw string) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Raw(raw)
	})
}

func (c chainG[T]) Go(step ...int) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Go(step...)
	})
}

func (c chainG[T]) From(vid any) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.From(vid)
	})
}

func (c chainG[T]) Over(edgeType ...string) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Over(edgeType...)
	})
}

func (c chainG[T]) Where(query string, args ...any) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Where(query, args...)
	})
}

func (c chainG[T]) Or(query string, args ...any) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Or(query, args...)
	})
}

func (c chainG[T]) Not(query string, args ...any) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Not(query, args...)
	})
}

func (c chainG[T]) Xor(query string, args ...any) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Xor(query, args...)
	})
}

func (c chainG[T]) Sample(sampleList ...int) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Sample(sampleList...)
	})
}

func (c chainG[T]) Fetch(name string, vid any) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Fetch(name, vid)
	})
}

func (c chainG[T]) FetchMulti(names []string, vid any) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.FetchMulti(names, vid)
	})
}

func (c chainG[T]) Lookup(name string) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Lookup(name)
	})
}

func (c chainG[T]) GroupBy(expr string) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.GroupBy(expr)
	})
}

func (c chainG[T]) Yield(expr string, distinct ...bool) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Yield(expr, distinct...)
	})
}

func (c chainG[T]) OrderBy(expr string) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.OrderBy(expr)
	})
}

func (c chainG[T]) Limit(limit int) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Limit(limit)
	})
}

func (c chainG[T]) InsertVertex(vertexes any, ifNotExists ...bool) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.InsertVertex(vertexes, ifNotExists...)
	})
}

func (c chainG[T]) UpdateVertex(vid any, propsUpdate any, opts ...clause.Option) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.UpdateVertex(vid, propsUpdate, opts...)
	})
}

func (c chainG[T]) UpsertVertex(vid any, propsUpdate any, opts ...clause.Option) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.UpsertVertex(vid, propsUpdate, opts...)
	})
}

func (c chainG[T]) DeleteVertex(vid any, withEdge ...bool) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.DeleteVertex(vid, withEdge...)
	})
}

func (c chainG[T]) InsertEdge(edges any, ifNotExists ...bool) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.InsertEdge(edges, ifNotExists...)
	})
}

func (c chainG[T]) UpdateEdge(edge any, propsUpdate any, opts ...clause.Option) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.UpdateEdge(edge, propsUpdate, opts...)
	})
}

func (c chainG[T]) UpsertEdge(edge any, propsUpdate any, opts ...clause.Option) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.UpsertEdge(edge, propsUpdate, opts...)
	})
}

func (c chainG[T]) DeleteEdge(edgeTypeName string, edge any) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.DeleteEdge(edgeTypeName, edge)
	})
}

func (c chainG[T]) When(query string, args ...any) chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.When(query, args...)
	})
}

func (c chainG[T]) Pipe() chainG[T] {
	return c.with(func(db *DB) *DB {
		return db.Pipe()
	})
}

type execG[T any] struct {
	g *g[T]
}

func (g execG[T]) NGQL(ctx context.Context) (string, error) {
	return g.g.apply(ctx).NGQL()
}

func (g execG[T]) RawResult(ctx context.Context) (*nebula.ResultSet, error) {
	return g.g.apply(ctx).RawResult()
}

func (g execG[T]) Exec(ctx context.Context) error {
	return g.g.apply(ctx).Exec()
}

func (g execG[T]) Find(ctx context.Context) ([]T, error) {
	var r []T
	err := g.g.apply(ctx).Find(&r)
	return r, err
}

func (g execG[T]) FindCol(ctx context.Context, col string) ([]T, error) {
	var r []T
	err := g.g.apply(ctx).FindCol(col, &r)
	return r, err
}

func (g execG[T]) Take(ctx context.Context) (T, error) {
	var r T
	err := g.g.apply(ctx).Take(&r)
	return r, err
}

func (g execG[T]) TakeCol(ctx context.Context, col string) (T, error) {
	var r T
	err := g.g.apply(ctx).TakeCol(col, &r)
	return r, err
}
