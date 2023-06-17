// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package aggregation_spotter

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// @notice Reenable keeper in whitelist
// @param keeper address of keeper to remove
type EnableKeeper struct {
	Keeper *ag_solanago.PublicKey

	// [0] = [WRITE] spotter
	//
	// [1] = [SIGNER] admin
	//
	// [2] = [WRITE] keeperPda
	//
	// [3] = [] keeperAcc
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewEnableKeeperInstructionBuilder creates a new `EnableKeeper` instruction builder.
func NewEnableKeeperInstructionBuilder() *EnableKeeper {
	nd := &EnableKeeper{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
	return nd
}

// SetKeeper sets the "keeper" parameter.
func (inst *EnableKeeper) SetKeeper(keeper ag_solanago.PublicKey) *EnableKeeper {
	inst.Keeper = &keeper
	return inst
}

// SetSpotterAccount sets the "spotter" account.
func (inst *EnableKeeper) SetSpotterAccount(spotter ag_solanago.PublicKey) *EnableKeeper {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(spotter).WRITE()
	return inst
}

// GetSpotterAccount gets the "spotter" account.
func (inst *EnableKeeper) GetSpotterAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetAdminAccount sets the "admin" account.
func (inst *EnableKeeper) SetAdminAccount(admin ag_solanago.PublicKey) *EnableKeeper {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(admin).SIGNER()
	return inst
}

// GetAdminAccount gets the "admin" account.
func (inst *EnableKeeper) GetAdminAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetKeeperPdaAccount sets the "keeperPda" account.
func (inst *EnableKeeper) SetKeeperPdaAccount(keeperPda ag_solanago.PublicKey) *EnableKeeper {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(keeperPda).WRITE()
	return inst
}

// GetKeeperPdaAccount gets the "keeperPda" account.
func (inst *EnableKeeper) GetKeeperPdaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetKeeperAccAccount sets the "keeperAcc" account.
func (inst *EnableKeeper) SetKeeperAccAccount(keeperAcc ag_solanago.PublicKey) *EnableKeeper {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(keeperAcc)
	return inst
}

// GetKeeperAccAccount gets the "keeperAcc" account.
func (inst *EnableKeeper) GetKeeperAccAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

func (inst EnableKeeper) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_EnableKeeper,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst EnableKeeper) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *EnableKeeper) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Keeper == nil {
			return errors.New("Keeper parameter is not set")
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
			return errors.New("accounts.KeeperPda is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.KeeperAcc is not set")
		}
	}
	return nil
}

func (inst *EnableKeeper) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("EnableKeeper")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Keeper", *inst.Keeper))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=4]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("  spotter", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("    admin", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("keeperPda", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("keeperAcc", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (obj EnableKeeper) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Keeper` param:
	err = encoder.Encode(obj.Keeper)
	if err != nil {
		return err
	}
	return nil
}
func (obj *EnableKeeper) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Keeper`:
	err = decoder.Decode(&obj.Keeper)
	if err != nil {
		return err
	}
	return nil
}

// NewEnableKeeperInstruction declares a new EnableKeeper instruction with the provided parameters and accounts.
func NewEnableKeeperInstruction(
	// Parameters:
	keeper ag_solanago.PublicKey,
	// Accounts:
	spotter ag_solanago.PublicKey,
	admin ag_solanago.PublicKey,
	keeperPda ag_solanago.PublicKey,
	keeperAcc ag_solanago.PublicKey) *EnableKeeper {
	return NewEnableKeeperInstructionBuilder().
		SetKeeper(keeper).
		SetSpotterAccount(spotter).
		SetAdminAccount(admin).
		SetKeeperPdaAccount(keeperPda).
		SetKeeperAccAccount(keeperAcc)
}