package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"taskbounty/x/task/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema  collections.Schema
	Params  collections.Item[types.Params]
	TaskSeq collections.Sequence
	Task    collections.Map[uint64, types.Task]
	TaskReward collections.Map[uint64, types.TaskReward]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params:  collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Task:    collections.NewMap(sb, types.TaskKey, "task", collections.Uint64Key, codec.CollValue[types.Task](cdc)),
		TaskSeq: collections.NewSequence(sb, types.TaskCountKey, "taskSequence"),
		TaskReward: collections.NewMap(sb, collections.NewPrefix(1), "task_reward", collections.Uint64Key, codec.CollValue[types.TaskReward](cdc)),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
