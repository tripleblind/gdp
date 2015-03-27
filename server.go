package gdp

type Server struct {
	SecretKey  []byte
	SharedKeys [][]byte
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
