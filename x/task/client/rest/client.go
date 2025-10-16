package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"

	taskTypes "taskbounty/x/task/types"
)

// TaskClient provides a REST client for the task module
type TaskClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewTaskClient creates a new TaskClient
func NewTaskClient(baseURL string) *TaskClient {
	return &TaskClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// GetParams fetches the module parameters
func (c *TaskClient) GetParams(ctx context.Context) (*taskTypes.QueryParamsResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/params", c.BaseURL)
	
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result taskTypes.QueryParamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTask fetches a task by ID
func (c *TaskClient) GetTask(ctx context.Context, id uint64) (*taskTypes.QueryGetTaskResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task/%d", c.BaseURL, id)
	
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result taskTypes.QueryGetTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListTasks fetches all tasks
func (c *TaskClient) ListTasks(ctx context.Context) (*taskTypes.QueryAllTaskResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task", c.BaseURL)
	
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result taskTypes.QueryAllTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTaskReward fetches a task reward by ID
func (c *TaskClient) GetTaskReward(ctx context.Context, id uint64) (*taskTypes.QueryGetTaskRewardResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task_reward/%d", c.BaseURL, id)
	
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result taskTypes.QueryGetTaskRewardResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListTaskRewards fetches all task rewards
func (c *TaskClient) ListTaskRewards(ctx context.Context) (*taskTypes.QueryAllTaskRewardResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task_reward", c.BaseURL)
	
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result taskTypes.QueryAllTaskRewardResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTaskRewardsByClaimant fetches all task rewards for a given claimant
func (c *TaskClient) GetTaskRewardsByClaimant(ctx context.Context, claimant string) (*taskTypes.QueryGetTaskRewardsByClaimantResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task_rewards/%s", c.BaseURL, url.PathEscape(claimant))
	
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result taskTypes.QueryGetTaskRewardsByClaimantResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateTask creates a new task
func (c *TaskClient) CreateTask(ctx context.Context, req *taskTypes.MsgCreateTask) (*taskTypes.MsgCreateTaskResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task", c.BaseURL)
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(endpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result taskTypes.MsgCreateTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateTask updates an existing task
func (c *TaskClient) UpdateTask(ctx context.Context, req *taskTypes.MsgUpdateTask) (*taskTypes.MsgUpdateTaskResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task/%d", c.BaseURL, req.Id)
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result taskTypes.MsgUpdateTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteTask deletes a task
func (c *TaskClient) DeleteTask(ctx context.Context, req *taskTypes.MsgDeleteTask) (*taskTypes.MsgDeleteTaskResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task/%d", c.BaseURL, req.Id)
	
	httpReq, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result taskTypes.MsgDeleteTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ClaimTask claims a task
func (c *TaskClient) ClaimTask(ctx context.Context, req *taskTypes.MsgClaimTask) (*taskTypes.MsgClaimTaskResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task/%d/claim", c.BaseURL, req.Id)
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(endpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result taskTypes.MsgClaimTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SubmitTask submits a completed task
func (c *TaskClient) SubmitTask(ctx context.Context, req *taskTypes.MsgSubmitTask) (*taskTypes.MsgSubmitTaskResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task/%d/submit", c.BaseURL, req.Id)
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(endpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result taskTypes.MsgSubmitTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ApproveTask approves a submitted task
func (c *TaskClient) ApproveTask(ctx context.Context, req *taskTypes.MsgApproveTask) (*taskTypes.MsgApproveTaskResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task/%d/approve", c.BaseURL, req.Id)
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(endpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result taskTypes.MsgApproveTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RejectTask rejects a submitted task
func (c *TaskClient) RejectTask(ctx context.Context, req *taskTypes.MsgRejectTask) (*taskTypes.MsgRejectTaskResponse, error) {
	endpoint := fmt.Sprintf("%s/taskbounty/task/v1/task/%d/reject", c.BaseURL, req.Id)
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(endpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result taskTypes.MsgRejectTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}