package gadk

import "testing"

func TestTrinaryValidTryte(t *testing.T) {
	type validTryteTC struct {
		in    rune
		valid bool
	}

	var validTryteCases = []validTryteTC{
		validTryteTC{in: 'A', valid: true},
		validTryteTC{in: 'Z', valid: true},
		validTryteTC{in: '9', valid: true},
		validTryteTC{in: '8', valid: false},
		validTryteTC{in: 'a', valid: false},
		validTryteTC{in: '-', valid: false},
		validTryteTC{in: 'Ɩ', valid: false},
	}

	for _, tc := range validTryteCases {
		if (IsValidTryte(tc.in) == nil) != tc.valid {
			t.Fatalf("ValidTryte(%q) should be %#v but is not", tc.in, tc.valid)
		}
	}
}

func TestTrinaryValidTrytes(t *testing.T) {
	type validTryteTC struct {
		in    Trytes
		valid bool
	}

	var validTryteCases = []validTryteTC{
		validTryteTC{in: "ABCDEFGHIJKLMNOPQRSTUVWXYZ9", valid: true},
		validTryteTC{in: "ABCDEFGHIJKLMNOPQRSTUVWXYZ90", valid: false},
		validTryteTC{in: "ABCDEFGHIJKLMNOPQRSTUVWXYZ9 ", valid: false},
		validTryteTC{in: "Ɩ", valid: false},
	}

	for _, tc := range validTryteCases {
		if (tc.in.IsValid() == nil) != tc.valid {
			t.Fatalf("ValidTrytes(%q) should be %#v but is not", tc.in, tc.valid)
		}
	}
}

func TestTrinaryValidTrit(t *testing.T) {
	type validTritTC struct {
		in    int8
		valid bool
	}

	var validTritCases = []validTritTC{
		validTritTC{in: -1, valid: true},
		validTritTC{in: 0, valid: true},
		validTritTC{in: 1, valid: true},
		validTritTC{in: -2, valid: false},
		validTritTC{in: 2, valid: false},
	}

	for _, tc := range validTritCases {
		if (IsValidTrit(tc.in) == nil) != tc.valid {
			t.Fatalf("ValidTrit(%q) should be %#v but is not", tc.in, tc.valid)
		}
	}
}

func TestTrinaryValidTrits(t *testing.T) {
	type validTritsTC struct {
		in    Trits
		valid bool
	}

	var validTritsCases = []validTritsTC{
		validTritsTC{in: Trits{0}, valid: true},
		validTritsTC{in: Trits{-1}, valid: true},
		validTritsTC{in: Trits{1}, valid: true},
		validTritsTC{in: Trits{0, -1, 1}, valid: true},
		validTritsTC{in: Trits{2, -1, 1}, valid: false},
	}

	for _, tc := range validTritsCases {
		if (tc.in.IsValid() == nil) != tc.valid {
			t.Fatalf("ValidTrits(%q) should be %#v but is not", tc.in, tc.valid)
		}
	}
}

func TestTrinaryConvert(t *testing.T) {
	trits := Trits{0, 1, -1, 1, 1, -1, -1, 1, 1, 0, 0, 1, 0, 1, 1}
	invalid := []int8{1, -1, 2, 0, 1, -1}

	if _, err := ToTrits(invalid); err == nil {
		t.Error("ToTrits is incorrect")
	}

	if _, err := ToTrytes("A_AAA"); err == nil {
		t.Error("ToTrytes is incorrect")
	}

	i := trits.Int()
	if i != 6562317 {
		t.Error("Int() is illegal.", i)
	}
	trits2 := Int2Trits(6562317, 15)
	if !trits.Equal(trits2) {
		t.Error("Int2Trits() is illegal.", trits2)
	}
	trits22 := Int2Trits(-1024, 7)
	if !trits22.Equal(Trits{-1, 1, 0, 1, -1, -1, -1}) {
		t.Error("Int2Trits() is illegal.")
	}
	try := trits.Trytes()
	if try != "UVKIL" {
		t.Error("Int() is illegal.", try)
	}
	trits3 := try.Trits()
	if !trits.Equal(trits3) {
		t.Error("Trits() is illegal.", trits3)
	}
}

func TestTrinaryNormalize(t *testing.T) {
	var bundleHash Trytes = "DEXRPLKGBROUQMKCLMRPG9HFKCACDZ9AB9HOJQWERTYWERJNOYLW9PKLOGDUPC9DLGSUH9UHSKJOASJRU"
	no := []int8{-13, -13, -13, -13, -11, 12, 11, 7, 2, -9, -12, -6, -10, 13, 11, 3, 12, 13, -9, -11, 7, 0, 8, 6,
		11, 3, 1, 13, 13, 13, 7, 1, 2, 0, 8, -12, 10, -10, -4, 5, -9, -7, -2, -4, 5, -9, 10, -13, -12, -2, 12, -4,
		0, -11, -5, 12, -12, 7, 4, -6, -11, 3, 0, 4, 12, 7, -8, -6, 8, 0, -6, 8, -8, 11, 10, -12, 1, -8, 10, -9, -6}
	norm := bundleHash.Normalize()
	for i := range no {
		if no[i] != norm[i] {
			t.Fatal("normalization is incorrect.")
		}
	}
}
