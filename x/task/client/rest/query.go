package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
	
	"taskbounty/x/task/types"
)

// RegisterRoutes registers the REST routes for the task module.
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerQueryRoutes(clientCtx, r)
	registerTxRoutes(clientCtx, r)
}

// registerQueryRoutes registers the REST routes for the query service.
func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/taskbounty/task/v1/params", queryParamsHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/taskbounty/task/v1/task", listTaskHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/taskbounty/task/v1/task/{id}", getTaskHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/taskbounty/task/v1/task_reward", listTaskRewardHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/taskbounty/task/v1/task_reward/{id}", getTaskRewardHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/taskbounty/task/v1/task_rewards/{claimant}", getTaskRewardsByClaimantHandler(clientCtx)).Methods("GET")
}

// registerTxRoutes registers the REST routes for the transaction service.
func registerTxRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/taskbounty/task/v1/task", createTaskHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/taskbounty/task/v1/task/{id}", updateTaskHandler(clientCtx)).Methods("PUT")
	r.HandleFunc("/taskbounty/task/v1/task/{id}", deleteTaskHandler(clientCtx)).Methods("DELETE")
	r.HandleFunc("/taskbounty/task/v1/task/{id}/claim", claimTaskHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/taskbounty/task/v1/task/{id}/submit", submitTaskHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/taskbounty/task/v1/task/{id}/approve", approveTaskHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/taskbounty/task/v1/task/{id}/reject", rejectTaskHandler(clientCtx)).Methods("POST")
}

// queryParamsHandler returns the module parameters
func queryParamsHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := clientCtx.QueryWithData("custom/task/v1/params", nil)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		PostProcessResponse(w, clientCtx, res)
	}
}

// listTaskHandler returns a list of all tasks
func listTaskHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := clientCtx.QueryWithData("custom/task/v1/task", nil)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		PostProcessResponse(w, clientCtx, res)
	}
}

// getTaskHandler returns a specific task
func getTaskHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/task/v1/task/%s", id), nil)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		PostProcessResponse(w, clientCtx, res)
	}
}
// listTaskRewardHandler returns a list of all task rewards
func listTaskRewardHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := clientCtx.QueryWithData("custom/task/v1/task_reward", nil)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		PostProcessResponse(w, clientCtx, res)
	}
}

// getTaskRewardHandler returns a specific task reward
func getTaskRewardHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/task/v1/task_reward/%s", id), nil)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		PostProcessResponse(w, clientCtx, res)
	}
}

// getTaskRewardsByClaimantHandler returns all task rewards for a given claimant
func getTaskRewardsByClaimantHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		claimant := vars["claimant"]
		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/task/v1/task_rewards/%s", claimant), nil)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		PostProcessResponse(w, clientCtx, res)
	}
}

// Task transaction handlers

func createTaskHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MsgCreateTask
		if !ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := BaseReqFromRequest(r)
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgCreateTask(req.Creator, req.Title, req.Description, req.Bounty)
		txWrite := TxResponseGenerator{
			ClientCtx: clientCtx,
			TxBuilder: clientCtx.TxConfig.NewTxBuilder(),
		}

		tx, err := txWrite.BuildTx(baseReq, msg)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		response := txWrite.FinalizeTx(baseReq, tx)
		PostProcessResponse(w, clientCtx, response)
	}
}

func updateTaskHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req types.MsgUpdateTask
		if !ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := BaseReqFromRequest(r)
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Parse task ID from URL
		taskID, ok := ParseUint64OrReturnBadRequest(w, id)
		if !ok {
			return
		}
		req.Id = taskID

		msg := types.NewMsgUpdateTask(req.Creator, req.Id, req.Title, req.Description, req.Bounty)
		txWrite := TxResponseGenerator{
			ClientCtx: clientCtx,
			TxBuilder: clientCtx.TxConfig.NewTxBuilder(),
		}

		tx, err := txWrite.BuildTx(baseReq, msg)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		response := txWrite.FinalizeTx(baseReq, tx)
		PostProcessResponse(w, clientCtx, response)
	}
}

func deleteTaskHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req types.MsgDeleteTask
		if !ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := BaseReqFromRequest(r)
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Parse task ID from URL
		taskID, ok := ParseUint64OrReturnBadRequest(w, id)
		if !ok {
			return
		}
		req.Id = taskID

		msg := types.NewMsgDeleteTask(req.Creator, req.Id)
		txWrite := TxResponseGenerator{
			ClientCtx: clientCtx,
			TxBuilder: clientCtx.TxConfig.NewTxBuilder(),
		}

		tx, err := txWrite.BuildTx(baseReq, msg)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		response := txWrite.FinalizeTx(baseReq, tx)
		PostProcessResponse(w, clientCtx, response)
	}
}

func claimTaskHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req types.MsgClaimTask
		if !ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := BaseReqFromRequest(r)
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Parse task ID from URL
		taskID, ok := ParseUint64OrReturnBadRequest(w, id)
		if !ok {
			return
		}
		req.Id = taskID

		msg := types.NewMsgClaimTask(req.Claimant, req.Id)
		txWrite := TxResponseGenerator{
			ClientCtx: clientCtx,
			TxBuilder: clientCtx.TxConfig.NewTxBuilder(),
		}

		tx, err := txWrite.BuildTx(baseReq, msg)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		response := txWrite.FinalizeTx(baseReq, tx)
		PostProcessResponse(w, clientCtx, response)
	}
}

func submitTaskHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req types.MsgSubmitTask
		if !ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := BaseReqFromRequest(r)
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Parse task ID from URL
		taskID, ok := ParseUint64OrReturnBadRequest(w, id)
		if !ok {
			return
		}
		req.Id = taskID

		msg := types.NewMsgSubmitTask(req.Claimant, req.Id, req.Proof)
		txWrite := TxResponseGenerator{
			ClientCtx: clientCtx,
			TxBuilder: clientCtx.TxConfig.NewTxBuilder(),
		}

		tx, err := txWrite.BuildTx(baseReq, msg)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		response := txWrite.FinalizeTx(baseReq, tx)
		PostProcessResponse(w, clientCtx, response)
	}
}

func approveTaskHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req types.MsgApproveTask
		if !ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := BaseReqFromRequest(r)
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Parse task ID from URL
		taskID, ok := ParseUint64OrReturnBadRequest(w, id)
		if !ok {
			return
		}
		req.Id = taskID

		msg := types.NewMsgApproveTask(req.Approver, req.Id, req.TxHash)
		txWrite := TxResponseGenerator{
			ClientCtx: clientCtx,
			TxBuilder: clientCtx.TxConfig.NewTxBuilder(),
		}

		tx, err := txWrite.BuildTx(baseReq, msg)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		response := txWrite.FinalizeTx(baseReq, tx)
		PostProcessResponse(w, clientCtx, response)
	}
}

func rejectTaskHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req types.MsgRejectTask
		if !ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := BaseReqFromRequest(r)
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Parse task ID from URL
		taskID, ok := ParseUint64OrReturnBadRequest(w, id)
		if !ok {
			return
		}
		req.Id = taskID

		msg := types.NewMsgRejectTask(req.Rejecter, req.Id, req.Reason)
		txWrite := TxResponseGenerator{
			ClientCtx: clientCtx,
			TxBuilder: clientCtx.TxConfig.NewTxBuilder(),
		}

		tx, err := txWrite.BuildTx(baseReq, msg)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		response := txWrite.FinalizeTx(baseReq, tx)
		PostProcessResponse(w, clientCtx, response)
	}
}