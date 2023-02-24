package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func requestToNeptune(c echo.Context) error {

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	// curl -X POST --data-binary 'query=select ?s ?p ?o where {?s ?p ?o} limit 10' https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql

	sparqlString := c.FormValue("sparqlstring")
	params := url.Values{}
	//params.Add("query", `select ?s ?p ?o where {?s ?p ?o} limit 10`)
	params.Add("query", "select "+sparqlString)
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
	//thing := fmt.Sprint("resp.Body: %v\n", resp.Body)
	returned := fmt.Sprint(resp.Body)
	return c.String(http.StatusOK, returned)

	//sparqlString := c.FormValue("sqarqlstring")
	//sparqlString := "hi"
	//return c.String(http.StatusOK, sparqlString)
	//sparqlString := c.FormValue("sparqlstring")
	//constructedSparqlQuery := "query=select " + sparqlString
	//returnRDF := bytes.NewBuffer([]byte(constructedSparqlQuery))
	//resp, err := http.Post("https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", "data-binary", returnRDF)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()
	//body, err2 := ioutil.ReadAll(resp.Body)
	//if err2 != nil {
	//	log.Fatal(err)
	//}
	//sb := string(body)
	//return c.String(http.StatusOK, constructedSparqlQuery)
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
