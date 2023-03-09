package main

import (
	"encoding/json"
	"fmt"
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func requestToNeptune(c echo.Context) error {

	sparqlString := c.FormValue("sparqlquery")

	limit := c.FormValue("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return c.String(http.StatusBadRequest, "Limit needs to be a number")
	}

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

	params := url.Values{}
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

	var jsonMap map[string]interface{}

	json.Unmarshal([]byte(data), &jsonMap)

	return c.JSON(http.StatusOK, jsonMap)
}

func fetchDiscovery(c echo.Context) error {
	keyword := c.Request().URL.Query().Get("keyword")
	source := strings.ToUpper(c.Request().URL.Query().Get("source"))

	if source == "" {
		source = "ALL"
	}

	response, err := http.Get("https://discovery.nationalarchives.gov.uk/API/search/records?sps.heldByCode=" + source + "&sps.searchQuery=" + keyword)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonMap map[string]interface{}

	json.Unmarshal([]byte(responseData), &jsonMap)

	return c.JSONPretty(http.StatusOK, jsonMap, " ")
}

/*
func getEntities(c echo.Context) error {

	keyword := c.QueryParam("keyword")
	page := c.QueryParam("page")

	off, err := strconv.Atoi(page)

	if err != nil {
		log.Fatal(err)
	}

	off = max(1, off)

	o := strconv.Itoa((off - 1) * 10)

	constructedQuery := "prefix tanc: <http://tanc.manchester.ac.uk/> SELECT DISTINCT ?o (count(?text) as ?count) WHERE { ?s <http://tanc.manchester.ac.uk/text> ?text. ?s tanc:mentions ?o FILTER (regex(str(?o), '" + keyword + "', 'i'))} GROUP BY ?o ORDER BY DESC(?count) OFFSET " + o + " LIMIT 10"

	// updated curl -X POST --data-binary 'query= prefix schema: <https://schema.org/> prefix ver:   <http://purl.org/linked-data/version#> prefix xsd:   <http://www.w3.org/2001/XMLSchema#> prefix rdfs:  <http://www.w3.org/2000/01/rdf-schema#> prefix edm:   <http://www.europeana.eu/schemas/edm/> prefix rst:   <http://id.loc.gov/vocabulary/relationshipSubType/> prefix rdau:  <http://rdaregistry.info/Elements/u/> prefix dct:   <http://purl.org/dc/terms/> prefix rdf:   <http://www.w3.org/1999/02/22-rdf-syntax-ns#> prefix cat:   <http://cat.nationalarchives.gov.uk/> prefix time:  <http://www.w3.org/2006/time#> prefix odrl:  <http://www.w3.org/ns/odrl/2/> prefix prov:  <http://www.w3.org/ns/prov#> prefix iso639-2: <http://id.loc.gov/vocabulary/iso639-2/> select distinct ?o (count(?text) as ?count) where {?s dct:abstract ?text .?s dct:abstract ?o filter (regex(str(?o), "part", "i"))}group by ?o order by desc(?count) offset 3 limit 10' https://ohos-live-data-neptune.cluster-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql

	//constructedQuery := "prefix schema: <https://schema.org/> prefix ver:   <http://purl.org/linked-data/version#> prefix xsd:   <http://www.w3.org/2001/XMLSchema#> prefix rdfs:  <http://www.w3.org/2000/01/rdf-schema#> prefix edm:   <http://www.europeana.eu/schemas/edm/> prefix rst:   <http://id.loc.gov/vocabulary/relationshipSubType/> prefix rdau:  <http://rdaregistry.info/Elements/u/> prefix dct:   <http://purl.org/dc/terms/> prefix rdf:   <http://www.w3.org/1999/02/22-rdf-syntax-ns#> prefix cat:   <http://cat.nationalarchives.gov.uk/> prefix time:  <http://www.w3.org/2006/time#> prefix odrl:  <http://www.w3.org/ns/odrl/2/> prefix prov:  <http://www.w3.org/ns/prov#> prefix iso639-2: <http://id.loc.gov/vocabulary/iso639-2/> select distinct ?o (count(?text) as ?count) where {?s dct:abstract ?text .?s dct:abstract ?o filter (regex(str(?o), " + keyword + ", 'i'))}group by ?o order by desc(?count) offset " + o + " limit 10"

	params := url.Values{}
	params.Add("query", constructedQuery)
	body := strings.NewReader(params.Encode())

	response, err := http.NewRequest("POST", "https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", body)

	if err != nil {
		fmt.Print(err.Error())
		fmt.Print("aaaahhhh")
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	//constructedQueryTwo := "prefix tanc: <http://tanc.manchester.ac.uk/> SELECT (count(*) as ?count) WHERE {SELECT DISTINCT ?o (count(?text) as ?count) WHERE { ?s <http://tanc.manchester.ac.uk/text> ?text. ?s tanc:mentions ?o FILTER (regex(str(?o), '" + keyword + "', 'i'))} GROUP BY ?o ORDER BY DESC(?count)}"

	//curl -X POST --data-binary 'query= prefix schema: <https://schema.org/> prefix ver:   <http://purl.org/linked-data/version#> prefix xsd:   <http://www.w3.org/2001/XMLSchema#> prefix rdfs:  <http://www.w3.org/2000/01/rdf-schema#> prefix edm:   <http://www.europeana.eu/schemas/edm/> prefix rst:   <http://id.loc.gov/vocabulary/relationshipSubType/> prefix rdau:  <http://rdaregistry.info/Elements/u/> prefix dct:   <http://purl.org/dc/terms/> prefix rdf:   <http://www.w3.org/1999/02/22-rdf-syntax-ns#> prefix cat:   <http://cat.nationalarchives.gov.uk/> prefix time:  <http://www.w3.org/2006/time#> prefix odrl:  <http://www.w3.org/ns/odrl/2/> prefix prov:  <http://www.w3.org/ns/prov#> prefix iso639-2: <http://id.loc.gov/vocabulary/iso639-2/> SELECT (count(*) as ?count) WHERE {SELECT DISTINCT ?o (count(?text) as ?count) WHERE { ?s dct:abstract ?text. ?s dct:abstract ?o FILTER (regex(str(?o), "part", "i"))} GROUP BY ?o ORDER BY DESC(?count)}' https://ohos-live-data-neptune.cluster-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql

	//constructedQueryTwo := "prefix schema: <https://schema.org/> prefix ver:   <http://purl.org/linked-data/version#> prefix xsd:   <http://www.w3.org/2001/XMLSchema#> prefix rdfs:  <http://www.w3.org/2000/01/rdf-schema#> prefix edm:   <http://www.europeana.eu/schemas/edm/> prefix rst:   <http://id.loc.gov/vocabulary/relationshipSubType/> prefix rdau:  <http://rdaregistry.info/Elements/u/> prefix dct:   <http://purl.org/dc/terms/> prefix rdf:   <http://www.w3.org/1999/02/22-rdf-syntax-ns#> prefix cat:   <http://cat.nationalarchives.gov.uk/> prefix time:  <http://www.w3.org/2006/time#> prefix odrl:  <http://www.w3.org/ns/odrl/2/> prefix prov:  <http://www.w3.org/ns/prov#> prefix iso639-2: <http://id.loc.gov/vocabulary/iso639-2/> SELECT (count(*) as ?count) WHERE {SELECT DISTINCT ?o (count(?text) as ?count) WHERE { ?s dct:abstract ?text. ?s dct:abstract ?o FILTER (regex(str(?o), " + keyword + ", 'i'))} GROUP BY ?o ORDER BY DESC(?count)}"

	params2 := url.Values{}
	params2.Add("query", constructedQueryTwo)
	body2 := strings.NewReader(params2.Encode())

	response2, err := http.NewRequest("POST", "https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", body2)

	if err != nil {
		fmt.Print(err.Error())
		fmt.Print("eeeek")
		os.Exit(1)
	}

	responseData2, err := ioutil.ReadAll(response2.Body)

	if err != nil {
		log.Fatal(err)
	}

	out := map[string]interface{}{}
	json.Unmarshal([]byte(responseData), &out)

	out2 := map[string]interface{}{}
	json.Unmarshal([]byte(responseData2), &out2)

	out["count"] = out2

	outputJSON, _ := json.Marshal(out)

	var jsonMap map[string]interface{}

	json.Unmarshal([]byte(outputJSON), &jsonMap)

	return c.JSON(http.StatusOK, jsonMap)
}
*/
func getEntity(c echo.Context) error {
	entity := c.Request().URL.Query().Get("entity")
	//page := strings.ToUpper(c.Request().URL.Query().Get("page"))
	page := c.Request().URL.Query().Get("page")

	fmt.Print("entity" + entity)
	fmt.Print("page" + page)

	off, err := strconv.Atoi(page)

	off = max(1, off)

	o := strconv.Itoa((off - 1) * 10)

	constructedQuery := "SELECT DISTINCT ?text (group_concat(?mentioned;separator=' ') as ?m)  WHERE { ?s <http://tanc.manchester.ac.uk/mentions> ?mentioned. ?s <http://tanc.manchester.ac.uk/text> ?text. ?s <http://tanc.manchester.ac.uk/mentions> <https://en.wikipedia.org/wiki/" + entity + ">.} GROUP BY ?text ORDER BY ?text OFFSET " + o + "LIMIT 10"

	params := url.Values{}
	params.Add("query", constructedQuery)
	body := strings.NewReader(params.Encode())

	response, err := http.NewRequest("POST", "https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", body)

	if err != nil {
		fmt.Print(err.Error())
		fmt.Print("erm")
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	constructedQueryTwo := "SELECT (count(?o) as ?count) WHERE { ?s <http://tanc.manchester.ac.uk/text> ?o. ?s <http://tanc.manchester.ac.uk/mentions> <https://en.wikipedia.org/wiki/" + entity + ">.}"

	params2 := url.Values{}
	params2.Add("query", constructedQueryTwo)
	body2 := strings.NewReader(params.Encode())

	response2, err := http.NewRequest("POST", "https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", body2)

	if err != nil {
		fmt.Print(err.Error())
		fmt.Print("nope")
		os.Exit(1)
	}

	responseData2, err := ioutil.ReadAll(response2.Body)

	out := map[string]interface{}{}
	json.Unmarshal([]byte(responseData), &out)

	out2 := map[string]interface{}{}
	json.Unmarshal([]byte(responseData2), &out2)

	out["count"] = out2

	outputJSON, _ := json.Marshal(out)

	var jsonMap map[string]interface{}

	json.Unmarshal([]byte(outputJSON), &jsonMap)

	return c.JSON(http.StatusOK, jsonMap)

}

func movingImages(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	page := c.QueryParam("page")

	off, err := strconv.Atoi(page)

	if err != nil {
		log.Fatal(err)
	}

	off = max(1, off)

	o := strconv.Itoa((off - 1) * 10)

	constructedQuery := "prefix ns0: <http://id.loc.gov/ontologies/bibframe/> prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> prefix xsd: <http://www.w3.org/2001/XMLSchema#> prefix ns1: <http://id.loc.gov/ontologies/bflc/> prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> select ?title (group_concat(?topic;separator=' ||| ')as ?topics) where {?s ns0:title _:title ._:title ns0:mainTitle ?title filter (regex(str(?title), '" + keyword + "', 'i')) .?s ns0:subject _:subject ._:subject rdfs:label ?topic .} group by ?title order by ?title OFFSET " + o + " LIMIT 10"
	return c.String(http.StatusTeapot, constructedQuery)
	/*
		params := url.Values{}
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

		var jsonMap map[string]interface{}

		json.Unmarshal([]byte(data), &jsonMap)

		return c.JSONPretty(http.StatusOK, jsonMap, " ")*/
}

/*
	response, err := http.NewRequest("POST", "https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", body)

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	returnString := string(responseData)

	return c.String(http.StatusTeapot, returnString)

	//var jsonMap map[string]interface{}

	//json.Unmarshal(responseData, &jsonMap)

	//return c.JSON(http.StatusOK, jsonMap)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//respneded, err := json.Marshal(responseData)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//return c.JSON(http.StatusOK, respneded)

	//var jsonMap map[string]interface{}

	//json.Unmarshal([]byte(responseData), &jsonMap)

	//return c.Stringhttp.StatusOK, responseData)
	//return c.JSONPretty(http.StatusOK, jsonMap, " ")
}*/

func entitytest(c echo.Context) error {
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error { //make sure its alive
		return c.HTML(http.StatusOK, "Hello, you've reached the Go API that lets you talk to the Neptune database. Well done!")
	})

	e.POST("/sparql", requestToNeptune) //for actual requests

	e.GET("/discovery", fetchDiscovery)

	//e.GET("/entities", getEntities)

	e.GET("/testentities", entitytest)

	e.GET("/movingImages", movingImages)

	//e.GET("/entities/{entity}", getEntity)

	e.Logger.Fatal(e.Start(":9000"))

}
