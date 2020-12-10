package bloomc

import (
	"bufio"
	"hash"
	"hash/fnv"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Bloom struct {
	counter []uint64
	size    uint64
	keys    int

	hashFNV hash.Hash64
	mu      sync.RWMutex
}

func New(size uint64, keys int) *Bloom {
	return &Bloom{
		counter: make([]uint64, size),
		keys:    keys,
		size:    size,

		hashFNV: fnv.New64(),
	}
}

func (b *Bloom) Reset() {
	b.counter = make([]uint64, len(b.counter))
}

func (b *Bloom) Test(data []byte) bool {
	for n := 0; n < b.keys; n++ {
		b.mu.RLock()
		if b.counter[b.hashData(data, n)%b.size] == 0 {
			b.mu.RUnlock()
			return false
		}
		b.mu.RUnlock()
	}

	return true
}

func (b *Bloom) Add(data []byte) {
	for n := 0; n < b.keys; n++ {
		b.mu.Lock()
		b.counter[b.hashData(data, n)%b.size]++
		b.mu.Unlock()
	}
}

func (b *Bloom) Del(data []byte) {
	for n := 0; n < b.keys; n++ {
		bitNum := b.hashData(data, n) % b.size

		b.mu.Lock()
		if b.counter[bitNum] > 0 {
			b.counter[bitNum]--
		}
		b.mu.Unlock()
	}
}

func (b *Bloom) LoadFile(file string) error {
	fh, err := os.OpenFile(file, os.O_CREATE|os.O_RDONLY, 0664)
	if err != nil {
		return err
	}
	defer fh.Close()

	b.Reset()

	s := bufio.NewScanner(fh)

	var i int
	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if len(line) == 0 {
			continue
		}

		n, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			continue
		}

		b.counter[i] = uint64(n)
		i++
	}

	return nil
}

func (b *Bloom) SaveFile(file string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	fh, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer fh.Close()

	for _, bit := range b.counter {
		if _, err = fh.WriteString(strconv.Itoa(int(bit)) + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bloom) hashData(data []byte, i int) uint64 {
	algo := b.hashFNV

	algo.Reset()
	_, _ = algo.Write(data)
	_, _ = algo.Write([]byte{
		byte(0xff & i),
		byte(0xff & (i >> 8)),
		byte(0xff & (i >> 16)),
		byte(0xff & (i >> 24)),
	})
	return algo.Sum64()
}
