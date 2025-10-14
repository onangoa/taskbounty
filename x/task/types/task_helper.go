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

