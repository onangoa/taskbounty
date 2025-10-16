package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{
		MinBounty:             sdk.NewCoin("stake", math.NewInt(1000)),
		MaxBounty:             sdk.NewCoin("stake", math.NewInt(1000000)),
		MaxTitleLength:        100,
		MaxDescriptionLength:  1000,
		ProofTypes:            []string{"ipfs", "url", "text"},
		AutoApproveThreshold:  5,
		TaskExpiry:            86400 * 30,
		ClaimDeadline:         86400 * 7,
		SubmissionDeadline:    86400 * 14,
	}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.MinBounty.IsZero() || p.MinBounty.IsNegative() {
		return fmt.Errorf("min bounty must be positive")
	}
	if p.MaxBounty.IsZero() || p.MaxBounty.IsNegative() {
		return fmt.Errorf("max bounty must be positive")
	}
	if p.MinBounty.Amount.GT(p.MaxBounty.Amount) {
		return fmt.Errorf("min bounty cannot be greater than max bounty")
	}
	if p.MinBounty.Denom != p.MaxBounty.Denom {
		return fmt.Errorf("min and max bounty must have the same denom")
	}
	if p.MaxTitleLength == 0 {
		return fmt.Errorf("max title length must be positive")
	}
	if p.MaxDescriptionLength == 0 {
		return fmt.Errorf("max description length must be positive")
	}
	if len(p.ProofTypes) == 0 {
		return fmt.Errorf("at least one proof type must be specified")
	}
	if p.AutoApproveThreshold == 0 {
		return fmt.Errorf("auto approve threshold must be positive")
	}

	return nil
}
