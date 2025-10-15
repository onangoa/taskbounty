package keeper

import (
"context"
"fmt"
"time"

"taskbounty/x/task/types"

"cosmossdk.io/collections"
errorsmod "cosmossdk.io/errors"
sdk "github.com/cosmos/cosmos-sdk/types"
sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// claiming of a task by a user
func (k msgServer) ClaimTask(ctx context.Context, msg *types.MsgClaimTask) (*types.MsgClaimTaskResponse, error) {
if _, err := k.addressCodec.StringToBytes(msg.Claimant); err != nil {
return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
}
// Get task
task, err := k.Task.Get(ctx, msg.Id)
if err != nil {
	if errors.Is(err, collections.ErrNotFound) {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("task %d not found", msg.Id))
	}
	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get task")
}


}
