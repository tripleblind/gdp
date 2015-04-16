package gdp

import "crypto/hmac"

type Puzzle struct {
	server *Server
	L      int
	I1     int
	T0     int
	H0, M0 []byte
}

func (p *Puzzle) Verify(ax, hsol []byte) bool {

	result := p.server.WithGuide(p.I1).F6(p.H0, ax, p.L, 0)

	return hmac.Equal(hsol, result)

}
