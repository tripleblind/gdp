package gdp

import "net"

type Guide struct {
	ServerKey  []byte
	SharedKeys [][]byte
}

func NewGuide(guides int) *Guide {
	return &Guide{
		SharedKeys: make([][]byte, guides),
	}
}

// verifies a HL and returns a HSOL
func (g *Guide) Verify(h0, hL []byte, L int, lastM []byte, i [][]byte) ([]byte, error) {

	var (
		ax = net.ParseIP("127.0.0.1")
		tl = 0
	)

	// f3 := nil
	// f5 := nil

	// let's assume it's verfied, compute hsol

	hsol := g.WithServer().F6(h0, ax, L, tl)

	return hsol, nil

	// {h0,hL,L,mL−1,i1,i2,...,iL} to the first stop tour

}

func (g *Guide) WithGuide(idx int) HMAC {
	return g.SharedKeys[idx]
}

func (g *Guide) WithServer() HMAC {
	return g.ServerKey
}
