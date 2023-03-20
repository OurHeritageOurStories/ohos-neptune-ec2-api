package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
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

type RDFHeadResponse struct {
	Vars []string
}

// full struct for count
type resultsCountStruct struct {
	Head    RDFHeadResponse
	Results ResultsCount
}

type ResultsCount struct {
	Bindings []ResultsBindings
}

type ResultsBindings struct {
	Count resultsBindingCount
}

type resultsBindingCount struct {
	Datatype string `json: datatype`
	Type     string `json: type`
	Value    string `json: value`
}

// struct for results for title/topic search

type TitleTopicStruct struct {
	Head    RDFHeadResponse
	Results TitleTopicBindingsStruct
}

type TitleTopicBindingsStruct struct {
	Bindings []BindingsResultsTitleTopic
}

type BindingsResultsTitleTopic struct {
	Title struct {
		Type  string `json: type`
		Value string `json: value`
	}
	Topic struct {
		Type  string `json: type`
		Value string `json: value`
	}
}

/*type BindingsResultsTitleTopic struct {
	Title  TitleTopicStructValues
	Topics TitleTopicStructValues
}

type TitleTopicStructValues struct {
	Type  string `json: type`
	Value string `json: value`
}*/

// struct for returning a keyword search

type keywordReturnStruct struct {
	Keywords     KeywordStruct
	TotalCount   int    `json: totalResults`
	FirstPage    string `json: firstPage`
	PreviousPage string `json: previousPage`
	CurrentPage  string `json: currentPage`
	NextPage     string `json: nextPage`
	LastPage     string `json: lastPage`
	Results      []TitleTopicBindingsStruct
}

type KeywordStruct struct {
	Page    string `json: page`
	Keyword string `json: keyword`
}

/*
func getResponseFromNeptune(query string, target interface{}) error {
	params := url.Values{}
	params.Add("query", query)
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

	return json.NewDecoder(req.Body).Decode(target)
}
*/

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

	providedParams := c.QueryParams()

	//default values
	keyword := "glasgow"
	pageKeyword := "1"
	pageInt := 1

	if len(providedParams) != 2 {
		return c.String(http.StatusUnprocessableEntity, "Missing a required param")
	} else {
		keyword = providedParams.Get("keyword")
		pageKeyword = providedParams.Get("page")
	}

	pageInt, err := strconv.Atoi(pageKeyword)
	if err != nil {
		return c.String(http.StatusBadRequest, "Page numbers need to be a number")
	}

	if pageInt < 1 {
		return c.String(http.StatusRequestedRangeNotSatisfiable, "If you are looking for the first page, please request page 1 (this is to align with the UI). If you are looking for page -1 or lower, those don't exist.")
	}

	off := max(1, pageInt)

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

	//testjson := new(resultsCountStruct)

	defer respCount.Body.Close()

	dataCount, _ := ioutil.ReadAll(respCount.Body)
	var jsonMapCount map[string]interface{}
	json.Unmarshal([]byte(dataCount), &jsonMapCount)

	jsonMap["count"] = jsonMapCount
	//countResultJson := resultsCountStruct{}
	//json.Marshal([]byte(dataCount), countResultJson)

	//stick the total number of results to the list

	//jsonMap["TotalCount"] = countResultJson.Results.Bindings[0].Count.Value
	return c.JSONPretty(http.StatusOK, jsonMap, " ")
}

func helloResponse(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, you've reached the Go API that lets you talk to the Neptune database. Well done!")
}

func neptest(c echo.Context) error {
	keyword := c.QueryParam("keyword")
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

	//testjson := new(resultsCountStruct)

	var structCount resultsCountStruct
	dataCount, _ := ioutil.ReadAll(respCount.Body)
	if err := json.Unmarshal(dataCount, &structCount); err != nil {
		fmt.Print("no")
	}

	defer respCount.Body.Close()
	return c.JSONPretty(http.StatusOK, structCount.Results.Bindings[0].Count.Value, " ")
	/*esultTest := new(resultsCountStruct)
	getResponseFromNeptune(constructedQueryCount, resultTest)
	return c.String(http.StatusTeapot, resultTest.Results.Bindings[0].Count.Value)
	//return c.String(http.StatusOK, "angry inch")
	//return c.String(http.StatusTeapot, resultTest.Results.Bindings[0].Count.Value)*/
}

