package keeper

import (
	"context"
	"errors"
	"fmt"

	"taskbounty/x/task/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateTask(ctx context.Context, msg *types.MsgCreateTask) (*types.MsgCreateTaskResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Get module parameters
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get params")
	}

	nextId, err := k.TaskSeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get next id")
	}

	// Get current timestamp
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime().Unix()

	var task = types.Task{
		Id:          nextId,
		Creator:     msg.Creator,
		Title:       msg.Title,
		Description: msg.Description,
		Bounty:      msg.Bounty,
		Status:      types.TASK_STATUS_OPEN, // Default to open status
		Claimant:    "",
		Proof:       "",
		Approver:    "",
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}

	// Validate the task
	if err := task.Validate(params); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if err = k.Task.Set(
		ctx,
		nextId,
		task,
	); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to set task")
	}

	return &types.MsgCreateTaskResponse{
		Id: nextId,
	}, nil
}

func (k msgServer) UpdateTask(ctx context.Context, msg *types.MsgUpdateTask) (*types.MsgUpdateTaskResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Get module parameters
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get params")
	}

	// Checks that the element exists
	val, err := k.Task.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get task")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// Get current timestamp
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime().Unix()

	// Create updated task with new values
	task := types.Task{
		Creator:     msg.Creator,
		Id:          msg.Id,
		Title:       msg.Title,
		Description: msg.Description,
		Bounty:      msg.Bounty,
		Status:      msg.Status,
		Claimant:    msg.Claimant,
		Proof:       msg.Proof,
		Approver:    msg.Approver,
		CreatedAt:   val.CreatedAt,
		UpdatedAt:   currentTime,
	}

	// Validate the status transition
	if !types.IsValidTransition(val.Status, task.Status) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid status transition from %s to %s", types.TaskStatusToString(val.Status), types.TaskStatusToString(task.Status)))
	}

	// Validate the updated task
	if err := task.Validate(params); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if err := k.Task.Set(ctx, msg.Id, task); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update task")
	}

	return &types.MsgUpdateTaskResponse{}, nil
}

func (k msgServer) DeleteTask(ctx context.Context, msg *types.MsgDeleteTask) (*types.MsgDeleteTaskResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Checks that the element exists
	val, err := k.Task.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get task")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// Check if task can be deleted (only in open or closed status)
	if val.Status != types.TASK_STATUS_OPEN && val.Status != types.TASK_STATUS_CLOSED {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("cannot delete task in %s status", types.TaskStatusToString(val.Status)))
	}

	if err := k.Task.Remove(ctx, msg.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete task")
	}

	return &types.MsgDeleteTaskResponse{}, nil
}
