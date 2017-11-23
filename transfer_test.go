package gadk

import (
	"testing"
)

const apiServer = "http://78.46.250.88:15555"

var (
	seed Trytes = "VOQHWAPIKQNYQZRYRMJIYSLPBVLFOTPJMQKKNYDANFTG9ICYDLRUJPCDDWDLD9YEGIKISSHWWHKOWONMN"
)

func TestTransfer1(t *testing.T) {
	var err error
	var adr Address
	var adrs []Address
	for i := 0; i < 5; i++ {
		api := NewAPI(apiServer, nil)
		adr, adrs, err = GetUsedAddress(api, seed, 2)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Error(err)
	}
	t.Log(adr, adrs)
	if len(adrs) < 1 {
		t.Error("GetUsedAddress is incorrect")
	}

	var bal Balances
	for i := 0; i < 5; i++ {
		api := NewAPI(apiServer, nil)
		_, bal, err = GetInputs(api, seed, 0, 10, 1000, 2)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Error(err)
	}
	t.Log(bal)
	if len(bal) < 1 {
		t.Error("GetInputs is incorrect")
	}

}
func TestTransfer2(t *testing.T) {
	var err error
	trs := []Transfer{
		Transfer{
			Address: "ZTBTQDHNZBVJXOJSIMQPUHQORZFALAHWRBYJQMRFTVSDDRLICVGBOEEXIJSMNNSWEVICVAMEZPBVASNSETGEIMKSGA",
			Value:   20,
			Tag:     "MOUDAMEPO",
		},
	}

	var bdl Bundle
	for i := 0; i < 5; i++ {
		api := NewAPI(apiServer, nil)
		bdl, err = PrepareTransfers(api, seed, trs, nil, "", 2)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Error(err)
	}
	if len(bdl) < 3 {
		for _, tx := range bdl {
			t.Log(tx.Trytes())
		}
		t.Fatal("PrepareTransfers is incorrect len(bdl)=", len(bdl))
	}
	if err = bdl.IsValid(); err != nil {
		t.Error(err)
	}
	name, pow := GetBestPoW()
	t.Log("using PoW: ", name)

	for i := 0; i < 5; i++ {
		api := NewAPI(apiServer, nil)
		bdl, err = Send(api, seed, 2, trs, pow)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Error(err)
	}
	for _, tx := range bdl {
		t.Log(tx.Trytes())
	}
}
