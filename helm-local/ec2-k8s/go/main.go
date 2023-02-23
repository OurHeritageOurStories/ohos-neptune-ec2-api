package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func requestToNeptune(c echo.Context) error {

	//sparqlString := c.FormValue("sqarqlstring")
	//sparqlString := "hi"
	//return c.String(http.StatusOK, sparqlString)
	sparqlString := c.FormValue("sqarqlstring")
	constructedSparqlQuery := "query=select " + sparqlString
	returnRDF := bytes.NewBuffer([]byte(constructedSparqlQuery))
	resp, err := http.Post("https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", "data-binary", returnRDF)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatal(err)
	}
	sb := string(body)
	return c.String(http.StatusOK, constructedSparqlQuery+": "+sb)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, you've reached the Go API that lets you talk to the Neptune database. Well done!")
	})

	e.POST("/sparql", requestToNeptune)

	e.Logger.Fatal(e.Start(":9000"))

}
