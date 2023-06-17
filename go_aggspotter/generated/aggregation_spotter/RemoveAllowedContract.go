// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package aggregation_spotter

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// @notice Removing contract from whitelist
// @param _contract address of contract to remove
type RemoveAllowedContract struct {
	Contract *ag_solanago.PublicKey

	// [0] = [WRITE] spotter
	//
	// [1] = [SIGNER] admin
	//
	// [2] = [WRITE] contractPda
	//
	// [3] = [] contractAcc
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewRemoveAllowedContractInstructionBuilder creates a new `RemoveAllowedContract` instruction builder.
func NewRemoveAllowedContractInstructionBuilder() *RemoveAllowedContract {
	nd := &RemoveAllowedContract{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
	return nd
}

// SetContract sets the "contract" parameter.
func (inst *RemoveAllowedContract) SetContract(contract ag_solanago.PublicKey) *RemoveAllowedContract {
	inst.Contract = &contract
	return inst
}

// SetSpotterAccount sets the "spotter" account.
func (inst *RemoveAllowedContract) SetSpotterAccount(spotter ag_solanago.PublicKey) *RemoveAllowedContract {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(spotter).WRITE()
	return inst
}

// GetSpotterAccount gets the "spotter" account.
func (inst *RemoveAllowedContract) GetSpotterAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetAdminAccount sets the "admin" account.
func (inst *RemoveAllowedContract) SetAdminAccount(admin ag_solanago.PublicKey) *RemoveAllowedContract {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(admin).SIGNER()
	return inst
}

// GetAdminAccount gets the "admin" account.
func (inst *RemoveAllowedContract) GetAdminAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetContractPdaAccount sets the "contractPda" account.
func (inst *RemoveAllowedContract) SetContractPdaAccount(contractPda ag_solanago.PublicKey) *RemoveAllowedContract {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(contractPda).WRITE()
	return inst
}

// GetContractPdaAccount gets the "contractPda" account.
func (inst *RemoveAllowedContract) GetContractPdaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetContractAccAccount sets the "contractAcc" account.
func (inst *RemoveAllowedContract) SetContractAccAccount(contractAcc ag_solanago.PublicKey) *RemoveAllowedContract {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(contractAcc)
	return inst
}

// GetContractAccAccount gets the "contractAcc" account.
func (inst *RemoveAllowedContract) GetContractAccAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

func (inst RemoveAllowedContract) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_RemoveAllowedContract,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst RemoveAllowedContract) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *RemoveAllowedContract) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Contract == nil {
			return errors.New("Contract parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Spotter is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Admin is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.ContractPda is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.ContractAcc is not set")
		}
	}
	return nil
}

func (inst *RemoveAllowedContract) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("RemoveAllowedContract")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Contract", *inst.Contract))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=4]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("    spotter", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("      admin", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("contractPda", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("contractAcc", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (obj RemoveAllowedContract) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Contract` param:
	err = encoder.Encode(obj.Contract)
	if err != nil {
		return err
	}
	return nil
}
func (obj *RemoveAllowedContract) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Contract`:
	err = decoder.Decode(&obj.Contract)
	if err != nil {
		return err
	}
	return nil
}

// NewRemoveAllowedContractInstruction declares a new RemoveAllowedContract instruction with the provided parameters and accounts.
func NewRemoveAllowedContractInstruction(
	// Parameters:
	contract ag_solanago.PublicKey,
	// Accounts:
	spotter ag_solanago.PublicKey,
	admin ag_solanago.PublicKey,
	contractPda ag_solanago.PublicKey,
	contractAcc ag_solanago.PublicKey) *RemoveAllowedContract {
	return NewRemoveAllowedContractInstructionBuilder().
		SetContract(contract).
		SetSpotterAccount(spotter).
		SetAdminAccount(admin).
		SetContractPdaAccount(contractPda).
		SetContractAccAccount(contractAcc)
}