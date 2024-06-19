package mymodule

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/inanodabasi62/projects/berf/x/mymodule/keeper"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		// Handle your module messages here.
		default:
			return nil, errors.Wrap(errors.ErrUnknownRequest, "unrecognized mymodule message type")
		}
	}
}
