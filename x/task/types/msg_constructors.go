package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgCreateTask(creator string, title string, description string, bounty sdk.Coin) *MsgCreateTask {
	return &MsgCreateTask{
		Creator:     creator,
		Title:       title,
		Description: description,
		Bounty:      bounty,
	}
}

func NewMsgUpdateTask(creator string, id uint64, title string, description string, bounty sdk.Coin) *MsgUpdateTask {
	return &MsgUpdateTask{
		Creator:     creator,
		Id:          id,
		Title:       title,
		Description: description,
		Bounty:      bounty,
	}
}

func NewMsgDeleteTask(creator string, id uint64) *MsgDeleteTask {
	return &MsgDeleteTask{
		Creator: creator,
		Id:      id,
	}
}

func NewMsgClaimTask(claimant string, id uint64) *MsgClaimTask {
	return &MsgClaimTask{
		Claimant: claimant,
		Id:       id,
	}
}

func NewMsgSubmitTask(claimant string, id uint64, proof TaskProof) *MsgSubmitTask {
	return &MsgSubmitTask{
		Claimant: claimant,
		Id:       id,
		Proof:    proof,
	}
}

func NewMsgApproveTask(approver string, id uint64, txHash string) *MsgApproveTask {
	return &MsgApproveTask{
		Approver: approver,
		Id:       id,
		TxHash:   txHash,
	}
}

func NewMsgRejectTask(rejecter string, id uint64, reason string) *MsgRejectTask {
	return &MsgRejectTask{
		Rejecter: rejecter,
		Id:       id,
		Reason:   reason,
	}
}