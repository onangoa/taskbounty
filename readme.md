# taskbounty
**taskbounty** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Additionally, Ignite CLI offers a frontend scaffolding feature (based on Vue) to help you quickly build a web frontend for your blockchain:

Use: `ignite scaffold vue`
This command can be run within your scaffolded blockchain project.


For more information see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/username/taskbounty@latest! | sudo bash
```
`username/taskbounty` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/ignite/installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.com/invite/ignitecli)


# TaskBounty — Blockchain-Based Task Bounty System

## Architecture & Design Decisions

### 1. High-Level Architecture
TaskBounty follows the standard Cosmos architecture:
- **Application Layer:** Custom task module  
- **SDK Core:** Handles transactions, accounts, and bank logic  
- **Tendermint Core:** Provides consensus and networking  
- **REST / gRPC / CLI:** For easy interaction

### 2. Key Design Decisions
1. **Task State Machine** – Ensures predictable task progression and prevents illegal transitions.  
2. **Token Economics** – Bounties are locked in the creator’s account until task completion.  
3. **Approval Logic** – Only task creators can approve, validated via message signers.  
4. **Proof Validation** – Supports on-chain text, URLs, or file hashes stored off-chain via IPFS.  
5. **Storage Design** – Uses collections and indexed maps for efficient queries and pagination.

### 3. Security Model
- Role-based permissions at the message level  
- Input validation and transition checks  
- Immutable records of proofs and rewards

---

## Smart Contract Logic

### Module Structure
```
x/task/
├── keeper/           # Core logic for state reads and writes
├── types/            # Data structures and validation
├── client/           # CLI and API interface
└── module/           # Registration and initialization
```
The **Keeper** manages blockchain state and validation.  
The **Types** folder defines entities like Task, TaskProof, and TaskReward.

### Core Data Structures
- **Task:** Title, description, bounty, and lifecycle state  
- **TaskProof:** Proof of completion by claimant  
- **TaskReward:** Record of bounty for completed tasks

### Keeper and Message Handlers
Handlers include:
- `CreateTask`
- `ClaimTask`
- `SubmitTask`
- `ApproveTask`
- `RejectTask`

Each validates inputs, checks permissions, and updates the chain atomically.

### Task Lifecycle
```
Open → Claimed → Submitted → Approved → Closed
```

### Rewards and Governance
Upon approval, a `TaskReward` record is created.  
Actual payments are handled by the **Cosmos bank module**, ensuring security and traceability.

---

## Local Development and Dockerized Deployment

### Local Development Setup

#### 1. Setup Local Testnet
```bash
go run ./cmd/taskbountyd/main.go init dev-node --chain-id taskbounty-dev

go run ./cmd/taskbountyd/main.go keys add creator
go run ./cmd/taskbountyd/main.go keys add claimant

go run ./cmd/taskbountyd/main.go add-genesis-account $(go run ./cmd/taskbountyd/main.go keys show creator -a) 1000000000stake,1000000000token
go run ./cmd/taskbountyd/main.go add-genesis-account $(go run ./cmd/taskbountyd/main.go keys show claimant -a) 1000000000stake,1000000000token

go run ./cmd/taskbountyd/main.go gentx creator 100000000stake --chain-id taskbounty-dev
go run ./cmd/taskbountyd/main.go collect-gentxs
go run ./cmd/taskbountyd/main.go start --api.enable
```

---

### Dockerized Deployment

#### 1. Clone Repository
```bash
git clone https://github.com/onangoa/taskbounty.git
cd taskbounty
```

#### 2. Build and Run
```bash
docker compose -f docker-compose.yml up -d --build
```

#### 3. Verify Setup
```bash
docker ps
curl http://localhost:26657/status
```

#### 4. Run Tests Inside Container
```bash
docker exec -it taskbounty-node-real sh
taskbountyd status
taskbountyd keys list --keyring-backend test
taskbountyd query task list
```

CLI walkthrough.

---

## Quick Demo Commands (Inside Container)
```bash
taskbountyd tx task create "Demo Task" "Finish this task" "1000stake"   --from validator   --chain-id taskbounty-dev   --keyring-backend test   --gas auto   --gas-adjustment 1.5   --gas-prices 0.025stake -y

taskbountyd tx task claim 0 --from claimant --chain-id taskbounty-dev --keyring-backend test --gas 100000 --gas-prices 0.025stake -y
taskbountyd tx task submit 0 "proof_hash_123" "url" "https://github.com/example/repo/pull/123" --from claimant --chain-id taskbounty-dev --keyring-backend test --gas 100000 --gas-prices 0.025stake -y
taskbountyd tx task approve 0 "approval_tx_hash_$(date +%s)" --from validator --chain-id taskbounty-dev --keyring-backend test --gas 100000 --gas-prices 0.025stake -y
```

---

## End


