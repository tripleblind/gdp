package gdp

import "crypto/hmac"

type Puzzle struct {
	server *Server
	L      int
	I1     int
	T0     int
	H0, M0 []byte
}

func (s *Server) NewPuzzle(ax []byte, l int) *Puzzle {

	p := &Puzzle{
		server: s,
		L:      l,
		I1:     dice(Guides),
		T0:     0,
	}

	p.H0 = s.WithSecret().F1(ax, p.L, p.I1, p.T0)
	p.M0 = s.WithSecret().F2(ax, p.L, p.I1, p.T0, p.H0)

	return p

}

func (s *Server) NewPuzzleFromJSON() (*Puzzle, error) {
	return nil, nil // TBD
}

func (p *Puzzle) Verify(ax, hsol []byte) bool {

	result := p.server.WithGuide(p.I1).F6(p.H0, ax, p.L, 0)

	return hmac.Equal(hsol, result)

}
