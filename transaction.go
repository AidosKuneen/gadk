package gadk

import (
	"encoding/json"
	"errors"
	"time"
)

//Transaction is transaction  structure for adk.
type Transaction struct {
	SignatureMessageFragment Trytes
	Address                  Address
	Value                    int64 `json:",string"`
	Tag                      Trytes
	Timestamp                time.Time
	CurrentIndex             int64 `json:",string"`
	LastIndex                int64 `json:",string"`
	Bundle                   Trytes
	TrunkTransaction         Trytes
	BranchTransaction        Trytes
	Nonce                    Trytes
}

//errors for tx.
var (
	ErrInvalidTransactionType = errors.New("invalid transaction type")
	ErrInvalidTransactionHash = errors.New("invalid transaction hash")
	ErrInvalidTransaction     = errors.New("malformed transaction")
)

// Trinary sizes and offsets of a transaction
const (
	SignatureMessageFragmentTrinaryOffset = 0
	SignatureMessageFragmentTrinarySize   = 6561
	AddressTrinaryOffset                  = SignatureMessageFragmentTrinaryOffset + SignatureMessageFragmentTrinarySize
	AddressTrinarySize                    = 243
	ValueTrinaryOffset                    = AddressTrinaryOffset + AddressTrinarySize
	ValueTrinarySize                      = 81
	TagTrinaryOffset                      = ValueTrinaryOffset + ValueTrinarySize
	TagTrinarySize                        = 81
	TimestampTrinaryOffset                = TagTrinaryOffset + TagTrinarySize
	TimestampTrinarySize                  = 27
	CurrentIndexTrinaryOffset             = TimestampTrinaryOffset + TimestampTrinarySize
	CurrentIndexTrinarySize               = 27
	LastIndexTrinaryOffset                = CurrentIndexTrinaryOffset + CurrentIndexTrinarySize
	LastIndexTrinarySize                  = 27
	BundleTrinaryOffset                   = LastIndexTrinaryOffset + LastIndexTrinarySize
	BundleTrinarySize                     = 243
	TrunkTransactionTrinaryOffset         = BundleTrinaryOffset + BundleTrinarySize
	TrunkTransactionTrinarySize           = 243
	BranchTransactionTrinaryOffset        = TrunkTransactionTrinaryOffset + TrunkTransactionTrinarySize
	BranchTransactionTrinarySize          = 243
	NonceTrinaryOffset                    = BranchTransactionTrinaryOffset + BranchTransactionTrinarySize
	NonceTrinarySize                      = 243

	transactionTrinarySize = SignatureMessageFragmentTrinarySize + AddressTrinarySize +
		ValueTrinarySize + TagTrinarySize + TimestampTrinarySize +
		CurrentIndexTrinarySize + LastIndexTrinarySize + BundleTrinarySize +
		TrunkTransactionTrinarySize + BranchTransactionTrinarySize +
		NonceTrinarySize
)

//NewTransaction makes tx from trits.
func NewTransaction(trytes Trytes) (*Transaction, error) {
	t := Transaction{}
	if err := checkTx(trytes); err != nil {
		return nil, err
	}
	var err error
	t.SignatureMessageFragment = trytes[SignatureMessageFragmentTrinaryOffset/3 : SignatureMessageFragmentTrinarySize/3]
	t.Address, err = trytes[AddressTrinaryOffset/3 : (AddressTrinaryOffset+AddressTrinarySize)/3].ToAddress()
	if err != nil {
		return nil, err
	}
	t.Value = trytes[ValueTrinaryOffset/3 : (ValueTrinaryOffset+ValueTrinarySize)/3].Trits().Int()
	t.Tag = trytes[TagTrinaryOffset/3 : (TagTrinaryOffset+TagTrinarySize)/3]
	timestamp := trytes[TimestampTrinaryOffset/3 : (TimestampTrinaryOffset+TimestampTrinarySize)/3].Trits().Int()
	t.Timestamp = time.Unix(timestamp, 0)
	t.CurrentIndex = trytes[CurrentIndexTrinaryOffset/3 : (CurrentIndexTrinaryOffset+CurrentIndexTrinarySize)/3].Trits().Int()
	t.LastIndex = trytes[LastIndexTrinaryOffset/3 : (LastIndexTrinaryOffset+LastIndexTrinarySize)/3].Trits().Int()
	t.Bundle = trytes[BundleTrinaryOffset/3 : (BundleTrinaryOffset+BundleTrinarySize)/3]
	t.TrunkTransaction = trytes[TrunkTransactionTrinaryOffset/3 : (TrunkTransactionTrinaryOffset+TrunkTransactionTrinarySize)/3]
	t.BranchTransaction = trytes[BranchTransactionTrinaryOffset/3 : (BranchTransactionTrinaryOffset+BranchTransactionTrinarySize)/3]
	t.Nonce = trytes[NonceTrinaryOffset/3 : (NonceTrinaryOffset+NonceTrinarySize)/3]
	return &t, nil
}

