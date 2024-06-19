package evm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/evm"
	evmTypes "github.com/tharsis/ethermint/x/evm/types"
)

func NewEvmModule(cdc codec.Marshaler) evm.AppModule {
	return evm.NewAppModule(cdc, evm.NewKeeper(cdc, sdk.NewKVStoreKey(evm.StoreKey), sdk.NewKVStoreKey(evm.TStoreKey), evmTypes.DefaultParams()))
}
