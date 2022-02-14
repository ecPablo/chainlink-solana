// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package ocr_2

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// ProposePayees is the `proposePayees` instruction.
type ProposePayees struct {
	TokenMint *ag_solanago.PublicKey
	Payees    *[]ag_solanago.PublicKey

	// [0] = [WRITE] proposal
	//
	// [1] = [SIGNER] authority
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewProposePayeesInstructionBuilder creates a new `ProposePayees` instruction builder.
func NewProposePayeesInstructionBuilder() *ProposePayees {
	nd := &ProposePayees{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 2),
	}
	return nd
}

// SetTokenMint sets the "tokenMint" parameter.
func (inst *ProposePayees) SetTokenMint(tokenMint ag_solanago.PublicKey) *ProposePayees {
	inst.TokenMint = &tokenMint
	return inst
}

// SetPayees sets the "payees" parameter.
func (inst *ProposePayees) SetPayees(payees []ag_solanago.PublicKey) *ProposePayees {
	inst.Payees = &payees
	return inst
}

// SetProposalAccount sets the "proposal" account.
func (inst *ProposePayees) SetProposalAccount(proposal ag_solanago.PublicKey) *ProposePayees {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(proposal).WRITE()
	return inst
}

// GetProposalAccount gets the "proposal" account.
func (inst *ProposePayees) GetProposalAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// SetAuthorityAccount sets the "authority" account.
func (inst *ProposePayees) SetAuthorityAccount(authority ag_solanago.PublicKey) *ProposePayees {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(authority).SIGNER()
	return inst
}

// GetAuthorityAccount gets the "authority" account.
func (inst *ProposePayees) GetAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst ProposePayees) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_ProposePayees,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst ProposePayees) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *ProposePayees) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.TokenMint == nil {
			return errors.New("TokenMint parameter is not set")
		}
		if inst.Payees == nil {
			return errors.New("Payees parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Proposal is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Authority is not set")
		}
	}
	return nil
}

func (inst *ProposePayees) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("ProposePayees")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=2]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("TokenMint", *inst.TokenMint))
						paramsBranch.Child(ag_format.Param("   Payees", *inst.Payees))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=2]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta(" proposal", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("authority", inst.AccountMetaSlice[1]))
					})
				})
		})
}

func (obj ProposePayees) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `TokenMint` param:
	err = encoder.Encode(obj.TokenMint)
	if err != nil {
		return err
	}
	// Serialize `Payees` param:
	err = encoder.Encode(obj.Payees)
	if err != nil {
		return err
	}
	return nil
}
func (obj *ProposePayees) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `TokenMint`:
	err = decoder.Decode(&obj.TokenMint)
	if err != nil {
		return err
	}
	// Deserialize `Payees`:
	err = decoder.Decode(&obj.Payees)
	if err != nil {
		return err
	}
	return nil
}

// NewProposePayeesInstruction declares a new ProposePayees instruction with the provided parameters and accounts.
func NewProposePayeesInstruction(
	// Parameters:
	tokenMint ag_solanago.PublicKey,
	payees []ag_solanago.PublicKey,
	// Accounts:
	proposal ag_solanago.PublicKey,
	authority ag_solanago.PublicKey) *ProposePayees {
	return NewProposePayeesInstructionBuilder().
		SetTokenMint(tokenMint).
		SetPayees(payees).
		SetProposalAccount(proposal).
		SetAuthorityAccount(authority)
}