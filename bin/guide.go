package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/usersjs/gdp"
)

func s2b(s string) []byte {
	return []byte(s)
}

var guide *gdp.Guide

func init() {

	idx, err := strconv.Atoi(os.Args[1])

	if err != nil {
		panic(err)
	}

	name := gdp.Name(idx)

	guide = &gdp.Guide{
		Name:      name,
		ServerKey: s2b(fmt.Sprintf("server+%s", name)),
		SharedKeys: [][]byte{
			s2b("north"),
			// s2b("north"),
			// s2b("north"),
			// s2b("north"),
			s2b("east"),
			s2b("south"),
			s2b("west"),
		},
	}

	log.Printf("%s (%d) ready: %v", guide.Name, guide.Name, guide.SharedKeys)

}

func main() {

	router := gin.Default()

	router.GET("/verify", func(c *gin.Context) {

	})

	router.GET("/visit", func(c *gin.Context) {

		var tour gdp.Tour

		c.Writer.Header().Add("Server", guide.Name.String())

		if err := json.Unmarshal([]byte(c.Request.URL.Query().Get("t")), &tour); err != nil {
			c.String(http.StatusBadRequest, "Unmarshaling error: %s", err)
		} else {

			ax := gdp.ClientIdentity(c.Request.RemoteAddr)

			step, next, err := guide.VisitT(ax, &tour)

			if err != nil {
				c.String(http.StatusBadRequest, "Invalid request: %s", err)
			} else {

				out, _ := json.Marshal(step)

				query := url.Values{}
				query.Add("t", string(out))

				link := url.URL{
					Scheme:   "http",
					Host:     fmt.Sprintf("127.0.0.1:%d", 10000*(step.I[step.S]+1)),
					Path:     "visit",
					RawQuery: query.Encode(),
				}

				if next {

					step.Link = link.String()

					c.JSON(200, step)

					// c.Redirect(http.StatusTemporaryRedirect, link.String())

				} else {

					c.JSON(http.StatusOK, step)
				}

			}

		}

	})

	router.Run(fmt.Sprintf("127.0.0.1:%d", 10000*(guide.Name+1)))

}
