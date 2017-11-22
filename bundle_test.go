package gadk

import (
	"testing"
	"time"
)

func TestBundle(t *testing.T) {
	var bs Bundle
	adr := []Address{
		"PQTDJXXKSNYZGRJDXEHHMNCLUVOIRZC9VXYLSITYMVCQDQERAHAUZJKRNBQEUHOLEAXRUSQBNYVJWESYR",
		"KTXFP9XOVMVWIXEWMOISJHMQEXMYMZCUGEQNKGUNVRPUDPRX9IR9LBASIARWNFXXESPITSLYAQMLCLVTL",
		"KTXFP9XOVMVWIXEWMOISJHMQEXMYMZCUGEQNKGUNVRPUDPRX9IR9LBASIARWNFXXESPITSLYAQMLCLVTL",
		"GXZWHBLRGGY9BCWCAVTFGHCOEWDBFLBTVTIBOQICKNLCCZIPYGPESAPUPDNBDQYENNMJTWSWDHZTYEHAJ",
	}
	value := []int64{
		50, -100, 0, 50,
	}
	ts := []string{
		"2017-03-11 12:25:05 +0900 JST",
		"2017-03-11 12:25:18 +0900 JST",
		"2017-03-11 12:25:18 +0900 JST",
		"2017-03-11 12:25:28 +0900 JST",
	}
	var hash Trytes = "DEXRPLKGBROUQMKCLMRPG9HFKCACDZ9AB9HOJQWERTYWERJNOYLW9PKLOGDUPC9DLGSUH9UHSKJOASJRU"
	for i := 0; i < 4; i++ {
		tss, err := time.Parse("2006-01-02 15:04:05 -0700 MST", ts[i])
		if err != nil {
			t.Fatal(err)
		}
		bs.Add(1, adr[i], value[i], tss, "")
	}
	if bs.Hash() != hash {
		t.Error("hash of bundles is illegal.", bs.Hash())
	}
	bs.Finalize([]Trytes{})
	send, receive := bs.Categorize(adr[1])
	if len(send) != 1 || len(receive) != 1 {
		t.Error("Categorize is incorrect")
	}
}
