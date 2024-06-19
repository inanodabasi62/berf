package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/inanodabasi62/projects/berf/x/mymodule/types"
)

type Keeper struct {
	cdc      codec.Marshaler
	storeKey sdk.StoreKey
}

func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

func (k Keeper) GetParams(ctx sdk.Context) types.GenesisState {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.ParamStoreKey))
	var gs types.GenesisState
	k.cdc.MustUnmarshalBinaryBare(b, &gs)
	return gs
}

func (k Keeper) SetParams(ctx sdk.Context, gs types.GenesisState) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(&gs)
	store.Set([]byte(types.ParamStoreKey), b)
}
