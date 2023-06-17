// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package aggregation_spotter

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// @notice execute approved operation
type ExecuteOperation struct {

	// [0] = [] testProgram
	// ··········· Callee program id
	//
	// [1] = [WRITE, SIGNER] executer
	//
	// [2] = [] operationAcc
	//
	// [3] = [WRITE] operationPda
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewExecuteOperationInstructionBuilder creates a new `ExecuteOperation` instruction builder.
func NewExecuteOperationInstructionBuilder() *ExecuteOperation {
	nd := &ExecuteOperation{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
	return nd
}

// SetTestProgramAccount sets the "testProgram" account.
// Callee program id
func (inst *ExecuteOperation) SetTestProgramAccount(testProgram ag_solanago.PublicKey) *ExecuteOperation {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(testProgram)
	return inst
}

// GetTestProgramAccount gets the "testProgram" account.
// Callee program id
func (inst *ExecuteOperation) GetTestProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetExecuterAccount sets the "executer" account.
func (inst *ExecuteOperation) SetExecuterAccount(executer ag_solanago.PublicKey) *ExecuteOperation {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(executer).WRITE().SIGNER()
	return inst
}

// GetExecuterAccount gets the "executer" account.
func (inst *ExecuteOperation) GetExecuterAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetOperationAccAccount sets the "operationAcc" account.
func (inst *ExecuteOperation) SetOperationAccAccount(operationAcc ag_solanago.PublicKey) *ExecuteOperation {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(operationAcc)
	return inst
}

// GetOperationAccAccount gets the "operationAcc" account.
func (inst *ExecuteOperation) GetOperationAccAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetOperationPdaAccount sets the "operationPda" account.
func (inst *ExecuteOperation) SetOperationPdaAccount(operationPda ag_solanago.PublicKey) *ExecuteOperation {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(operationPda).WRITE()
	return inst
}

// GetOperationPdaAccount gets the "operationPda" account.
func (inst *ExecuteOperation) GetOperationPdaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

func (inst ExecuteOperation) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_ExecuteOperation,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst ExecuteOperation) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *ExecuteOperation) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.TestProgram is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Executer is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.OperationAcc is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.OperationPda is not set")
		}
	}
	return nil
}

func (inst *ExecuteOperation) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("ExecuteOperation")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=4]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta(" testProgram", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("    executer", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("operationAcc", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("operationPda", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (obj ExecuteOperation) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *ExecuteOperation) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewExecuteOperationInstruction declares a new ExecuteOperation instruction with the provided parameters and accounts.
func NewExecuteOperationInstruction(
	// Accounts:
	testProgram ag_solanago.PublicKey,
	executer ag_solanago.PublicKey,
	operationAcc ag_solanago.PublicKey,
	operationPda ag_solanago.PublicKey) *ExecuteOperation {
	return NewExecuteOperationInstructionBuilder().
		SetTestProgramAccount(testProgram).
		SetExecuterAccount(executer).
		SetOperationAccAccount(operationAcc).
		SetOperationPdaAccount(operationPda)
}