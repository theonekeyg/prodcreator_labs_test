package spotter

import (
	"context"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/gagliardetto/solana-go"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"go_aggspotter/generated/aggregation_spotter"
)

// AggregationSpotter client is a wrapper around solana-go-rpc.Client.
// It adds some convenience methods for AggregationSpotter API.
type AggregationSpotter struct {
	rpcClient *rpc.Client
	wsClient  *ws.Client
	pk        solana.PrivateKey
}

// NewAggregationSpotterClient creates a new AggregationSpotter client.
func NewSpotterWithEndpoints(rpcEndpoint string, wsEndpoint string, pkFile string) *AggregationSpotter {
	rpcClient  := rpc.New(rpcEndpoint)
	pk, err := solana.PrivateKeyFromSolanaKeygenFile(pkFile)

	if err != nil {
		panic(err)
	}

	wsClient, err := ws.Connect(context.Background(), wsEndpoint)

	return &AggregationSpotter {
		rpcClient: rpcClient,
		wsClient: wsClient,
		pk: pk,
	}
}

func NewSpotter(pkFile string) *AggregationSpotter {
	// return NewWithEndpoints(rpc.DevNet_RPC, rpc.DevNet_WS, pkFile)
	return NewSpotterWithEndpoints(rpc.LocalNet_RPC, rpc.LocalNet_WS, pkFile)
}

// Sends and confirms a transaction.
func (s *AggregationSpotter) sendTransaction(tx *solana.Transaction) solana.Signature {
	sig, err := confirm.SendAndConfirmTransaction(
		context.TODO(),
		s.rpcClient,
		s.wsClient,
		tx,
	)

	if err != nil {
		panic(err)
	}

	return sig
}

// Sends a transaction and returns transaction logs.
func (s *AggregationSpotter) sendTransactionWithLogs(tx *solana.Transaction) []string {
	sig := s.sendTransaction(tx)
	log.Debug("sig:", sig)

	chainTx, err := s.rpcClient.GetTransaction(context.TODO(), sig, nil)
	if err != nil {
		panic(err)
	}

	logs := chainTx.Meta.LogMessages
	return logs
}

// Sends and confirms a transaction. After the transaction is confirmed, check logs for
// ProposalApproved event. If found, return true, else return false.
func (s *AggregationSpotter) sendAndCheckForAllowace(tx *solana.Transaction) bool {
	logs := s.sendTransactionWithLogs(tx)
	log.Debug("logs:", logs)

	// Search for "ProposalApproved(**Solana Public Key**)" regex in each string in logs
	// If found, return true, else return false

	regex := regexp.MustCompile(`ProposalApproved\((\w+)\)`)
	for _, log := range logs {
		if regex.MatchString(log) {
			return true
		}
	}

	return false

}

