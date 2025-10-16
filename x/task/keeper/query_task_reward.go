package keeper

import (
	"context"
	"errors"

	"taskbounty/x/task/types"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) GetTaskReward(ctx context.Context, req *types.QueryGetTaskRewardRequest) (*types.QueryGetTaskRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	reward, err := q.k.TaskReward.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetTaskRewardResponse{TaskReward: reward}, nil
}

func (q queryServer) ListTaskReward(ctx context.Context, req *types.QueryAllTaskRewardRequest) (*types.QueryAllTaskRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	rewards, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.TaskReward,
		req.Pagination,
		func(_ uint64, value types.TaskReward) (types.TaskReward, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTaskRewardResponse{TaskReward: rewards, Pagination: pageRes}, nil
}

func (q queryServer) GetTaskRewardsByClaimant(ctx context.Context, req *types.QueryGetTaskRewardsByClaimantRequest) (*types.QueryGetTaskRewardsByClaimantResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var rewards []types.TaskReward

	// loop all rewards and filter by claimant
	err := q.k.TaskReward.Walk(ctx, nil, func(key uint64, reward types.TaskReward) (bool, error) {
		if reward.Claimant == req.Claimant {
			rewards = append(rewards, reward)
		}
		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetTaskRewardsByClaimantResponse{TaskRewards: rewards}, nil
}
