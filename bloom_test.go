package bloomc

import (
	"testing"
)

func TestHash(t *testing.T) {
	cases := []struct {
		key    string
		bitNum uint64
	}{
		{``, 53},
		{`hello`, 12},
		{`world`, 21},
		{`1000`, 63},
		{`2000`, 35},
		{`10000`, 5},
		{`100000`, 24},
		{`golang`, 38},
		{"\n\n\n", 13},
		{"\t\t\t", 43},
	}

	b := New(64, 1)

	for i, c := range cases {
		if bn := b.hashData([]byte(c.key), i) % 64; bn != c.bitNum {
			t.Errorf(`expected bit number %d for key %s, given %d`, c.bitNum, c.key, bn)
		}
	}
}

func TestA(t *testing.T) {
	b := New(64, 10)

	b.Add([]byte(`hello`))

	if !b.Test([]byte(`hello`)) {
		t.Error(`test key, expected true`)
	}

	if b.Test([]byte(`golang'`)) {
		t.Error(`test key, expected false`)
	}

	b.Del([]byte(`hello`))

	if b.Test([]byte(`hello`)) {
		t.Error(`test key, expected true`)
	}
}