func checkTx(trytes Trytes) error {
	if err := trytes.IsValid(); err != nil {
		return errors.New("invalid transaction " + err.Error())
	}
	if len(trytes) != transactionTrinarySize/3 {
		return errors.New("invalid trits counts in transaction")
	}

	if trytes[2279:2295] != "9999999999999999" {
		return errors.New("invalid value in transaction")
	}
	return nil
}

func (t *Transaction) parser(trits Trits) error {
	var err error
	t.SignatureMessageFragment = trits[SignatureMessageFragmentTrinaryOffset:SignatureMessageFragmentTrinarySize].Trytes()
	t.Address, err = trits[AddressTrinaryOffset : AddressTrinaryOffset+AddressTrinarySize].Trytes().ToAddress()
	if err != nil {
		return err
	}
	t.Value = trits[ValueTrinaryOffset : ValueTrinaryOffset+ValueTrinarySize].Int()
	t.Tag = trits[TagTrinaryOffset : TagTrinaryOffset+TagTrinarySize].Trytes()
	timestamp := trits[TimestampTrinaryOffset : TimestampTrinaryOffset+TimestampTrinarySize].Int()
	t.Timestamp = time.Unix(timestamp, 0)
	t.CurrentIndex = trits[CurrentIndexTrinaryOffset : CurrentIndexTrinaryOffset+CurrentIndexTrinarySize].Int()
	t.LastIndex = trits[LastIndexTrinaryOffset : LastIndexTrinaryOffset+LastIndexTrinarySize].Int()
	t.Bundle = trits[BundleTrinaryOffset : BundleTrinaryOffset+BundleTrinarySize].Trytes()
	t.TrunkTransaction = trits[TrunkTransactionTrinaryOffset : TrunkTransactionTrinaryOffset+TrunkTransactionTrinarySize].Trytes()
	t.BranchTransaction = trits[BranchTransactionTrinaryOffset : BranchTransactionTrinaryOffset+BranchTransactionTrinarySize].Trytes()
	t.Nonce = trits[NonceTrinaryOffset : NonceTrinaryOffset+NonceTrinarySize].Trytes()

	return nil
}

//Trytes converts the transaction to Trytes.
func (t *Transaction) Trytes() Trytes {
	tr := make(Trits, transactionTrinarySize)
	copy(tr, t.SignatureMessageFragment.Trits())
	copy(tr[AddressTrinaryOffset:], Trytes(t.Address).Trits())
	copy(tr[ValueTrinaryOffset:], Int2Trits(t.Value, ValueTrinarySize))
	copy(tr[TagTrinaryOffset:], t.Tag.Trits())
	copy(tr[TimestampTrinaryOffset:], Int2Trits(t.Timestamp.Unix(), TimestampTrinarySize))
	copy(tr[CurrentIndexTrinaryOffset:], Int2Trits(t.CurrentIndex, CurrentIndexTrinarySize))
	copy(tr[LastIndexTrinaryOffset:], Int2Trits(t.LastIndex, LastIndexTrinarySize))
	copy(tr[BundleTrinaryOffset:], t.Bundle.Trits())
	copy(tr[TrunkTransactionTrinaryOffset:], t.TrunkTransaction.Trits())
	copy(tr[BranchTransactionTrinaryOffset:], t.BranchTransaction.Trits())
	copy(tr[NonceTrinaryOffset:], t.Nonce.Trits())
	return tr.Trytes()
}

//HasValidNonce checks t's hash has valid MinWeightMagnitude.
func (t *Transaction) HasValidNonce() bool {
	h := t.Hash()
	for i := len(h) - 1; i > len(h)-1-MinWeightMagnitude/3; i-- {
		if h[i] != '9' {
			return false
		}
	}
	return true
}

//HasValidNonceMWM checks t's hash has valid MinWeightMagnitude.
func (t *Transaction) HasValidNonceMWM(mwm int) bool {
	h := t.Hash()
	for i := len(h) - 1; i > len(h)-1-mwm/3; i-- {
		if h[i] != '9' {
			return false
		}
	}
	return true
}

//Hash returns the hash of the transaction.
func (t *Transaction) Hash() Trytes {
	return t.Trytes().Hash()
}

//UnmarshalJSON makes transaction struct from json.
func (t *Transaction) UnmarshalJSON(b []byte) error {
	var s Trytes
	var err error
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}
	if err = t.parser(s.Trits()); err != nil {
		return err
	}
	return nil
}

//MarshalJSON makes trytes ([]byte) from a transaction.
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Trytes() + `"`), nil
}
