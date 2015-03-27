package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net"
)

type Guide struct {
	ServerKey  []byte
	SharedKeys [][]byte
}

func NewGuide() *Guide {
	return &Guide{
		SharedKeys: make([][]byte, Guides),
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

	hsol := g.withServer().F6(h0, ax, L, tl)

	return hsol, nil

	// {h0,hL,L,mL−1,i1,i2,...,iL} to the first stop tour

}

func (g *Guide) withGuide(idx int) HMAC {
	return g.SharedKeys[idx]
}

func (g *Guide) withServer() HMAC {
	return g.ServerKey
}

type Server struct {
	SecretKey  []byte
	SharedKeys [][]byte
}

func (s *Server) Verify(h0, hSol, t0, tL []byte, i1 int) error {

	// var (
	// 	ax = net.ParseIP("127.0.0.1")
	// )

	// f6 := s.withGuide(i1).Generate(h0, ax)

	return nil

}

func (s *Server) withGuide(idx int) HMAC {
	return s.SharedKeys[idx]
}

func (s *Server) withSecret() HMAC {
	return s.SecretKey
}

const (
	KeySize = 16
	Guides  = 4
)

func dice(up int) int {

	if n, err := rand.Int(rand.Reader, big.NewInt(int64(up))); err != nil {
		panic(err)
	} else {
		return int(n.Int64())
	}

}

func keygen(size int) (buf []byte) {

	buf = make([]byte, size)

	n, err := rand.Read(buf)

	if n != size {
		err = fmt.Errorf("unexpected random bytes, want %d got %d", size, n)
	}

	if err != nil {
		panic(err)
	}

	return buf

}

func main() {

	server := &Server{
		SecretKey: keygen(KeySize),
	}

	guides := []*Guide{
		NewGuide(),
		NewGuide(),
		NewGuide(),
		NewGuide(),
	}

	for i, e := range guides {

		server.SharedKeys = append(server.SharedKeys, e.ServerKey)

		for i2, _ := range guides {

			if i != i2 {

				if key := guides[i2].SharedKeys[i]; key != nil {
					log.Printf("guide %d: COPY shared key from %d", i2, i)
					e.SharedKeys[i2] = key
				} else {
					log.Printf("guide %d: MAKE shared key for %d", i2, i)
					e.SharedKeys[i2] = keygen(KeySize)
				}

			}

		}

	}

	log.Println("symetry self-test for server-guide-shared key")

	// debug the key relationships are setup symetrical
	for i, _ := range guides {

		msg := []byte("this is a secret")

		a := server.withGuide(i).generate(msg)
		b := guides[i].withServer().generate(msg)

		log.Printf("guide %d with server : %x == %x = %v",
			i,
			a,
			b,
			bytes.Equal(a, b),
		)

		for i2, _ := range guides {

			if i != i2 {

				a = guides[i].withGuide(i2).generate(msg)
				b = guides[i2].withGuide(i).generate(msg)

				log.Printf("guide %d with guide %d: %x == %x = %v",
					i,
					i2,
					a,
					b,
					bytes.Equal(a, b),
				)

			}

		}

	}

	// puzzle generation

	var (
		ax = net.ParseIP("127.0.0.1")
		L  = 3
		i1 = dice(Guides)
		t0 = 0
	)

	log.Printf("[server] requesting tour of length %d, starting at %d", L, i1)

	h0 := server.withSecret().F1(ax, L, i1, t0)
	m0 := server.withSecret().F2(ax, L, i1, t0, h0)

	var (
		allH [][]byte
		allM = [][]byte{m0}
	)

	// begin the tour

	var (
		ts   = 0
		S    = 1            // current stop
		iS   = i1           // current index
		iSp1 = dice(Guides) // successor index
	)

	next := func() bool {

		next := S < L

		S = S + 1
		iS = iSp1

		if S == L-1 {
			iSp1 = i1
		} else {
			iSp1 = dice(Guides)
		}

		return next

	}

	for {

		h := guides[iS].withGuide(i1).F3(h0, ax, L, S, iS, iSp1)
		m := guides[iS].withGuide(iSp1).F4(allM[len(allM)-1], ax, L, S, iS, iSp1, ts)

		log.Printf("[  tour] visiting stop %d at index %d (h: %x, m: %x)", S, iS, h, m)

		allH = append(allH, h)
		allM = append(allM, m)

		// TODO: verify TS

		if !next() {
			break
		}

	}

	hl := F5(allH...)

	log.Printf("[client] tour complete, generating hl: %x", hl)

	// tour verifies this ...

	hsol, err := guides[i1].Verify(h0, nil, L, nil, allM)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[  tour] validate hl and return hsol %x at %d", hsol, i1)

	// client passes hsol etc to the server

	result := server.withGuide(i1).F6(h0, ax, L, 0)

	log.Printf("[server] %x = %x == %v", hsol, result, bytes.Equal(hsol, result))

}
