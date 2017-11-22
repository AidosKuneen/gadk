package gadk

//Various constants for gadk.
const (
	TryteAlphabet      = "9ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	MinTryteValue      = -13
	MaxTryteValue      = 13
	SignatureSize      = 6561
	HashSize           = 243
	MinWeightMagnitude = 18 //must be 18.
	Depth              = 3
)

//Units for adk token.
const (
	Ki = 1000
	Mi = 1000000
	Gi = 1000000000
	Ti = 1000000000000
	Pi = 1000000000000000
)

var (
	//emptySig represents empty signature.
	emptySig Trytes
	//EmptyHash represents empty hash.
	EmptyHash Trytes = "999999999999999999999999999999999999999999999999999999999999999999999999999999999"
	//EmptyAddress represents empty address.
	EmptyAddress Address = "999999999999999999999999999999999999999999999999999999999999999999999999999999999"
)

func init() {
	bytes := make([]byte, SignatureSize/3)
	for i := 0; i < SignatureSize/3; i++ {
		bytes[i] = '9'
	}
	emptySig = Trytes(bytes)
}
