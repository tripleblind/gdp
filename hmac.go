package gdp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"

	// "github.com/codahale/blake2"
)

type HMAC []byte

func i2b(i int) []byte {

	var b = make([]byte, 1)
	binary.PutUvarint(b, uint64(i))

	return b

}

func (h HMAC) F1(ax []byte, L, i1, t0 int) []byte {
	return h.Generate(ax, i2b(L), i2b(i1), i2b(t0))
}

func (h HMAC) F2(ax []byte, L, i1, t0 int, h0 []byte) []byte {
	return h.Generate(ax, i2b(L), i2b(i1), i2b(t0), h0)
}

func (h HMAC) F3(h0, ax []byte, L, s, is, isp1 int) []byte {
	return h.Generate(h0, ax, i2b(L), i2b(s), i2b(is), i2b(isp1))
}

func (h HMAC) F4(msm1, ax []byte, L, s, is, isp1, ts int) []byte {
	return h.Generate(msm1, ax, i2b(L), i2b(s), i2b(is), i2b(isp1), i2b(ts))
}

func F5(e ...[]byte) []byte {

	if len(e) == 0 {
		panic("no elements passed to F5")
	}

	var (
		l   = len(e[0])
		buf []byte
	)

	for i := 0; i < len(e); i++ {

		if buf == nil {
			buf = e[i]
		} else {

			// xor current e and buf
			for i2 := 0; i2 < l; i2++ {
				buf[i2] = buf[i2] ^ e[i][i2]
			}

		}

	}

	return buf

}

func (h HMAC) F6(h0, ax []byte, L, tl int) []byte {
	return h.Generate(h0, ax, i2b(L), i2b(tl))
}

func (h HMAC) Generate(parts ...[]byte) []byte {

	hmac := hmac.New(sha1.New, h)

	for _, e := range parts {

		if _, err := hmac.Write(e); err != nil {
			panic(err)
		}

	}

	return hmac.Sum(nil)

}
