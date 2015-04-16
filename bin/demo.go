package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net"

	"github.com/usersjs/gdp"
)

const (
	KeySize = 16
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

var (
	server *gdp.Server
	guides []*gdp.Guide
)

func init() {

	server = &gdp.Server{
		SecretKey: keygen(KeySize),
	}

	guides = []*gdp.Guide{
		gdp.NewGuide(gdp.North),
		gdp.NewGuide(gdp.East),
		gdp.NewGuide(gdp.South),
		gdp.NewGuide(gdp.West),
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

		a := server.WithGuide(i).Generate(msg)
		b := guides[i].WithServer().Generate(msg)

		log.Printf("guide %d with server : %x == %x = %v",
			i,
			a,
			b,
			bytes.Equal(a, b),
		)

		for i2, _ := range guides {

			if i != i2 {

				a = guides[i].WithGuide(i2).Generate(msg)
				b = guides[i2].WithGuide(i).Generate(msg)

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

}

func main() {

	server.TourLength = 10

	// puzzle generation

	var (
		ax = net.ParseIP("127.0.0.1")
		p  = server.NewPuzzle(ax)
	)

	log.Printf("[server] requesting tour of length %d, starting at %d", p.L, p.I1)

	// state

	var (
		allH [][]byte
		allI []int
		S    = 1
		IS   = p.I1
		TS   = p.T0
		M    = p.M0
	)

	for {

		reply, next, err := guides[IS].Visit(ax, &gdp.Request{
			H0:   p.H0,
			L:    p.L,
			S:    S,
			ISM1: IS,
			TSM1: TS,
			MSM1: M,
			I1:   p.I1,
			IS:   IS,
		})

		if err != nil {
			panic(err)
		}

		allH = append(allH, reply.HS)
		allI = append(allI, reply.ISP1)

		IS = reply.ISP1
		S = S + 1

		// h := guides[iS].WithGuide(p.I1).F3(p.H0, ax, p.L, S, iS, iSp1)

		// m := guides[iS].WithGuide(iSp1).F4(allM[len(allM)-1], ax, p.L, S, iS, iSp1, ts)

		// log.Printf("[  tour] visiting stop %d at index %d (%s) (h: %x, m: %x)", S, iS, guides[iS].Name, h, m)

		// allH = append(allH, h)
		// allM = append(allM, m)
		// allI = append(allI, iS)

		// // TODO: verify TS

		if !next {
			break
		}

	}

	hl := gdp.F5(allH...)

	log.Printf("[client] tour complete, generating hl: %x", hl)

	// tour verifies this ...

	hsol, err := guides[p.I1].VerifyTour(p.H0, nil, p.L, nil, allI)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[  tour] validate hl and return hsol %x at %d", hsol, p.I1)

	// client passes hsol etc to the server

	log.Printf("[server] VERIFY: %v", p.Verify(ax, hsol))

	// result := server.WithGuide(p.I1).F6(p.H0, ax, p.L, 0)

	// log.Printf("[server] %x = %x == %v", hsol, result, bytes.Equal(hsol, result))

}
