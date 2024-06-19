package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewWasmModule(cdc codec.Marshaler) wasm.AppModule {
	return wasm.NewAppModule(cdc, wasm.NewKeeper(cdc, sdk.NewKVStoreKey(wasm.StoreKey), sdk.NewKVStoreKey(wasm.TStoreKey), wasmTypes.DefaultWasmConfig()))
}
