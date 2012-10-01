package hash

import (
	"github.com/dgryski/dgohash"
)

type Hash uint32

func (h Hash) Intgr() uint32 {
	return uint32(h)
}

func (h Hash) Str() string {
	lookup := []byte{
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9',
		'A',
		'B',
		'C',
		'D',
		'E',
		'F',
	}
	return string([]byte{
		lookup[(h>>28)&0xf],
		lookup[(h>>24)&0xf],
		lookup[(h>>20)&0xf],
		lookup[(h>>16)&0xf],
		lookup[(h>>12)&0xf],
		lookup[(h>>8)&0xf],
		lookup[(h>>4)&0xf],
		lookup[h&0xf],
	})
}

func Rehash(str string) (Hash, bool) {
	lookup := map[rune]uint32{
		'0': 0x0,
		'1': 0x1,
		'2': 0x2,
		'3': 0x3,
		'4': 0x4,
		'5': 0x5,
		'6': 0x6,
		'7': 0x7,
		'8': 0x8,
		'9': 0x9,
		'A': 0xa,
		'B': 0xb,
		'C': 0xc,
		'D': 0xd,
		'E': 0xe,
		'F': 0xf,
		'a': 0xa,
		'b': 0xb,
		'c': 0xc,
		'd': 0xd,
		'e': 0xe,
		'f': 0xf,
	}
	var result Hash
	for _, b := range str {
		if val, ok := lookup[b]; ok {
			result = Hash((uint32(result) << 4) | val)
		} else {
			return Hash(0), false
		}
	}
	return result, true
}

func Gethash(str string) Hash {
	m := dgohash.NewMurmur3_x86_32()
	m.Write([]byte(str))
	return Hash(m.Sum32())
}
