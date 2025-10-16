package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"taskbounty/x/task/types"
)

// GetTxCmd returns the transaction commands for the task module
func GetTaskTxCmd() *cobra.Command {
	taskTxCmd := &cobra.Command{
		Use:                        "task",
		Short:                      "Task transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	taskTxCmd.AddCommand(
		GetCmdCreateTask(),
		GetCmdUpdateTask(),
		GetCmdDeleteTask(),
		GetCmdClaimTask(),
		GetCmdSubmitTask(),
		GetCmdApproveTask(),
		GetCmdRejectTask(),
	)

	return taskTxCmd
}

// GetTaskQueryCmd returns the query commands for the task module
func GetTaskQueryCmd() *cobra.Command {
	taskQueryCmd := &cobra.Command{
		Use:                        "task",
		Short:                      "Task query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	taskQueryCmd.AddCommand(
		GetCmdQueryTask(),
		GetCmdQueryTasks(),
		GetCmdQueryTaskReward(),
		GetCmdQueryTaskRewards(),
		GetCmdQueryTaskRewardsByClaimant(),
	)

	return taskQueryCmd
}

// GetCmdCreateTask implements the create task command handler
func GetCmdCreateTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [title] [description] [bounty]",
		Short: "Create a new task with a bounty",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title := args[0]
			description := args[1]
			bountyStr := args[2]

			// Parse bounty
			bounty, err := sdk.ParseCoinNormalized(bountyStr)
			if err != nil {
				return fmt.Errorf("invalid bounty format: %v", err)
			}

			msg := types.NewMsgCreateTask(
				clientCtx.GetFromAddress().String(),
				title,
				description,
				bounty,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdUpdateTask implements the update task command handler
func GetCmdUpdateTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [id] [title] [description] [bounty]",
		Short: "Update an existing task",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task id: %v", err)
			}

			title := args[1]
			description := args[2]
			bountyStr := args[3]

			// Parse bounty
			bounty, err := sdk.ParseCoinNormalized(bountyStr)
			if err != nil {
				return fmt.Errorf("invalid bounty format: %v", err)
			}

			msg := types.NewMsgUpdateTask(
				clientCtx.GetFromAddress().String(),
				id,
				title,
				description,
				bounty,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdDeleteTask implements the delete task command handler
func GetCmdDeleteTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task id: %v", err)
			}

			msg := types.NewMsgDeleteTask(
				clientCtx.GetFromAddress().String(),
				id,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdClaimTask implements the claim task command handler
func GetCmdClaimTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [id]",
		Short: "Claim a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task id: %v", err)
			}

			msg := types.NewMsgClaimTask(
				clientCtx.GetFromAddress().String(),
				id,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdSubmitTask implements the submit task command handler
func GetCmdSubmitTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit [id] [proof-type] [proof-data]",
		Short: "Submit a completed task with proof",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task id: %v", err)
			}

			proofType := args[1]
			proofData := args[2]

			// Create a TaskProof
			proof := types.TaskProof{
				Type:      proofType,
				Data:      proofData,
				Timestamp: 0, // Will be set by the server
			}

			msg := types.NewMsgSubmitTask(
				clientCtx.GetFromAddress().String(),
				id,
				proof,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdApproveTask implements the approve task command handler
func GetCmdApproveTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve [id] [tx-hash]",
		Short: "Approve a submitted task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task id: %v", err)
			}

			txHash := args[1]

			msg := types.NewMsgApproveTask(
				clientCtx.GetFromAddress().String(),
				id,
				txHash,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdRejectTask implements the reject task command handler
func GetCmdRejectTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject [id] [reason]",
		Short: "Reject a submitted task",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task id: %v", err)
			}

			// Join all remaining arguments as the reason
			reason := strings.Join(args[1:], " ")

			msg := types.NewMsgRejectTask(
				clientCtx.GetFromAddress().String(),
				id,
				reason,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTask implements the query task command handler
func GetCmdQueryTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Query a task by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task id: %v", err)
			}

			res, err := queryClient.GetTask(cmd.Context(), &types.QueryGetTaskRequest{Id: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTasks implements the query tasks command handler
func GetCmdQueryTasks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Query all tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.ListTask(cmd.Context(), &types.QueryAllTaskRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "tasks")
	return cmd
}

// GetCmdQueryTaskReward implements the query task reward command handler
func GetCmdQueryTaskReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward [id]",
		Short: "Query a task reward by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task id: %v", err)
			}

			res, err := queryClient.GetTaskReward(cmd.Context(), &types.QueryGetTaskRewardRequest{Id: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTaskRewards implements the query task rewards command handler
func GetCmdQueryTaskRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards",
		Short: "Query all task rewards",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.ListTaskReward(cmd.Context(), &types.QueryAllTaskRewardRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "task-rewards")
	return cmd
}

// GetCmdQueryTaskRewardsByClaimant implements the query task rewards by claimant command handler
func GetCmdQueryTaskRewardsByClaimant() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards-by-claimant [claimant]",
		Short: "Query all task rewards for a specific claimant",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			claimant := args[0]

			res, err := queryClient.GetTaskRewardsByClaimant(cmd.Context(), &types.QueryGetTaskRewardsByClaimantRequest{Claimant: claimant})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}