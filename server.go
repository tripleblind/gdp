package tour

import (
	"crypto/hmac"
	"fmt"
	"log"
)

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

		s.WithSecret().F1(
			ax,
			t.L,
			t.I[0],
			t.T[0],
		),
	}

	t.M = [][]byte{

		s.WithGuide(t.I[0]).F2(
			ax,
			t.L,
			t.I[0],
			t.T[0],
			t.H[0],
		),
	}

	return t

}

func (s *Server) Verify(ax []byte, t *Tour) error {

	h0 := s.WithSecret().F1(ax, t.L, t.I[0], t.T[0])

	if !hmac.Equal(h0, t.H[0]) {
		return fmt.Errorf("Invalid H0")
	} else {
		log.Println("H0 ok")
	}

	ts := 0

	hsol := s.WithGuide(t.I[0]).F6(
		t.H[0],
		ax,
		t.L,
		ts,
	)

	if !hmac.Equal(hsol, t.Sol) {
		return fmt.Errorf("Invalid hsol")
	} else {
		log.Println("hsol ok")
	}

	return nil

}

func (s *Server) WithGuide(idx int) HMAC {
	return s.SharedKeys[idx]
}

func (s *Server) WithSecret() HMAC {
	return s.SecretKey
}
