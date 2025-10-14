package keeper_test

import (
	"context"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"taskbounty/x/task/keeper"
	"taskbounty/x/task/types"
)

func createNTask(keeper keeper.Keeper, ctx context.Context, n int) []types.Task {
	items := make([]types.Task, n)
	for i := range items {
		iu := uint64(i)
		items[i].Id = iu
		items[i].Title = strconv.Itoa(i)
		items[i].Description = strconv.Itoa(i)
		items[i].Bounty = sdk.NewInt64Coin(`token`, int64(i+100))
		items[i].Status = int64(i)
		items[i].Claimant = strconv.Itoa(i)
		items[i].Proof = strconv.Itoa(i)
		items[i].Approver = strconv.Itoa(i)
		_ = keeper.Task.Set(ctx, iu, items[i])
		_ = keeper.TaskSeq.Set(ctx, iu)
	}
	return items
}

func TestTaskQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNTask(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetTaskRequest
		response *types.QueryGetTaskResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetTaskRequest{Id: msgs[0].Id},
			response: &types.QueryGetTaskResponse{Task: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetTaskRequest{Id: msgs[1].Id},
			response: &types.QueryGetTaskResponse{Task: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetTaskRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetTask(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestTaskQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNTask(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllTaskRequest {
		return &types.QueryAllTaskRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListTask(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Task), step)
			require.Subset(t, msgs, resp.Task)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListTask(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Task), step)
			require.Subset(t, msgs, resp.Task)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListTask(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Task)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListTask(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
