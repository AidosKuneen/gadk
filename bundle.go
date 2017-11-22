package gadk

import (
	"errors"
	"fmt"
	"time"
)

func pad(orig Trytes, size int) Trytes {
	out := make([]byte, size)
	copy(out, []byte(orig))
	for i := len(orig); i < size; i++ {
		out[i] = '9'
	}
	return Trytes(out)
}

//Bundle is transactions are bundled
// (or grouped) together during the creation of a transfer.
type Bundle []Transaction

//Add adds one bundle to bundle slice tempolary.
//For now elements which are not specified are filled with trits 0.
func (bs *Bundle) Add(num int, address Address, value int64, timestamp time.Time, tag Trytes) {
	if tag == "" {
		tag = EmptyHash
	}
	for i := 0; i < num; i++ {
		var v int64
		if i == 0 {
			v = value
		}
		b := Transaction{
			SignatureMessageFragment: emptySig,
			Address:                  address,
			Value:                    v,
			Tag:                      pad(tag, TagTrinarySize/3),
			Timestamp:                timestamp,
			CurrentIndex:             int64(len(*bs) - 1),
			LastIndex:                0,
			Bundle:                   EmptyHash,
			TrunkTransaction:         EmptyHash,
			BranchTransaction:        EmptyHash,
			Nonce:                    EmptyHash,
		}
		*bs = append(*bs, b)
	}
}

//Finalize filled sigs,bundlehash and indices elements in bundle.
func (bs Bundle) Finalize(sig []Trytes) {
	h := bs.Hash()
	for i := range bs {
		if len(sig) > i && sig[i] != "" {
			bs[i].SignatureMessageFragment = pad(sig[i], SignatureMessageFragmentTrinarySize/3)
		}
		bs[i].CurrentIndex = int64(i)
		bs[i].LastIndex = int64(len(bs) - 1)
		bs[i].Bundle = h
	}
}

//Hash calculates hash of Bundle.
func (bs Bundle) Hash() Trytes {
	c := NewCurl()
	buf := make(Trits, 243+81*3)
	for i, b := range bs {
		copy(buf, Trytes(b.Address).Trits())
		copy(buf[243:], Int2Trits(b.Value, 81))
		copy(buf[243+81:], b.Tag.Trits())
		copy(buf[243+81+81:], Int2Trits(b.Timestamp.Unix(), 27))
		copy(buf[243+81+81+27:], Int2Trits(int64(i), 27))            //CurrentIndex
		copy(buf[243+81+81+27+27:], Int2Trits(int64(len(bs)-1), 27)) //LastIndex
		c.Absorb(buf.Trytes())
	}
	return c.Squeeze()
}

//Categorize Categorizes a list of transfers into sent and received.
//It is important to note that zero value transfers (which for example,
//is being used for storing addresses in the Mesh), are seen as received in this function.
func (bs Bundle) Categorize(adr Address) (send Bundle, received Bundle) {
	send = make(Bundle, 0, len(bs))
	received = make(Bundle, 0, len(bs))
	for _, b := range bs {
		if b.Address != adr {
			continue
		}
		if b.Value >= 0 {
			received = append(received, b)
		} else {
			send = append(send, b)
		}
	}
	return
}

//IsValid checks the validity of Bundle.
//It checks total balance==0 and its signature.
//You must call Finalize() beforehand.
func (bs Bundle) IsValid() error {
	var total int64
	sigs := make(map[Address][]Trytes)
	for index, b := range bs {
		total += b.Value
		if b.CurrentIndex != int64(index) {
			return fmt.Errorf("CurrentIndex of index %d is not correct", b.CurrentIndex)
		}
		if b.LastIndex != int64(len(bs)-1) {
			return fmt.Errorf("LastIndex of index %d is not correct", b.CurrentIndex)
		}
		if b.Value >= 0 {
			continue
		}
		sigs[b.Address] = append(sigs[b.Address], b.SignatureMessageFragment)
		// Find the subsequent txs with the remaining signature fragment
		for i := index; i < len(bs)-1; i++ {
			tx := bs[i+1]
			// Check if new tx is part of the signature fragment
			if tx.Address == b.Address && tx.Value == 0 {
				sigs[tx.Address] = append(sigs[tx.Address], tx.SignatureMessageFragment)
			}
		}
	}
	// Validate the signatures
	h := bs.Hash()
	for adr, sig := range sigs {
		if !IsValidSig(adr, sig, h) {
			return errors.New("invalid signature")
		}
	}
	if total != 0 {
		return errors.New("total balance of Bundle is not 0")
	}
	return nil
}
