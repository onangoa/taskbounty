package task

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"taskbounty/testutil/sample"
	tasksimulation "taskbounty/x/task/simulation"
	"taskbounty/x/task/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	taskGenesis := types.GenesisState{
		Params:   types.DefaultParams(),
		TaskList: []types.Task{{Id: 0, Creator: sample.AccAddress()}, {Id: 1, Creator: sample.AccAddress()}}, TaskCount: 2,
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&taskGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateTask          = "op_weight_msg_task"
		defaultWeightMsgCreateTask int = 100
	)

	var weightMsgCreateTask int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateTask, &weightMsgCreateTask, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTask = defaultWeightMsgCreateTask
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTask,
		tasksimulation.SimulateMsgCreateTask(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateTask          = "op_weight_msg_task"
		defaultWeightMsgUpdateTask int = 100
	)

	var weightMsgUpdateTask int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateTask, &weightMsgUpdateTask, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateTask = defaultWeightMsgUpdateTask
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateTask,
		tasksimulation.SimulateMsgUpdateTask(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteTask          = "op_weight_msg_task"
		defaultWeightMsgDeleteTask int = 100
	)

	var weightMsgDeleteTask int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteTask, &weightMsgDeleteTask, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteTask = defaultWeightMsgDeleteTask
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteTask,
		tasksimulation.SimulateMsgDeleteTask(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
