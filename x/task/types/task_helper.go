package types

import (
"fmt"
"strings"
"time"

sdk "github.com/cosmos/cosmos-sdk/types"
)

// converts a TaskStatus to its string representation
func TaskStatusToString(status TaskStatus) string {
switch status {
case TASK_STATUS_UNDEFINED:
return "undefined"
case TASK_STATUS_OPEN:
return "open"
case TASK_STATUS_CLAIMED:
return "claimed"
case TASK_STATUS_SUBMITTED:
return "submitted"
case TASK_STATUS_APPROVED:
return "approved"
case TASK_STATUS_REJECTED:
return "rejected"
case TASK_STATUS_CLOSED:
return "closed"
default:
return "unknown"
}
}

// converts a string to a TaskStatus
func StringToTaskStatus(status string) TaskStatus {
switch strings.ToLower(status) {
case "undefined":
return TASK_STATUS_UNDEFINED
case "open":
return TASK_STATUS_OPEN
case "claimed":
return TASK_STATUS_CLAIMED
case "submitted":
return TASK_STATUS_SUBMITTED
case "approved":
return TASK_STATUS_APPROVED
case "rejected":
return TASK_STATUS_REJECTED
case "closed":
return TASK_STATUS_CLOSED
default:
return TASK_STATUS_UNDEFINED
}
}

// checks if the given status is a valid TaskStatus
func IsValidTaskStatus(status TaskStatus) bool {
return status >= TASK_STATUS_UNDEFINED && status <= TASK_STATUS_CLOSED
}

// list of valid status transitions
func GetValidTransitions() []TaskTransition {
return []TaskTransition{
{From: TASK_STATUS_UNDEFINED, To: TASK_STATUS_OPEN},
{From: TASK_STATUS_OPEN, To: TASK_STATUS_CLAIMED},
{From: TASK_STATUS_OPEN, To: TASK_STATUS_CLOSED},
{From: TASK_STATUS_CLAIMED, To: TASK_STATUS_SUBMITTED},
{From: TASK_STATUS_CLAIMED, To: TASK_STATUS_OPEN}, // -> Revert claim
{From: TASK_STATUS_SUBMITTED, To: TASK_STATUS_APPROVED},
{From: TASK_STATUS_SUBMITTED, To: TASK_STATUS_REJECTED},
{From: TASK_STATUS_REJECTED, To: TASK_STATUS_CLAIMED}, // when resubmited
{From: TASK_STATUS_REJECTED, To: TASK_STATUS_OPEN},    // when reopened
{From: TASK_STATUS_APPROVED, To: TASK_STATUS_CLOSED}, // when closed after approval
}
}

func IsValidTransition(from, to TaskStatus) bool {
for _, transition := range GetValidTransitions() {
if transition.From == from && transition.To == to {
return true
}
}
return false
}

func GetTransitionsFrom(status TaskStatus) []TaskStatus {
var transitions []TaskStatus
for _, transition := range GetValidTransitions() {
if transition.From == status {
transitions = append(transitions, transition.To)
}
}
return transitions
}

func (t Task) Validate(params Params) error {
if strings.TrimSpace(t.Title) == "" {
return fmt.Errorf("title cannot be empty")
}
if uint32(len(t.Title)) > params.MaxTitleLength {
return fmt.Errorf("title exceeds maximum length of %d", params.MaxTitleLength)
}

	if strings.TrimSpace(t.Description) == "" {
		return fmt.Errorf("description cannot be empty")
	}
	if uint32(len(t.Description)) > params.MaxDescriptionLength {
		return fmt.Errorf("description exceeds maximum length of %d", params.MaxDescriptionLength)
	}

	if t.Bounty.IsZero() {
		return fmt.Errorf("bounty cannot be zero")
	}
	if t.Bounty.IsNegative() {
		return fmt.Errorf("bounty cannot be negative")
	}
	if t.Bounty.Amount.LT(params.MinBounty.Amount) {
		return fmt.Errorf("bounty amount is below minimum of %s", params.MinBounty.String())
	}
	if !params.MaxBounty.IsZero() && t.Bounty.Amount.GT(params.MaxBounty.Amount) {
		return fmt.Errorf("bounty amount exceeds maximum of %s", params.MaxBounty.String())
	}
	if t.Bounty.Denom != params.MinBounty.Denom {
		return fmt.Errorf("bounty denom must be %s", params.MinBounty.Denom)
	}
	if !IsValidTaskStatus(t.Status) {
		return fmt.Errorf("invalid task status: %s", TaskStatusToString(t.Status))
	}
	if strings.TrimSpace(t.Creator) == "" {
		return fmt.Errorf("creator cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(t.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %s", err)
	}
	if t.CreatedAt <= 0 {
		return fmt.Errorf("created_at must be positive")
	}
	if t.UpdatedAt < t.CreatedAt {
		return fmt.Errorf("updated_at cannot be before created_at")
	}
	if strings.TrimSpace(t.Claimant) != "" {
		if _, err := sdk.AccAddressFromBech32(t.Claimant); err != nil {
			return fmt.Errorf("invalid claimant address: %s", err)
		}
	}
	if strings.TrimSpace(t.Approver) != "" {
		if _, err := sdk.AccAddressFromBech32(t.Approver); err != nil {
			return fmt.Errorf("invalid approver address: %s", err)
		}
	}

	return nil
}

