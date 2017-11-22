package gadk

//constants for Sizes.
const (
	stateSize = 729
)

var (
	transformC func(Trits)
	truthTable = [11]int8{1, 0, -1, 0, 1, -1, 0, 0, -1, 1, 0}
	indices    [stateSize + 1]int
)

func init() {
	for i := 0; i < stateSize; i++ {
		p := -365
		if indices[i] < 365 {
			p = 364
		}
		indices[i+1] = indices[i] + p
	}
}

// Curl is a sponge function with an internal state of size StateSize.
// b = r + c, b = StateSize, r = HashSize, c = StateSize - HashSize
type Curl struct {
	state Trits
}

// NewCurl initializes a new instance with an empty state.
func NewCurl() *Curl {
	c := &Curl{
		state: make(Trits, stateSize),
	}
	return c
}

//Squeeze do Squeeze in sponge func.
func (c *Curl) Squeeze() Trytes {
	ret := c.state[:HashSize].Trytes()
	c.Transform()

	return ret
}

// Absorb fills the internal state of the sponge with the given trits.
func (c *Curl) Absorb(inn Trytes) {
	in := inn.Trits()
	var lenn int
	for i := 0; i < len(in); i += lenn {
		lenn = 243
		if len(in)-i < 243 {
			lenn = len(in) - i
		}
		copy(c.state, in[i:i+lenn])
		c.Transform()
	}
}

// Transform does Transform in sponge func.
func (c *Curl) Transform() {
	if transformC != nil {
		transformC(c.state)
		return
	}
	var cpy [stateSize]int8
	for r := 27; r > 0; r-- {
		copy(cpy[:], c.state)
		c.state = c.state[:stateSize]
		for i := 0; i < stateSize; i++ {
			t1 := indices[i]
			t2 := indices[i+1]
			c.state[i] = truthTable[cpy[t1]+(cpy[t2]<<2)+5]
		}
	}
}

// Reset the internal state of the Curl sponge by filling it with all
// 0's.
func (c *Curl) Reset() {
	for i := range c.state {
		c.state[i] = 0
	}
}

//Hash returns hash of t.
func (t Trytes) Hash() Trytes {
	c := NewCurl()
	c.Absorb(t)
	return c.Squeeze()
}
