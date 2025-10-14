package keeper

import (
	"context"
	"errors"
	"fmt"

	"taskbounty/x/task/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateTask(ctx context.Context, msg *types.MsgCreateTask) (*types.MsgCreateTaskResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	nextId, err := k.TaskSeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get next id")
	}

	var task = types.Task{
		Id:          nextId,
		Creator:     msg.Creator,
		Title:       msg.Title,
		Description: msg.Description,
		Bounty:      msg.Bounty,
		Status:      msg.Status,
		Claimant:    msg.Claimant,
		Proof:       msg.Proof,
		Approver:    msg.Approver,
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

	var task = types.Task{
		Creator:     msg.Creator,
		Id:          msg.Id,
		Title:       msg.Title,
		Description: msg.Description,
		Bounty:      msg.Bounty,
		Status:      msg.Status,
		Claimant:    msg.Claimant,
		Proof:       msg.Proof,
		Approver:    msg.Approver,
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

	if err := k.Task.Remove(ctx, msg.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete task")
	}

	return &types.MsgDeleteTaskResponse{}, nil
}
