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

	server.TourLength = 3

}

func main() {

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {

		ax := gdp.ClientIdentity(c.Request.RemoteAddr)

		tour := server.NewTour(ax)

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
