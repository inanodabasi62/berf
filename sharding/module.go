package sharding

import (
	"encoding/json"

	"cosmossdk.io/api/tendermint/abci"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func NewShardingModule(cdc codec.Marshaler) ShardingModule {
	return ShardingModule{cdc: cdc}
}

type ShardingModule struct {
	cdc codec.Marshaler
}

func (am ShardingModule) RegisterInvariants(ir sdk.InvariantRegistry) {}
func (am ShardingModule) Route() string                               { return "sharding" }
func (am ShardingModule) NewHandler() sdk.Handler                     { return NewHandler() }
func (am ShardingModule) QuerierRoute() string                        { return "sharding" }
func (am ShardingModule) NewQuerierHandler() sdk.Querier              { return NewQuerier() }
func (am ShardingModule) RegisterServices(cfg module.Configurator)    {}
func (am ShardingModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
func (am ShardingModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return nil
}
