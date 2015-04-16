package gdp

import (
	"crypto/hmac"
	"fmt"
	"log"
)

type Name int

func (n Name) String() string {

	switch n {
	case North:
		return "north"
	case East:
		return "east"
	case South:
		return "south"
	case West:
		return "west"
	default:
		panic(fmt.Sprintf("Unknown tour guide index %q", n))
	}

}

const (
	North  Name = 0
	East        = 1
	South       = 2
	West        = 3
	Guides      = 4
)

type Guide struct {
	Name       Name
	ServerKey  []byte
	SharedKeys [][]byte
}

func NewGuide(name Name) *Guide {
	return &Guide{
		Name:       name,
		SharedKeys: make([][]byte, Guides),
	}
}

func (g *Guide) Verify(ax []byte, tour *Tour) error {

	for s := 1; s == tour.L; s++ {

		h := g.WithGuide(tour.I[s]).F3(
			tour.H[0],
			ax,
			tour.L,
			s,
			tour.I[s],
			tour.I[s+1],
		)

		if !hmac.Equal(h, tour.H[s]) {
			err := fmt.Errorf("Invalid H at stop %d: see %x, calculated %x", s, tour.H[s], h)
			return err
		}

	}

	ts := 0

	tour.Sol = g.WithServer().F6(
		tour.H[0],
		ax,
		tour.L,
		ts,
	)

	return nil

}

func (g *Guide) Visit(ax []byte, prev *Tour) (*Tour, bool, error) {

	log.Printf("Visit at stop %d: I %d T %d H %d", prev.S, len(prev.I), len(prev.T), len(prev.H))

	initial := 0

	// TODO: verify invariants (maybe!)

	// verify signatures

	var m []byte

	if prev.S == initial {

		log.Println("Generating initial M")

		m = g.WithServer().F2(
			ax,
			prev.L,
			prev.I[0],
			prev.T[0],
			prev.H[0],
		)

	} else {

		log.Println("Generating followup M")

		m = g.WithGuide(prev.I[prev.S-1]).F4(
			prev.M[prev.S-1],
			ax,
			prev.L,
			prev.S-1,
			prev.I[prev.S-1],
			prev.I[prev.S],
			prev.T[prev.S-1],
		)

	}

	prevM := prev.M[prev.S]

	if !hmac.Equal(m, prevM) {
		err := fmt.Errorf("Invalid M (%x != %x)", m, prevM)
		return nil, false, err
	} else {
		log.Println("Signatures verified")
	}

	final := prev.S == prev.L-1

	var nextI int

	if final {
		nextI = prev.I[0] // last step
	} else {
		nextI = dice(Guides)
	}

	next := *prev

	next.T = append(prev.T, 0)
	next.I = append(prev.I, nextI)

	// calculate new signatures

	ts := 0

	next.S = next.S + 1

	next.H = append(next.H, g.WithGuide(prev.I[0]).F3(
		prev.H[0],
		ax,
		prev.L,
		prev.S,
		prev.I[prev.S],
		nextI,
	))

	log.Printf("Calculating H for S %d: %x", next.S, next.H[next.S])

	next.M = append(next.M, g.WithGuide(prev.I[prev.S]).F4(
		prev.M[prev.S],
		ax,
		prev.L,
		prev.S,
		prev.I[prev.S],
		nextI,
		ts,
	))

	return &next, !final, nil

}

func (g *Guide) WithGuide(idx int) HMAC {
	return g.SharedKeys[idx]
}

func (g *Guide) WithServer() HMAC {
	return g.ServerKey
}
