package gdp

import (
	"crypto/hmac"
	"fmt"
	"log"
	"net"
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

func (g *Guide) VisitT(ax []byte, prev *Tour) (*Tour, bool, error) {

	initial := 0

	// TODO: verify invariants

	// ...

	// TODO: verify signatures

	var m []byte

	if prev.S == initial {
		m = g.WithServer().F2(ax, prev.L, prev.I[0], prev.T[0], prev.H[0])
	} else {
		m = g.WithGuide(prev.I[prev.S-1]).F4(prev.M[prev.S-1], ax, prev.L, prev.S-1, prev.I[prev.S-1], prev.I[prev.S], prev.T[prev.S-1])
	}

	prevM := prev.M[prev.S]

	if !hmac.Equal(m, prevM) {
		log.Printf("Invalid M (%x != %x)", m, prevM)
	} else {
		log.Println(" Valid M")
	}

	final := prev.S == prev.L

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

	next.H = append(next.H, g.WithGuide(prev.I[0]).F3(prev.H[0], ax, prev.L, prev.S, prev.I[prev.S], nextI))
	next.M = append(next.M, g.WithGuide(prev.I[prev.S]).F4(prev.M[prev.S], ax, prev.L, prev.S, prev.I[prev.S], nextI, ts))

	next.S = next.S + 1

	return &next, !final, nil

}

func (g *Guide) Visit(ax []byte, r *Request) (*Reply, bool, error) {

	// verify invariants

	if int(g.Name) != r.IS {
		return nil, false, fmt.Errorf("Wrong guide %q for index %q", g.Name, r.IS)
	}

	// verify signatures

	var m []byte

	if r.S == 1 {
		m = g.WithServer().F2(ax, r.L, r.I1, 0, r.H0)
	} else {
		m = g.WithGuide(r.ISM1).Generate(r.MSM1) // F4(r.MSM1, ax, r.L, r.S, r.ISM1, r.IS, r.TSM1)
	}

	log.Printf("WE ARE AT STOP IS %d", int(g.Name))

	if !hmac.Equal(m, r.MSM1) {
		log.Printf("WARNING, checksums wrong")
		// return nil, false, fmt.Errorf("Cannot verify (%d) %x == %x = %v", r.S, m, r.MSM1)
	} else {
		log.Printf("ATTENTION, checksums right")
	}

	// TODO: verify timestamp

	var (
		next  = r.S < r.L
		isp1  int
		reply = &Reply{
			TS: 0,
		}
	)

	if r.S == r.L-1 {
		isp1 = r.I1
	} else {
		isp1 = dice(Guides)
	}

	reply.ISP1 = isp1

	log.Printf("NEXT STOP IS %d", isp1)

	log.Printf("Using MSM1 %x", r.MSM1)

	reply.HS = g.WithGuide(r.I1).F3(r.H0, ax, r.L, r.S, r.IS, isp1)
	reply.MS = g.WithGuide(isp1).Generate(r.MSM1) // F4(r.MSM1, ax, r.L, r.S, r.IS, isp1, reply.TS)

	log.Printf("[  tour] visiting (last: %v) stop %d at index %d (%s) (h: %x, m: %x)", next == false, r.S, r.IS, g.Name, reply.HS, reply.MS)

	return reply, next, nil

}

// verifies a HL and returns a HSOL
func (g *Guide) VerifyTour(h0, hl []byte, L int, lastM []byte, i []int) ([]byte, error) {

	var (
		ax = net.ParseIP("127.0.0.1")
		tl = 0
	)

	// verify

	var h [][]byte

	for s, is := range i {

		var isp1 int

		if s == len(i)-1 {
			isp1 = i[0]
		} else {
			isp1 = i[s+1]
		}

		h = append(h, g.WithServer().F3(h0, ax, L, s, is, isp1))

	}

	f5 := F5(h...)

	// TODO: this does not work yet!

	log.Printf("VERIFY %x = %x == %v", hl, f5, hmac.Equal(hl, f5))

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