func (t Task) CanClaim(claimant string) error {
	if t.Status != TASK_STATUS_OPEN {
		return fmt.Errorf("task is not open for claiming")
	}
	if t.Creator == claimant {
		return fmt.Errorf("creator cannot claim their own task")
	}
	if strings.TrimSpace(t.Claimant) != "" {
		return fmt.Errorf("task is already claimed by %s", t.Claimant)
	}

	return nil
}

func (t Task) CanSubmit(claimer string) error {
	if t.Status != TASK_STATUS_CLAIMED {
		return fmt.Errorf("task is not in claimed status")
	}
	if t.Claimant != claimer {
		return fmt.Errorf("only the current claimant can submit the task")
	}

	return nil
}

func (t Task) CanApprove(approver string) error {
	if t.Status != TASK_STATUS_SUBMITTED {
		return fmt.Errorf("task is not in submitted status")
	}
	if t.Creator != approver {
		return fmt.Errorf("only the creator can approve the task")
	}

	return nil
}

func (t Task) CanReject(rejecter string) error {
	if t.Status != TASK_STATUS_SUBMITTED {
		return fmt.Errorf("task is not in submitted status")
	}
	if t.Creator != rejecter {
		return fmt.Errorf("only the creator can reject the task")
	}

	return nil
}

func (t Task) IsExpired(params Params, currentTime time.Time) bool {
	if params.TaskExpiry == 0 {
		return false
	}

	expiryTime := time.Unix(t.CreatedAt, 0).Add(time.Duration(params.TaskExpiry) * time.Second)
	return currentTime.After(expiryTime)
}

func (t Task) IsClaimExpired(params Params, currentTime time.Time) bool {
	if params.ClaimDeadline == 0 || t.Claimant == "" {
		return false
	}
	claimTime := time.Unix(t.UpdatedAt, 0)
	expiryTime := claimTime.Add(time.Duration(params.ClaimDeadline) * time.Second)
	return currentTime.After(expiryTime)
}

func (t Task) IsSubmissionExpired(params Params, currentTime time.Time) bool {
	if params.SubmissionDeadline == 0 || t.Claimant == "" {
		return false
	}
	claimTime := time.Unix(t.UpdatedAt, 0)
	expiryTime := claimTime.Add(time.Duration(params.SubmissionDeadline) * time.Second)
	return currentTime.After(expiryTime)
}

