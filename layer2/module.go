package layer2

import (
	"encoding/json"

	"cosmossdk.io/api/tendermint/abci"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func NewLayer2Module(cdc codec.Marshaler) Layer2Module {
	return Layer2Module{cdc: cdc}
}

type Layer2Module struct {
	cdc codec.Marshaler
}

func (am Layer2Module) RegisterInvariants(ir sdk.InvariantRegistry) {}
func (am Layer2Module) Route() string                               { return "layer2" }
func (am Layer2Module) NewHandler() sdk.Handler                     { return NewHandler() }
func (am Layer2Module) QuerierRoute() string                        { return "layer2" }
func (am Layer2Module) NewQuerierHandler() sdk.Querier              { return NewQuerier() }
func (am Layer2Module) RegisterServices(cfg module.Configurator)    {}
func (am Layer2Module) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
func (am Layer2Module) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return nil
}
