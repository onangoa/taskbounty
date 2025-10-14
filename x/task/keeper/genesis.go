package keeper

import (
	"context"

	"taskbounty/x/task/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.TaskList {
		if err := k.Task.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.TaskSeq.Set(ctx, genState.TaskCount); err != nil {
		return err
	}
	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = k.Task.Walk(ctx, nil, func(key uint64, elem types.Task) (bool, error) {
		genesis.TaskList = append(genesis.TaskList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.TaskCount, err = k.TaskSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
