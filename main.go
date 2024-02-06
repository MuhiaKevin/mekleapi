package main

import (
	"context"
	"fmt"
	"path"

	"merkleapi/config"
	myhttp "merkleapi/http"
	"merkleapi/provider"
	"merkleapi/provider/file"
	"merkleapi/provider/pg"
	"merkleapi/state"
	"merkleapi/types"
	"merkleapi/updates"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xssnick/tonutils-go/address"
)

func main() {
	config.LoadConfig()

	pool, err := pgxpool.New(context.Background(), config.Config.Database)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	var sp provider.StateProvider = &file.StateProvider{
		Path: path.Join(config.Config.DataDir, "state.json"),
	}
	var ip provider.ItemProvider = pg.NewItemProvider(pool)
	var np provider.NodeProvider = pg.NewNodeProvider(pool)

	currentState, err := sp.GetState()
	if err != nil {
		panic(err)
	}

	stateHolder := state.NewStateHolder(currentState)

	addrs := make(chan *address.Address, 16)
	newStates := make(chan *types.State, 16)

	if currentState.Address != nil {
		addrs <- currentState.Address.Address
	}

	go updates.Watcher(newStates, addrs, stateHolder, sp)

	var up updates.Recorder = &updates.FileUpdateRecorder{
		Base: path.Join(config.Config.DataDir, "upd"),
	}

	handler := &myhttp.Handler{
		StateProvider: sp,
		ItemProvider:  ip,
		NodeProvider:  np,

		StateHolder: stateHolder,

		Depth: config.Config.Depth,

		NewStates: newStates,
		Addresses: addrs,

		UpdateRecorder: up,
	}

	handler.RegisterHandlers(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Config.Port)))
}
