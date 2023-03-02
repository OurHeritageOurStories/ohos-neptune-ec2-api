package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func requestToNeptune(c echo.Context) error {
	sparqlString := c.FormValue("sparqlquery")

	limit := c.FormValue("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return c.String(http.StatusBadRequest, "Limit needs to be a number")
	}

	params := url.Values{}

	maxLimit := os.Getenv("LIMIT")
	maxLimitInt, err2 := strconv.Atoi(maxLimit)
	if err2 != nil {
		return c.String(http.StatusInternalServerError, "Max limit not a number")
	}

	var limitToUse int
	if limitInt > maxLimitInt {
		limitToUse = maxLimitInt
	} else {
		limitToUse = limitInt
	}

	constructedQuery := sparqlString + " LIMIT " + strconv.Itoa(limitToUse)

	params.Add("query", constructedQuery)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return c.String(http.StatusOK, "Search query"+sparqlString+"	Constructed query: "+constructedQuery+"	Results: "+string(data))
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error { //make sure its alive
		return c.HTML(http.StatusOK, "Hello, you've reached the Go API that lets you talk to the Neptune database. Well done!")
	})

	e.POST("/sparql", requestToNeptune) //for actual requests

	e.Logger.Fatal(e.Start(":9000"))

}
