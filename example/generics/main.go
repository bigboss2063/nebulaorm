package main

import (
	"context"
	"log"
	"time"

	"github.com/haysons/norm"
	"github.com/haysons/norm/clause"
)

type Player struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name"`
	Age  int    `norm:"prop:age"`
}

func (p Player) VertexID() string {
	return p.VID
}

func (p Player) VertexTagName() string {
	return "player"
}

type Team struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name"`
}

func (t Team) VertexID() string {
	return t.VID
}

func (t Team) VertexTagName() string {
	return "team"
}

type Serve struct {
	SrcID     string `norm:"edge_src_id"`
	DstID     string `norm:"edge_dst_id"`
	Rank      int    `norm:"edge_rank"`
	StartYear int64  `norm:"prop:start_year"`
	EndYear   int64  `norm:"prop:end_year"`
}

func (s Serve) EdgeTypeName() string {
	return "serve"
}

var db *norm.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conf := &norm.Config{
		Username:    "root",
		Password:    "nebula",
		SpaceName:   "test",
		Addresses:   []string{"127.0.0.1:9669"},
		ConnTimeout: 10 * time.Second,
	}
	var err error
	db, err = norm.Open(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Migrator().AutoMigrateVertexes(Player{}, Team{}); err != nil {
		log.Fatal(err)
	}
	if err = db.Migrator().AutoMigrateEdges(Serve{}); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	Insert(ctx)

	Query(ctx)

	Update(ctx)

	Delete(ctx)
}

func Insert(ctx context.Context) {
	player := &Player{
		VID:  "player1001",
		Name: "Kobe Bryant",
		Age:  33,
	}
	if err := norm.G[Player](db).InsertVertex(player).Exec(ctx); err != nil {
		log.Fatalf("insert player failed: %v", err)
	}
	team := &Team{
		VID:  "team1001",
		Name: "Lakers",
	}
	if err := norm.G[Team](db).InsertVertex(team).Exec(ctx); err != nil {
		log.Fatalf("insert team failed: %v", err)
	}
	serve := &Serve{
		SrcID:     "player1001",
		DstID:     "team1001",
		StartYear: time.Date(1996, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
		EndYear:   time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
	}
	if err := norm.G[Serve](db).InsertEdge(serve).Exec(ctx); err != nil {
		log.Fatalf("insert serve failed: %v", err)
	}
}

func Query(ctx context.Context) {
	player, err := norm.G[Player](db).
		Fetch("player", "player1001").
		Yield("vertex as v").
		TakeCol(ctx, "v")
	if err != nil {
		log.Fatalf("fetch player failed: %v", err)
	}
	log.Printf("player: %+v", player)

	team, err := norm.G[Team](db).
		Go().
		From("player1001").
		Over("serve").
		Yield("$$ as t").
		TakeCol(ctx, "t")
	if err != nil {
		log.Fatalf("fetch team failed: %v", err)
	}
	log.Printf("team: %+v", team)

	serve, err := norm.G[Serve](db).
		Go().
		From("player1001").
		Over("serve").
		Yield("edge as e").
		FindCol(ctx, "e")
	if err != nil {
		log.Fatalf("fetch serve failed: %v", err)
	}
	for _, s := range serve {
		log.Printf("serve: %+v", s)
	}

	type edgeCnt struct {
		Edge string `norm:"col:e"`
		Cnt  int    `norm:"col:cnt"`
	}
	edgesCnt, err := norm.G[edgeCnt](db).
		Go().
		From("player1001").
		Over("*").
		Yield("type(edge) as t").
		GroupBy("$-.t").
		Yield("$-.t as e, count(*) as cnt").
		Find(ctx)
	if err != nil {
		log.Fatalf("get edge cnt failed: %v", err)
	}
	for _, c := range edgesCnt {
		log.Printf("edge cnt: %+v\n", c)
	}

	type edgeVertex struct {
		ID string `norm:"col:id"`
		T  *Team  `norm:"col:t"`
	}
	edgeVertexes, err := norm.G[edgeVertex](db).
		Go().
		From("player1001").
		Over("serve").
		Yield("id($^) as id, $$ as t").
		Find(ctx)
	if err != nil {
		log.Fatalf("get edge vertex failed: %v", err)
	}
	for _, v := range edgeVertexes {
		log.Printf("id: %v, t: %+v", v.ID, v.T)
	}
}

func Update(ctx context.Context) {
	if err := norm.G[Player](db).UpdateVertex("player1001", &Player{Age: 23}).Exec(ctx); err != nil {
		log.Fatalf("update player failed: %v", err)
	}
	prop, err := norm.G[map[string]any](db).
		Fetch("player", "player1001").
		Yield("properties(vertex) as p").
		FindCol(ctx, "p")
	if err != nil {
		log.Fatalf("fetch player failed: %v", err)
	}
	log.Printf("vertex prop after update: %+v", prop)

	if err = norm.G[Serve](db).UpdateEdge(Serve{SrcID: "player1001", DstID: "team1001"}, &Serve{StartYear: 160123456}).Exec(ctx); err != nil {
		log.Fatalf("update edge serve failed: %v", err)
	}
	prop, err = norm.G[map[string]any](db).
		Fetch("serve", clause.Expr{Str: `"player1001"->"team1001"`}).
		Yield("properties(edge) as p").
		FindCol(ctx, "p")
	if err != nil {
		log.Fatalf("fetch serve failed: %v", err)
	}
	log.Printf("edge prop after update: %+v", prop)
}

func Delete(ctx context.Context) {
	if err := norm.G[Player](db).DeleteVertex("player1001").Exec(ctx); err != nil {
		log.Fatalf("delete player failed: %v", err)
	}
	player, err := norm.G[Player](db).
		Fetch("player", "player1001").
		Yield("vertex as v").
		TakeCol(ctx, "v")
	if err != nil {
		log.Printf("after delete, fetch player failed: %v", err)
	} else {
		log.Printf("after delete, fetch player: %+v", player)
	}

	if err = norm.G[Serve](db).DeleteEdge("serve", Serve{SrcID: "player1001", DstID: "team1001"}).Exec(ctx); err != nil {
		log.Fatalf("delete edge serve failed: %v", err)
	}
	serve, err := norm.G[Serve](db).
		Fetch("serve", clause.Expr{Str: `"player1001"->"team1001"`}).
		Yield("edge as e").
		TakeCol(ctx, "e")
	if err != nil {
		log.Printf("after delete, fetch server failed: %v", err)
	} else {
		log.Printf("after delete, fetch server: %+v", serve)
	}
}
