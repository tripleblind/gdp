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

func (g *Guide) Visit(ax []byte, r *Request) (*Reply, bool, error) {

	// TODO: verify signature

	// how?!?!?

	// g.WithGuide(r.ISM1).F4(msm1, ax, L, s, is, isp1, ts)

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

	reply.HS = g.WithGuide(r.I1).F3(r.H0, ax, r.L, r.S, r.IS, isp1)

	reply.MS = g.WithGuide(isp1).F4(r.MSM1, ax, r.L, r.S, r.IS, isp1, reply.TS)

	log.Printf("[  tour] visiting stop %d at index %d (%s) (h: %x, m: %x)", r.S, r.IS, g.Name, reply.HS, reply.MS)

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
