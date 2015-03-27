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

func (g *Guide) Visit() {

}

// verifies a HL and returns a HSOL
func (g *Guide) Verify(h0, hl []byte, L int, lastM []byte, i []int) ([]byte, error) {

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
