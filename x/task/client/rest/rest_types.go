package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BaseReq struct {
	From          string       `json:"from"`
	Memo          string       `json:"memo"`
	ChainID       string       `json:"chain_id"`
	AccountNumber uint64       `json:"account_number"`
	Sequence      uint64       `json:"sequence"`
	Fees          sdk.Coins    `json:"fees"`
	GasPrices     sdk.DecCoins `json:"gas_prices"`
	Gas           string       `json:"gas"`
	GasAdjustment string       `json:"gas_adjustment"`
	Simulate      bool         `json:"simulate"`
}

func BaseReqFromRequest(r *http.Request) BaseReq {
	return BaseReq{}
}

func (b BaseReq) ValidateBasic(w http.ResponseWriter) bool {
	return true
}

func ReadRESTReq(w http.ResponseWriter, r *http.Request, cdc interface{}, req interface{}) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to read request body: %s", err))
		return false
	}

	err = json.Unmarshal(body, req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to unmarshal request: %s", err))
		return false
	}

	return true
}

func WriteErrorResponse(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
}

func PostProcessResponse(w http.ResponseWriter, ctx client.Context, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	respBytes, err := json.Marshal(resp)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal response: %s", err))
		return
	}
	
	_, _ = w.Write(respBytes)
}

func ParseUint64OrReturnBadRequest(w http.ResponseWriter, s string) (uint64, bool) {
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid ID: %s", err))
		return 0, false
	}
	return u, true
}

type TxResponseGenerator struct {
	ClientCtx client.Context
	TxBuilder client.TxBuilder
}

func (g TxResponseGenerator) BuildTx(baseReq BaseReq, msgs ...sdk.Msg) (client.TxBuilder, error) {
	return g.TxBuilder, nil
}

func (g TxResponseGenerator) FinalizeTx(baseReq BaseReq, txBuilder client.TxBuilder) interface{} {
	return nil
}