func FilterTasks(tasks []Task, filter TaskFilter) []Task {
	var filteredTasks []Task

	for _, task := range tasks {
		if filter.Creator != "" && task.Creator != filter.Creator {
			continue
		}
		if filter.Claimant != "" && task.Claimant != filter.Claimant {
			continue
		}
		if filter.Approver != "" && task.Approver != filter.Approver {
			continue
		}
		if filter.Status != TASK_STATUS_UNDEFINED && task.Status != filter.Status {
			continue
		}
		if !filter.MinBounty.IsZero() && task.Bounty.Amount.LT(filter.MinBounty.Amount) {
			continue
		}
		if !filter.MaxBounty.IsZero() && task.Bounty.Amount.GT(filter.MaxBounty.Amount) {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}

	return filteredTasks
}

func SortTasks(tasks []Task, sort TaskSort) []Task {
	sortedTasks := make([]Task, len(tasks))
	copy(sortedTasks, tasks)

	switch sort.Field {
	case "id":
		if sort.Direction == "desc" {
			for i := 0; i < len(sortedTasks); i++ {
				for j := i + 1; j < len(sortedTasks); j++ {
					if sortedTasks[i].Id < sortedTasks[j].Id {
						sortedTasks[i], sortedTasks[j] = sortedTasks[j], sortedTasks[i]
					}
				}
			}
		} else {
			for i := 0; i < len(sortedTasks); i++ {
				for j := i + 1; j < len(sortedTasks); j++ {
					if sortedTasks[i].Id > sortedTasks[j].Id {
						sortedTasks[i], sortedTasks[j] = sortedTasks[j], sortedTasks[i]
					}
				}
			}
		}
	case "bounty":
		if sort.Direction == "desc" {
			for i := 0; i < len(sortedTasks); i++ {
				for j := i + 1; j < len(sortedTasks); j++ {
					if sortedTasks[i].Bounty.Amount.LT(sortedTasks[j].Bounty.Amount) {
						sortedTasks[i], sortedTasks[j] = sortedTasks[j], sortedTasks[i]
					}
				}
			}
		} else {
			for i := 0; i < len(sortedTasks); i++ {
				for j := i + 1; j < len(sortedTasks); j++ {
					if sortedTasks[i].Bounty.Amount.GT(sortedTasks[j].Bounty.Amount) {
						sortedTasks[i], sortedTasks[j] = sortedTasks[j], sortedTasks[i]
					}
				}
			}
		}
	case "status":
		if sort.Direction == "desc" {
			for i := 0; i < len(sortedTasks); i++ {
				for j := i + 1; j < len(sortedTasks); j++ {
					if sortedTasks[i].Status < sortedTasks[j].Status {
						sortedTasks[i], sortedTasks[j] = sortedTasks[j], sortedTasks[i]
					}
				}
			}
		} else {
			for i := 0; i < len(sortedTasks); i++ {
				for j := i + 1; j < len(sortedTasks); j++ {
					if sortedTasks[i].Status > sortedTasks[j].Status {
						sortedTasks[i], sortedTasks[j] = sortedTasks[j], sortedTasks[i]
					}
				}
			}
		}
	case "created_at":
		if sort.Direction == "desc" {
			for i := 0; i < len(sortedTasks); i++ {
				for j := i + 1; j < len(sortedTasks); j++ {
					if sortedTasks[i].CreatedAt < sortedTasks[j].CreatedAt {
						sortedTasks[i], sortedTasks[j] = sortedTasks[j], sortedTasks[i]
					}
				}
			}
		} else {
			for i := 0; i < len(sortedTasks); i++ {
				for j := i + 1; j < len(sortedTasks); j++ {
					if sortedTasks[i].CreatedAt > sortedTasks[j].CreatedAt {
						sortedTasks[i], sortedTasks[j] = sortedTasks[j], sortedTasks[i]
					}
				}
			}
		}
	}

	return sortedTasks
}

func (p TaskProof) Validate(params Params) error {
	if strings.TrimSpace(p.Hash) == "" {
		return fmt.Errorf("proof hash cannot be empty")
	}
	if strings.TrimSpace(p.Type) == "" {
		return fmt.Errorf("proof type cannot be empty")
	}
	allowed := false
	for _, proofType := range params.ProofTypes {
		if proofType == p.Type {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("proof type %s is not allowed", p.Type)
	}
	if p.Timestamp <= 0 {
		return fmt.Errorf("proof timestamp must be positive")
	}

	return nil
}

func (r TaskReward) Validate() error {
	if r.TaskId == 0 {
		return fmt.Errorf("task ID must be positive")
	}
	if strings.TrimSpace(r.Claimant) == "" {
		return fmt.Errorf("claimant cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(r.Claimant); err != nil {
		return fmt.Errorf("invalid claimant address: %s", err)
	}
	if r.Amount.IsZero() {
		return fmt.Errorf("reward amount cannot be zero")
	}
	if r.Amount.IsNegative() {
		return fmt.Errorf("reward amount cannot be negative")
	}
	if r.Timestamp <= 0 {
		return fmt.Errorf("reward timestamp must be positive")
	}

	return nil
}

func DefaultParams() Params {
	minBounty := sdk.NewCoin("stake", sdk.NewInt(1000)) // 1000 stake as minimum
	maxBounty := sdk.NewCoin("stake", sdk.NewInt(1000000)) // 1M stake as maximum

	return Params{
		MinBounty:             minBounty,
		MaxBounty:             maxBounty,
		MaxTitleLength:        100,
		MaxDescriptionLength:  1000,
		ProofTypes:            []string{"ipfs", "url", "text"},
		AutoApproveThreshold:  5,
		TaskExpiry:            86400 * 30, 
		ClaimDeadline:         86400 * 7,  
		SubmissionDeadline:    86400 * 14, 
	}
}

func CreateTaskReward(taskId uint64, claimant string, bounty sdk.Coin, txHash string, timestamp int64) TaskReward {
	return TaskReward{
		TaskId:    taskId,
		Claimant:  claimant,
		Amount:    bounty,
		Timestamp: timestamp,
		TxHash:    txHash,
	}
}

func ValidateRewardDistribution(task Task, reward TaskReward) error {
	if task.Status != TASK_STATUS_APPROVED {
		return fmt.Errorf("task must be approved to distribute rewards")
	}
	if task.Claimant != reward.Claimant {
		return fmt.Errorf("reward claimant %s does not match task claimant %s", reward.Claimant, task.Claimant)
	}
	if task.Id != reward.TaskId {
		return fmt.Errorf("reward task ID %d does not match task ID %d", reward.TaskId, task.Id)
	}
	if !task.Bounty.IsEqual(reward.Amount) {
		return fmt.Errorf("reward amount %s does not match task bounty %s", reward.Amount.String(), task.Bounty.String())
	}
	if err := reward.Validate(); err != nil {
		return fmt.Errorf("invalid reward: %s", err)
	}

	return nil
}

func CalculateRewardAmount(task Task, params Params, performanceScore float64) sdk.Coin {
	if performanceScore < 0.0 {
		performanceScore = 0.0
	} else if performanceScore > 1.0 {
		performanceScore = 1.0
	}
	baseReward := sdk.NewDecFromInt(task.Bounty.Amount).Mul(sdk.NewDec(performanceScore))
	rewardAmount := baseReward.TruncateInt()
	if rewardAmount.GT(task.Bounty.Amount) {
		rewardAmount = task.Bounty.Amount
	}
	if performanceScore >= 0.5 && rewardAmount.LT(params.MinBounty.Amount) {
		rewardAmount = params.MinBounty.Amount
	}

	return sdk.NewCoin(task.Bounty.Denom, rewardAmount)
}

func SplitTaskReward(reward TaskReward, recipients []string, weights []sdk.Int) []TaskReward {
	if len(recipients) != len(weights) {
		return nil
	}

	var rewards []TaskReward
	totalWeight := sdk.NewInt(0)

	// total weight
	for _, weight := range weights {
		totalWeight = totalWeight.Add(weight)
	}

	// rewards for each recipient
	for i, recipient := range recipients {
		if totalWeight.IsZero() {
			continue
		}
		portion := sdk.NewDecFromInt(weights[i]).Quo(sdk.NewDecFromInt(totalWeight))
		amount := sdk.NewDecFromInt(reward.Amount.Amount).Mul(portion).TruncateInt()
		if amount.IsZero() {
			continue
		}

		rewards = append(rewards, TaskReward{
			TaskId:    reward.TaskId,
			Claimant:  recipient,
			Amount:    sdk.NewCoin(reward.Amount.Denom, amount),
			Timestamp: reward.Timestamp,
			TxHash:    reward.TxHash,
		})
	}

	return rewards
}

func CheckAutoApproval(task Task, params Params, proof TaskProof, approvals uint32) bool {
	if params.AutoApproveThreshold == 0 {
		return false
	}
	if task.Status != TASK_STATUS_SUBMITTED {
		return false
	}
	if err := proof.Validate(params); err != nil {
		return false
	}
	return approvals >= params.AutoApproveThreshold
}

func EstimateTaskCompletionTime(task Task, params Params) time.Duration {
	bountyRatio := sdk.NewDecFromInt(task.Bounty.Amount).Quo(sdk.NewDecFromInt(params.MaxBounty.Amount))
	baseHours := 100.0
	estimatedHours := baseHours / bountyRatio.MustFloat64()
	if estimatedHours > 1000.0 {
		estimatedHours = 1000.0
	}
	
	return time.Duration(estimatedHours) * time.Hour
}

func GetTaskProgress(task Task) float64 {
	switch task.Status {
	case TASK_STATUS_UNDEFINED:
		return 0.0
	case TASK_STATUS_OPEN:
		return 0.0
	case TASK_STATUS_CLAIMED:
		return 0.25
	case TASK_STATUS_SUBMITTED:
		return 0.75
	case TASK_STATUS_APPROVED:
		return 1.0
	case TASK_STATUS_REJECTED:
		return 0.5
	case TASK_STATUS_CLOSED:
		if task.Claimant != "" {
			return 1.0
		}
		return 0.0
	default:
		return 0.0
	}
}