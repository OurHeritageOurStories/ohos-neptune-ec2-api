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
Methods unused at this stage
*/
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
func getEntity(c echo.Context) error {
	entity := c.QueryParam("entity")
	page := c.QueryParam("page")

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
*/

func movingImages(c echo.Context) error {
	// keyword := c.QueryParam("keyword")
	//page := c.QueryParam("page")
	// page := c.Request().URL.Query().Get("page")
	// page, err := c.Request().URL.Query().Get("page")
	//if len(c.Echo().URL().)
	//if len(self.Request().GET.keys()!=2){
	//	return c.NoContent(http.StatusTeapot)
	//}

	//c.ParseForm()

	//_, hasPage := c.Form["page"]
	/*
		params := c.QueryParams()

		numberOfParams := len(params)

		if numberOfParams != 2 {
			return c.NoContent(http.StatusTeapot)
		} else {
			page := params.Get("page")
		}

		if keyword == "" && page == "" {
			return c.NoContent(http.StatusTeapot)
		}*/

	providedParams := c.QueryParams()

	//default values
	keyword := "glasgow"
	page := "1"
	missingParams := false

	// This is a fairly blunt-instrument approach to dealing with missing
	// params. Once we have more than 2 and when they are optional, we should
	// switch this over to something else - maybe a switch, maybe something
	// else
	if len(providedParams) != 2 {
		missingParams = true
	} else {
		keyword = providedParams.Get("keyword")
		page = providedParams.Get("page")
	}

	off, err := strconv.Atoi(page)

	if err != nil {
		log.Fatal(err)
	}

	off = max(1, off)

	o := strconv.Itoa((off - 1) * 10)

	//get the page of results

	constructedQuery := "prefix ns0: <http://id.loc.gov/ontologies/bibframe/> prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> prefix xsd: <http://www.w3.org/2001/XMLSchema#> prefix ns1: <http://id.loc.gov/ontologies/bflc/> prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> select ?title (group_concat(?topic;separator=' ||| ')as ?topics) where {?s ns0:title _:title ._:title ns0:mainTitle ?title filter (regex(str(?title), '" + keyword + "', 'i')) .?s ns0:subject _:subject ._:subject rdfs:label ?topic .} group by ?title order by ?title OFFSET " + o + " LIMIT 10"
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

	//get the total number of results

	constructedQueryCount := "prefix ns0: <http://id.loc.gov/ontologies/bibframe/> prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> prefix xsd: <http://www.w3.org/2001/XMLSchema#> prefix ns1: <http://id.loc.gov/ontologies/bflc/> prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> select (count(*) as ?count) where {select ?title (group_concat(?topic;separator=' ||| ')as ?topics) where {?s ns0:title _:title ._:title ns0:mainTitle ?title filter (regex(str(?title), '" + keyword + "', 'i')) .?s ns0:subject _:subject ._:subject rdfs:label ?topic .} group by ?title order by desc(?count)}"
	paramsCount := url.Values{}
	paramsCount.Add("query", constructedQueryCount)
	bodyCount := strings.NewReader(paramsCount.Encode())

	reqCount, err := http.NewRequest("POST", "https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", bodyCount)
	if err != nil {
		log.Fatal(err)
	}
	reqCount.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	respCount, err := http.DefaultClient.Do(reqCount)
	if err != nil {
		log.Fatal(err)
	}

	defer respCount.Body.Close()

	dataCount, _ := ioutil.ReadAll(respCount.Body)
	var jsonMapCount map[string]interface{}
	json.Unmarshal([]byte(dataCount), &jsonMapCount)

	//stick the total number of results to the list

	jsonMap["count"] = jsonMapCount

	if missingParams {
		jsonMap["params"] = "Missing params, using default of keyword=glasgow, page=1"
	}

	return c.JSONPretty(http.StatusOK, jsonMap, " ")
}

func helloResponse(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, you've reached the Go API that lets you talk to the Neptune database. Well done!")
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.GET("/", helloResponse)

	e.POST("/sparql", requestToNeptune) //to pass requests directly through

	e.GET("/discovery", fetchDiscovery)

	//e.GET("/entities", getEntities)

	//e.GET("/entities/{entity}", getEntity)

	e.GET("/movingImages", movingImages)

	e.Logger.Fatal(e.Start(":9000"))

}
