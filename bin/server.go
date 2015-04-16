package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/usersjs/gdp"
)

func s2b(s string) []byte {
	return []byte(s)
}

var server *gdp.Server

func init() {

	server = &gdp.Server{
		SecretKey: s2b("server+server"),
		SharedKeys: [][]byte{
			s2b("server+north"),
			s2b("server+east"),
			s2b("server+south"),
			s2b("server+west"),
		},
	}

	server.TourLength = 10

}

func main() {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {

		ax := gdp.ClientIdentity(c.Request.RemoteAddr)

		tour := server.NewTour(ax)

		// puzzle := server.NewPuzzle(ax)

		// req := gdp.Request{ // TODO: add to *Puzzle
		// 	H0:   puzzle.H0,
		// 	L:    puzzle.L,
		// 	S:    1,
		// 	TSM1: puzzle.T0,
		// 	MSM1: puzzle.M0,
		// 	I1:   puzzle.I1,
		// 	IS:   puzzle.I1,
		// }

		out, _ := json.Marshal(tour)

		query := url.Values{}
		query.Add("t", string(out))

		link := url.URL{
			Scheme:   "http",
			Host:     fmt.Sprintf("127.0.0.1:%d", 10000*(tour.I[0]+1)),
			Path:     "visit",
			RawQuery: query.Encode(),
		}

		//		c.Writer.Header().Add("Link", fmt.Sprintf(`<%s>; rel="next"`, link.String()))

		tour.Link = link.String()

		// c.Redirect(http.StatusTemporaryRedirect, link.String())
		c.JSON(http.StatusPaymentRequired, tour)

	})

	router.Run("127.0.0.1:9000")

}