// Execute the needed script to create operation (from the frist account in keepersAccs array,
// vote for this operaiton with other keepers, and if operation is approved, execute it.
func (s *AggregationSpotter) ExecuteScript(keeperAccs []solana.PrivateKey, operationAcc solana.PublicKey) {

	if len(keeperAccs) < 2 {
		panic("keeperAccs must be at least 2")
	}

	pa := solana.MustPublicKeyFromBase58(PA_DEFAULT)
	testProgram := solana.MustPublicKeyFromBase58(TEST_PROGRAM)
	spotter, _ := getSpotterPda(pa)
	keeper1Acc := keeperAccs[0];
	keeper1, _ := getKeeperPda(keeper1Acc.PublicKey(), pa)
	operation, _ := getOperationPda(operationAcc, pa)
	systemProgram := solana.MustPublicKeyFromBase58(SYSTEM_PROGRAM)

	// Create CreateOperation instruction
	operationData := aggregation_spotter.OperationData {
		Contr: testProgram,
		Accounts: []solana.PublicKey {
			testProgram, s.pk.PublicKey(),
		},
	}

	recent, err := s.rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)

	if err != nil {
		panic(err)
	}

	log.Debug("recent:", recent)

	// Build instruction
	inst, err := aggregation_spotter.NewCreateOperationInstructionBuilder().
		SetOperationData(operationData).
		SetGasPrice(50000).
		SetSpotterAccount(spotter).
		SetKeeperAccAccount(keeper1Acc.PublicKey()).
		SetKeeperPdaAccount(keeper1).
		SetOperationAccAccount(operationAcc).
		SetOperationPdaAccount(operation).
		SetSystemProgramAccount(systemProgram).
		ValidateAndBuild()

	// Print accounts from instruction above
	fmt.Printf("Built CreateOperation Instruction: inst.spotter: %s; inst.keeperAcc: %s; inst.keeperPda: %s; inst.operationAcc: %s; inst.operationPda: %s; inst.systemProgram: %s\n",
		spotter, keeper1Acc.PublicKey(), keeper1, operationAcc, operation, systemProgram,
	)

	if err != nil {
		panic(err)
	}

	log.Debug("inst:", inst)

	// Wrap instruction into transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ inst },
		recent.Value.Blockhash,
		solana.TransactionPayer(keeper1Acc.PublicKey()),
	)

	if err != nil {
		panic(err)
	}

	log.Debug("tx:", tx)

	// Sign transaction
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if keeper1Acc.PublicKey().Equals(key) {
				return &keeper1Acc
			}
			return nil
		},
	)

	if err != nil {
		panic(err)
	}

	// Send and confirm transaction
	done := s.sendAndCheckForAllowace(tx)
	log.Debug("done:", done)


	if done == false {
		for _, keeperAcc := range keeperAccs[1:] {

			keeper, _ := getKeeperPda(keeperAcc.PublicKey(), pa)
			// vote for operation
			inst, err := aggregation_spotter.NewProposeOperationInstructionBuilder().
				SetGasPrice(50000).
				SetSpotterAccount(spotter).
				SetKeeperAccAccount(keeperAcc.PublicKey()).
				SetKeeperPdaAccount(keeper).
				SetOperationAccAccount(operationAcc).
				SetOperationPdaAccount(operation).
				ValidateAndBuild()
			fmt.Printf("inst.spotter: %s; inst.keeperAcc: %s; inst.keeperPda: %s; inst.operationAcc: %s; inst.operationPda: %s\n",
				spotter, keeperAcc.PublicKey(), keeper, operationAcc, operation,
			)

			if err != nil {
				panic(err)
			}
			log.Debug("inst:", inst)

			tx, err := solana.NewTransaction(
				[]solana.Instruction{ inst },
				recent.Value.Blockhash,
				solana.TransactionPayer(keeperAcc.PublicKey()),
			)
			log.Debug("tx:", tx)

			// Sign transaction
			_, err = tx.Sign(
				func(key solana.PublicKey) *solana.PrivateKey {
					if keeperAcc.PublicKey().Equals(key) {
						return &keeperAcc
					}
					return nil
				},
			)

			if err != nil {
				panic(err)
			}

			done = s.sendAndCheckForAllowace(tx)
			log.Debug("done:", done)

			if done { break }
		}
	} // !done


	// if done, execute operation
	if done {

		log.Info("Operation is approved, executing it...")

		inst, err := aggregation_spotter.NewExecuteOperationInstructionBuilder().
			SetTestProgramAccount(testProgram).
			SetExecuterAccount(s.pk.PublicKey()).
			SetOperationAccAccount(operationAcc).
			SetOperationPdaAccount(operation).
			ValidateAndBuild()

		fmt.Printf("inst.testProgram: %s; inst.executer: %s; inst.operationAcc: %s; inst.operationPda: %s\n",
			testProgram, s.pk.PublicKey(), operationAcc, operation,
		)

		if err != nil {
			panic(err)
		}

		log.Debug("inst:", inst)

		tx, err := solana.NewTransaction(
			[]solana.Instruction{ inst },
			recent.Value.Blockhash,
			solana.TransactionPayer(s.pk.PublicKey()),
		)

		if err != nil {
			panic(err)
		}

		log.Debug("tx:", tx)

		// Sign transaction

		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if s.pk.PublicKey().Equals(key) {
					return &s.pk
				}
				return nil
			},
		)

		if err != nil {
			panic(err)
		}

		logs := s.sendTransactionWithLogs(tx)
		log.Debug("logs:", logs)
	}
}

const PA_DEFAULT = "3XQdG1Zpk151xuGHSd6DUkNuh9m9i3M8ptxJEHNfZdJ2";
const SYSTEM_PROGRAM = "11111111111111111111111111111111";
const TEST_PROGRAM = "FBXArrxAnAZpmGmVYD8jQVauohCrXHaJBgTXzUPKRyHa";
const SPOTTER_SEED = "spotter";
const KEEPER_SEED = "keeper";
const CONTRACT_SEED = "contract";
const OPERATION_SEED = "operation";

func getSpotterPda(programId solana.PublicKey) (solana.PublicKey, byte) {
	// func FindProgramAddress(seed [][]byte, programID PublicKey) (PublicKey, uint8, error)
	spotter, bump, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(SPOTTER_SEED),
		},
		programId,
	)

	if err != nil {
		panic(err)
	}

	return spotter, bump
}

func getKeeperPda(keeper solana.PublicKey, programId solana.PublicKey) (solana.PublicKey, byte) {
	// func FindProgramAddress(seed [][]byte, programID PublicKey) (PublicKey, uint8, error)
	keeper, bump, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(KEEPER_SEED),
			keeper[:],
		},
		programId,
	)

	if err != nil {
		panic(err)
	}

	return keeper, bump
}

func getContractPda(contract solana.PublicKey, programId solana.PublicKey) (solana.PublicKey, byte) {
	// func FindProgramAddress(seed [][]byte, programID PublicKey) (PublicKey, uint8, error)
	contract, bump, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(CONTRACT_SEED),
			contract[:],
		},
		programId,
	)

	if err != nil {
		panic(err)
	}

	return contract, bump
}

func getOperationPda(operation solana.PublicKey, programId solana.PublicKey) (solana.PublicKey, byte) {
	// func FindProgramAddress(seed [][]byte, programID PublicKey) (PublicKey, uint8, error)
	operation, bump, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(OPERATION_SEED),
			operation[:],
		},
		programId,
	)

	if err != nil {
		panic(err)
	}

	return operation, bump
}
