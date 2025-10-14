package task

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"taskbounty/x/task/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListTask",
					Use:       "list-task",
					Short:     "List all task",
				},
				{
					RpcMethod:      "GetTask",
					Use:            "get-task [id]",
					Short:          "Gets a task by id",
					Alias:          []string{"show-task"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateTask",
					Use:            "create-task [title] [description] [bounty] [status] [claimant] [proof] [approver]",
					Short:          "Create task",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "title"}, {ProtoField: "description"}, {ProtoField: "bounty"}, {ProtoField: "status"}, {ProtoField: "claimant"}, {ProtoField: "proof"}, {ProtoField: "approver"}},
				},
				{
					RpcMethod:      "UpdateTask",
					Use:            "update-task [id] [title] [description] [bounty] [status] [claimant] [proof] [approver]",
					Short:          "Update task",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "title"}, {ProtoField: "description"}, {ProtoField: "bounty"}, {ProtoField: "status"}, {ProtoField: "claimant"}, {ProtoField: "proof"}, {ProtoField: "approver"}},
				},
				{
					RpcMethod:      "DeleteTask",
					Use:            "delete-task [id]",
					Short:          "Delete task",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
