## 1. Setup Local Testnet
First, set up your local development environment:
```bash
# Initialize the blockchain
go run ./cmd/taskbountyd/main.go init dev-node --chain-id taskbounty-dev
# Create two accounts: creator and claimant
go run ./cmd/taskbountyd/main.go keys add creator
go run ./cmd/taskbountyd/main.go keys add claimant
# Add genesis accounts with initial tokens
go run ./cmd/taskbountyd/main.go add-genesis-account $(go run
./cmd/taskbountyd/main.go keys show creator -a) 1000000000stake,1000000000token
go run ./cmd/taskbountyd/main.go add-genesis-account $(go run
./cmd/taskbountyd/main.go keys show claimant -a) 1000000000stake,1000000000token
# Create genesis transaction for the validator
go run ./cmd/taskbountyd/main.go gentx creator 100000000stake --chain-id
taskbounty-dev
# Collect genesis transactions
go run ./cmd/taskbountyd/main.go collect-gentxs
# Start the blockchain with REST API enabled
go run ./cmd/taskbountyd/main.go start --api.enable
```
# 1. Get creator address
CREATOR_ADDR=$(go run ./cmd/taskbountyd/main.go keys show creator -a)

# 2. Query creator balances
go run ./cmd/taskbountyd/main.go query bank balances $CREATOR_ADDR

# 1. Get the claimant address
CLAIMANT_ADDR=$(go run ./cmd/taskbountyd/main.go keys show claimant -a)

# 2. Query claimant's bank balance
go run ./cmd/taskbountyd/main.go query bank balances $CLAIMANT_ADDR


## 2. Complete Task Test
### Create a Task (as Creator)
```bash
# Terminal 1: Create a task with bounty
go run ./cmd/taskbountyd/main.go tx task create "Test Task" "Complete this task
for bounty" "1000stake" \
--from creator \
--chain-id taskbounty-dev \
--gas auto \
--gas-adjustment 1.5 \
--gas-prices 0.025stake \
--yes
```
### Step 2: Verify Task Creation
```bash
# Terminal 1: Query the task
# Query the task
go run ./cmd/taskbountyd/main.go query task get 0
```
### Step 3: Claim the Task (as Claimant)
```bash
#  Claim the task
go run ./cmd/taskbountyd/main.go tx task claim 3 \
--from claimant \
--chain-id taskbounty-dev \
--gas auto \
--gas-adjustment 1.5 \
--gas-prices 0.025stake \
--yes
```
### Step 4: Submit the Task (as Claimant)
```bash
#  Submit proof of work

go run ./cmd/taskbountyd/main.go tx task submit 3 "$(echo -n 'https://github.com/example/repo/pull/123' | sha256sum | cut -d' ' -f1)" "text" "https://github.com/example/repo/pull/123" \
--from claimant \
--chain-id taskbounty-dev \
--gas auto \
--gas-adjustment 1.5 \
--gas-prices 0.025stake \
--yes

```
### Step 5: Approve the Task (as Creator) - This triggers the bounty transfer
```bash
# Terminal 1: Approve and release bounty
98E1F9566A5F55702E2B1FB06BD1582AF0F4135C4FB0F8F45B651B2E570CE8F0
"TX_HASH_PLACEHOLDER"
go run ./cmd/taskbountyd/main.go tx task approve 3 98E1F9566A5F55702E2B1FB06BD1582AF0F4135C4FB0F8F45B651B2E570CE8F0 \
--from creator \
--chain-id taskbounty-dev \
--gas auto \
--gas-adjustment 1.5 \
--gas-prices 0.025stake \
--yes

```
## 3. Verify the On-Chain Transfer
### Check Task Status
```bash
# Verify task is approved
go run ./cmd/taskbountyd/main.go query task get 0
```
### Check Reward Distribution
```bash
# Check the reward record
go run ./cmd/taskbountyd/main.go query task reward 0
# Check claimant's balance after approval
```
### Check Creator's Balance


#Docker
docker exec -it taskbounty-node-real sh – use (docker ps)
# Check blockchain status
taskbountyd status
# List available keys
taskbountyd keys list --keyring-backend test
# Get validator address
VALIDATOR_ADDRESS=$(taskbountyd keys show validator -a --keyring-backend test)
CLAIMANT_ADDRESS=$(taskbountyd keys show claimant -a --keyring-backend test)
taskbountyd query bank balances $CLAIMANT_ADDRESS
# Check validator balance
taskbountyd query bank balances $VALIDATOR_ADDRESS
#fund claimant
taskbountyd tx bank send validator $CLAIMANT_ADDRESS 1000000stake --from
validator --chain-id taskbounty-dev --keyring-backend test --gas 100000 --gas-
prices 0.025stake -y
# Create a test task
echo "📝 Creating a test task..."
taskbountyd tx task create "Manual Test Task" "This is a manual test task" "1000stake" --from validator --chain-id taskbounty-dev --keyring-backend test --gas auto --gas-adjustment 1.5 --gas-prices 0.025stake -y

# claim a task
taskbountyd tx task claim 3 \
--from $CLAIMANT_ADDRESS \
--chain-id taskbounty-dev \
--keyring-backend test \
--gas 100000 \
--gas-prices 0.025stake \
-y

```

taskbountyd query task list

#submit proof of work
taskbountyd tx task submit 3 "proof_hash_123" "url" "https://github.com/example/repo/pull/123" --from $CLAIMANT_ADDRESS --chain-id taskbounty-dev --keyring-backend test --gas 100000 --gas-prices 0.025stake -y

#stats
taskbountyd query task get 3

#approve
taskbountyd tx task approve 3 "approval_tx_hash_$(date +%s)" --from $VALIDATOR_ADDRESS --chain-id taskbounty-dev --keyring-backend test --gas 100000 --gas-prices 0.025stake -y


#reward record and verify the balance changes:
taskbountyd query task reward 3


taskbountyd query bank balances $CLAIMANT_ADDRESS
taskbountyd query bank balances $VALIDATOR_ADDRESS