func movingImagesBetter(c echo.Context) error {

	//default params
	keyword := "glasgow"
	pageKeyword := "1"
	pageInt := 1
	numberOfResults := 0
	// maxPages := 1

	var jsonToReturn keywordReturnStruct

	userProvidedParams := c.QueryParams()

	//check if we've got both

	if len(userProvidedParams) != 2 {
		return c.String(http.StatusUnprocessableEntity, "You need to provide both a keyword and a page number")
	} else {
		keyword = userProvidedParams.Get("keyword")
		pageKeyword = userProvidedParams.Get("page")
		jsonToReturn.Keywords.Keyword = keyword
		jsonToReturn.Keywords.Page = pageKeyword
	}

	pageInt, err := strconv.Atoi(pageKeyword)
	if err != nil {
		return c.String(http.StatusBadRequest, "Page needs to be selected as a number")
	}

	if pageInt < 1 {
		return c.String(http.StatusRequestedRangeNotSatisfiable, "If you are looking for the first page, please request page 1 (this is to align with the UI). If you are looking for page -1 or lower, those don't exist.")
	}

	off := max(1, pageInt)

	offset := strconv.Itoa((off - 1) * 10)

	//check if there are any actual results

	countQuery := "prefix ns0: <http://id.loc.gov/ontologies/bibframe/> prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> prefix xsd: <http://www.w3.org/2001/XMLSchema#> prefix ns1: <http://id.loc.gov/ontologies/bflc/> prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> select (count(*) as ?count) where {select ?title (group_concat(?topic;separator=' ||| ')as ?topics) where {?s ns0:title _:title ._:title ns0:mainTitle ?title filter (regex(str(?title), '" + keyword + "', 'i')) .?s ns0:subject _:subject ._:subject rdfs:label ?topic .} group by ?title order by desc(?count)}"

	countParams := url.Values{}
	countParams.Add("query", countQuery)
	countBody := strings.NewReader(countParams.Encode())

	countReq, err := http.NewRequest("POST", "https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", countBody)

	if err != nil {
		log.Fatal(err)
		return c.String(http.StatusInternalServerError, "Something went wrong sending the request to the database.")
	}

	countReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	countResp, err := http.DefaultClient.Do(countReq)

	if err != nil {
		log.Fatal(err)
		return c.String(http.StatusInternalServerError, "Something went wrong getting the result from the database.")
	}

	// map the response to the "count" struct

	var countStruct resultsCountStruct

	countData, err := ioutil.ReadAll(countResp.Body)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong reading the number of responses.")
	}

	if err := json.Unmarshal(countData, &countStruct); err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong working with the number of responses")
	}

	numberOfResults, err = strconv.Atoi(countStruct.Results.Bindings[0].Count.Value)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong as the number of results isn't a number")
	}

	if numberOfResults < 1 {
		return c.String(http.StatusNoContent, "The search worked, there just aren't any results") // No Content doesn't send a return string, which makes sense, as that'd be some content
	}

	maxPages := int(math.Ceil(float64(numberOfResults) / 10)) //odd way to do it, but that's what the linter did to it

	//Now that we know the number of pages, we can fill in the various page options
	jsonToReturn.TotalCount = numberOfResults
	jsonToReturn.FirstPage = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/betterMovingImages?keyword=" + keyword + "&page=1"
	if off == 1 {
		jsonToReturn.PreviousPage = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/betterMovingImages?keyword=" + keyword + "&page=1"
	} else {
		jsonToReturn.PreviousPage = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/betterMovingImages?keyword=" + keyword + "&page=" + strconv.Itoa(off-1)
	}
	jsonToReturn.CurrentPage = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/betterMovingImages?keyword=" + keyword + "&page=" + pageKeyword
	if off == maxPages {
		jsonToReturn.NextPage = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/betterMovingImages?keyword=" + keyword + "&page=" + strconv.Itoa(maxPages)
	} else {
		jsonToReturn.NextPage = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/betterMovingImages?keyword=" + keyword + "&page=" + strconv.Itoa(off+1)
	}
	jsonToReturn.LastPage = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/betterMovingImages?keyword=" + keyword + "&page=" + strconv.Itoa(maxPages)

	defer countResp.Body.Close()

	//Now, we do the actual query

	mainSearchQuery := "prefix ns0: <http://id.loc.gov/ontologies/bibframe/> prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> prefix xsd: <http://www.w3.org/2001/XMLSchema#> prefix ns1: <http://id.loc.gov/ontologies/bflc/> prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> select ?title (group_concat(?topic;separator=' ||| ')as ?topics) where {?s ns0:title _:title ._:title ns0:mainTitle ?title filter (regex(str(?title), '" + keyword + "', 'i')) .?s ns0:subject _:subject ._:subject rdfs:label ?topic .} group by ?title order by ?title OFFSET " + offset + " LIMIT 10"

	mainSearchParams := url.Values{}
	mainSearchParams.Add("query", mainSearchQuery)
	mainSearchBody := strings.NewReader(countParams.Encode())

	mainSearchReq, err := http.NewRequest("POST", "https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", mainSearchBody)

	if err != nil {
		log.Fatal(err)
		return c.String(http.StatusInternalServerError, "Something went wrong sending the main request to the database.")
	}

	mainSearchReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mainSearchResp, err := http.DefaultClient.Do(mainSearchReq)

	if err != nil {
		log.Fatal(err)
		return c.String(http.StatusInternalServerError, "Something went wrong getting the main result from the database.")
	}

	// Map the main response to the struct

	var mainResultStruct TitleTopicStruct

	mainResultData, err := ioutil.ReadAll(mainSearchResp.Body)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong reading the main response.")
	}

	if err := json.Unmarshal(mainResultData, &mainResultStruct); err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong working with the main results.")
	}

	jsonToReturn.Results = mainResultStruct.Results

	return c.JSONPretty(http.StatusOK, mainResultStruct, " ")
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

	e.GET("/neptest", neptest)

	e.GET("/betterMovingImages", movingImagesBetter)

	e.Logger.Fatal(e.Start(":9000"))

}
