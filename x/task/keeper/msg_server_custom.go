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

if err := task.CanClaim(msg.Claimant); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
}

params, err := k.Params.Get(ctx)
if err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get params")
}

currentTime := sdk.UnixContext(ctx).Unix()
if task.IsExpired(params, time.Unix(currentTime, 0)) {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "task has expired")
}

task.Claimant = msg.Claimant
task.Status = types.TASK_STATUS_CLAIMED
task.UpdatedAt = currentTime

// Save the updated task
if err := k.Task.Set(ctx, task.Id, task); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update task")
}

return &types.MsgClaimTaskResponse{}, nil
}

func (k msgServer) SubmitTask(ctx context.Context, msg *types.MsgSubmitTask) (*types.MsgSubmitTaskResponse, error) {
if _, err := k.addressCodec.StringToBytes(msg.Claimant); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
}

task, err := k.Task.Get(ctx, msg.Id)
if err != nil {
	if errors.Is(err, collections.ErrNotFound) {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("task %d not found", msg.Id))
	}
	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get task")
}

if err := task.CanSubmit(msg.Claimant); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
}

params, err := k.Params.Get(ctx)
if err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get params")
}

if err := msg.Proof.Validate(params); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
}

currentTime := sdk.UnixContext(ctx).Unix()
if task.IsClaimExpired(params, time.Unix(currentTime, 0)) {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "task claim has expired")
}

proofStr := fmt.Sprintf("%s:%s:%d", msg.Proof.Hash, msg.Proof.Type, msg.Proof.Timestamp)
if msg.Proof.Data != "" {
	proofStr += fmt.Sprintf(":%s", msg.Proof.Data)
}

task.Proof = proofStr
task.Status = types.TASK_STATUS_SUBMITTED
task.UpdatedAt = currentTime

if err := k.Task.Set(ctx, task.Id, task); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update task")
}

return &types.MsgSubmitTaskResponse{}, nil
}

func (k msgServer) ApproveTask(ctx context.Context, msg *types.MsgApproveTask) (*types.MsgApproveTaskResponse, error) {
if _, err := k.addressCodec.StringToBytes(msg.Approver); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
}

task, err := k.Task.Get(ctx, msg.Id)
if err != nil {
	if errors.Is(err, collections.ErrNotFound) {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("task %d not found", msg.Id))
	}
	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get task")
}

if err := task.CanApprove(msg.Approver); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
}
params, err := k.Params.Get(ctx)
if err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get params")
}
currentTime := sdk.UnixContext(ctx).Unix()
if task.IsSubmissionExpired(params, time.Unix(currentTime, 0)) {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "task submission has expired")
}

reward := types.CreateTaskReward(task.Id, task.Claimant, task.Bounty, msg.TxHash, currentTime)
if err := types.ValidateRewardDistribution(task, reward); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
}

task.Approver = msg.Approver
task.Status = types.TASK_STATUS_APPROVED
task.UpdatedAt = currentTime

if err := k.Task.Set(ctx, task.Id, task); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update task")
}

if err := k.TaskReward.Set(ctx, task.Id, reward); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to store task reward")
}

return &types.MsgApproveTaskResponse{}, nil
}

func (k msgServer) RejectTask(ctx context.Context, msg *types.MsgRejectTask) (*types.MsgRejectTaskResponse, error) {
if _, err := k.addressCodec.StringToBytes(msg.Rejecter); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
}
task, err := k.Task.Get(ctx, msg.Id)
if err != nil {
	if errors.Is(err, collections.ErrNotFound) {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("task %d not found", msg.Id))
	}
	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get task")
}
if err := task.CanReject(msg.Rejecter); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
}
params, err := k.Params.Get(ctx)
if err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get params")
}
currentTime := sdk.UnixContext(ctx).Unix()
if task.IsSubmissionExpired(params, time.Unix(currentTime, 0)) {
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "task submission has expired")
}

task.Proof = fmt.Sprintf("REJECTED: %s", msg.Reason)
task.Status = types.TASK_STATUS_REJECTED
task.UpdatedAt = currentTime

if err := k.Task.Set(ctx, task.Id, task); err != nil {
	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update task")
}

return &types.MsgRejectTaskResponse{}, nil
}
