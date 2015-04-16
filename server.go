package gdp

type Server struct {
	SecretKey  []byte
	SharedKeys [][]byte
	TourLength int
}

func (s *Server) NewTour(ax []byte) *Tour {

	t := &Tour{
		L: s.TourLength,
		S: 0,
		I: []int{
			dice(Guides),
		},
		T: []int{
			0,
		},
	}

	t.H = [][]byte{
		s.WithSecret().F1(ax, t.L, t.I[0], t.T[0]),
	}

	t.M = [][]byte{
		s.WithGuide(t.I[0]).F2(ax, t.L, t.I[0], t.T[0], t.H[0]),
	}

	return t

}

func (s *Server) NewPuzzle(ax []byte) *Puzzle {

	p := &Puzzle{
		server: s,
		L:      s.TourLength,
		I1:     dice(Guides),
		T0:     0,
	}

	p.H0 = s.WithSecret().F1(ax, p.L, p.I1, p.T0)
	p.M0 = s.WithGuide(p.I1).F2(ax, p.L, p.I1, p.T0, p.H0)

	return p

}

func (s *Server) Verify(h0, hSol, t0, tL []byte, i1 int) error {

	// var (
	//  ax = net.ParseIP("127.0.0.1")
	// )

	// f6 := s.withGuide(i1).Generate(h0, ax)

	return nil

}

func (s *Server) WithGuide(idx int) HMAC {
	return s.SharedKeys[idx]
}

func (s *Server) WithSecret() HMAC {
	return s.SecretKey
}
