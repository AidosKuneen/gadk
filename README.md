[![GoDoc](https://godoc.org/github.com/AidosKuneen/gadk?status.svg)](https://godoc.org/github.com/AidosKuneen/gadk)
[![Build Status](https://travis-ci.org/AidosKuneen/gadk.svg?branch=master)](https://travis-ci.org/AidosKuneen/gadk)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/AidosKuneen/gadk/LICENSE)

gadk
=====

Client library for the ADK.


Install
====

You will need C compiler (gcc for linux, mingw for windows) to compile PoW routine in C.




Examples
====

```go

import "github.com/AidosKuneen/gadk"

//Trits
tritsFrom:=[]int8{1,-1,1,0,1,1,0,-1,0}
trits,err:=gadk.ToTrits(tritsFrom)

//Trytes
trytes:=trits.Trytes()
trytesFrom:="ABCDEAAC9ACB9PO..."
trytes2,err:=gadk.ToTrytes(trytesFrom)

//Hash
hash:=trytes.Hash()

//API
api := gadk.NewAPI("http://localhost:14265", nil)
resp, err := api.FindTransactions([]Trytes{"DEXRPL...SJRU"})

///Address
index:=0
security:=2
adr,err:=gadk.NewAddress(trytes,index,seciruty) //without checksum.
adrWithChecksum := adr.WithChecksum() //adrWithChecksum is trytes type.

//transaction
tx,err:=gadk.NewTransaction(trytes)
if tx.HasValidNonce(){...}
trytes2:=tx.trytes()

//create signature
key := gadk.NewKey(seed, index, security)
norm := bundleHash.Normalize()
sign := gadk.Sign(norm[:27], key[:6561/3])

//validate signature
if gadk.ValidateSig(adr, []Trytes{sign}, bundleHash) {...}

//send
trs := []gadk.Transfer{
	gadk.Transfer{
		Address: "KTXF...QTIWOWTY",
		Value:   20,
		Tag: "MOUDAMEPO",
	},
}
_, pow := gadk.GetBestPow()
bdl, err = gadk.Send(api, seed, security, trs, pow)
```

PoW(Proof of Work) Benchmarking
====

You can benchmark PoWs(by C,Go,SSE) by

```
    $ go test -v -run Pow
```

or if you want to add OpenCL PoW,

```
    $ go test -tags=gpu -v -run Pow
```

then it outputs like:

```
	$ go test -tags=gpu -v -run Pow
=== RUN   TestPowC
--- PASS: TestPowC (15.93s)
	pow_c_test.go:50: 1550 kH/sec on C PoW
=== RUN   TestPowCL
--- PASS: TestPowCL (17.45s)
	pow_cl_test.go:49: 332 kH/sec on GPU PoW
=== RUN   TestPowGo
--- PASS: TestPowGo (21.21s)
	pow_go_test.go:50: 1164 kH/sec on Go PoW
=== RUN   TestPowSSE
--- PASS: TestPowSSE (13.41s)
	pow_sse_test.go:52: 2292 kH/sec on SSE PoW
```



