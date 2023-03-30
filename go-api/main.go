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

	_ "github.com/OurHeritageOurStories/ohos-neptune-ec2-api/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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

//struct for results for title/topic/url/description in movingImages
type movingImagesTitleTopicUrlDesc struct {
	Head    RDFHeadResponse
	Results TitleTopicUrlDescriptionBindingStruct
}

type TitleTopicUrlDescriptionBindingStruct struct {
	Bindings []BindingsTitleTopicUrlDescription
}

type BindingsTitleTopicUrlDescription struct {
	Identifier 	IdentifierReturnValues
	Title       TitleTopicStructValues
	Description TitleTopicStructValues
	URL         URLReturnValues
	Topics      TitleTopicStructValues
}

type URLReturnValues struct {
	Datatype string `json:"datatype"`
	Type     string `json:"type"`
	Value    string `json:"value"`
}

type IdentifierReturnValues struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
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
	Title  TitleTopicStructValues `json:"title"`
	Topics TitleTopicStructValues `json:"topics"`
}

type TitleTopicStructValues struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// struct for returning a keyword search

type keywordReturnStruct struct {
	Id    KeywordStruct                      `json:"id"`
	Total int                                `json:"total"`
	First string                             `json:"first"`
	Prev  string                             `json:"prev"`
	Next  string                             `json:"next"`
	Last  string                             `json:"last"`
	Items []BindingsTitleTopicUrlDescription `json:"items"`
}

type EntityReturnStruct struct {
	Items []BindingsTitleTopicUrlDescription `json:"items"`
}

type KeywordStruct struct {
	Page    string `json:"page"`
	Keyword string `json:"keyword"`
}

//Discovery held by all result
type discoveryAllStruct struct {
	CagalogueLevels       []DiscoveryCodeCount
	ClosureStatuses       []DiscoveryCodeCount
	Count                 int `json:"count"`
	Departments           []DiscoveryCodeCount
	HeldByReps            []DiscoveryCodeCount
	NextBatchMark         string `json:"nextBatchMark"`
	Records               []DiscoveryRecordDetails
	ReferenceFirstLetters []DiscoveryCodeCount
	Repositories          []DiscoveryCodeCount
	Sources               []DiscoveryCodeCount
	TaxonomySubject       []DiscoveryCodeCount
	TimePeriods           []DiscoveryCodeCount
	TitleFirstLetters     []DiscoveryCodeCount
}

type DiscoveryRecordDetails struct {
	AdminHistory       string   `json:"adminHistory"`
	AltName            string   `json:"altName"`
	Arrangement        string   `json:"arrangement"`
	CatalogueLevel     int      `json:"catalogueLevel"`
	ClosureCode        string   `json:"closureCode"`
	ClosureStatus      string   `json:"closureStatus"`
	ClosureType        string   `json:"closureType"`
	Content            string   `json:"content"`
	Context            string   `json:"context"`
	CorpBodies         []string `json:"corpBodies"`
	CoveringDates      string   `json:"coveringDates"`
	Department         string   `json:"department"`
	Description        string   `json:"description"`
	DocumentType       string   `json:"documentType"`
	EndDate            string   `json:"endDate"`
	FormerReferenceDep string   `json:"formerReferenceDep"`
	FormerReferencePro string   `json:"formerReferencePro"`
	HeldBy             []string `json:"heldBy"`
	Id                 string   `json:"id"`
	MapDesignation     string   `json:"mapDesignation"`
	MapScale           string   `json:"mapScale"`
	Note               string   `json:"note"`
	NumEndDate         int      `json:"numEndDate"`
	NumStartDate       int      `json:"numStartDate"`
	OpeningDate        string   `json:"openingDate"`
	PhysicalCondition  string   `json:"physicalCondition"`
	Places             []string `json:"places"`
	Reference          string   `json:"reference"`
	Score              int      `json:"score"`
	Source             string   `json:"source"`
	StartDate          string   `json:"startDate"`
	Taxonomies         []string `json:"taxonomies"`
	Title              string   `json:"title"`
	UrlParameters      string   `json:"urlParameters"`
}

type DiscoveryCodeCount struct {
	Code  string `json:"code"`
	Count int    `json:"count"`
}

// TODO this should have more than one return
func buildMainSparqlQuery(keyword string, offset string) string {
	titleTopicURLDescription := "prefix ns0: <http://id.loc.gov/ontologies/bibframe/> prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> prefix xsd: <http://www.w3.org/2001/XMLSchema#> prefix ns1: <http://id.loc.gov/ontologies/bflc/> prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>select (?id as ?identifier) ?title ?url ?description (group_concat(?topic;separator=' ||| ')as ?topics) where {?s ns0:title _:title ._:title ns0:mainTitle ?title filter (regex(str(?title), '" + keyword + "', 'i')) .?s ns0:summary _:summary ._:summary rdfs:label ?description .?s ns0:subject _:subject ._:subject rdfs:label ?topic .?s ns0:hasInstance ?t .?t ns0:hasItem ?r .?r ns0:electronicLocator _:url ._:url rdf:value ?url . ?s ns0:adminMetadata _:adminData ._:adminData ns0:identifiedBy _:identifiedBy ._:identifiedBy rdf:value ?id .} group by ?title ?id ?description ?url order by ?title  OFFSET " + offset + " limit 10"
	return titleTopicURLDescription
}

func buildEntityMainSparqlQuery(id string) string {
	titleTopicURLDescription := "prefix ns0: <http://id.loc.gov/ontologies/bibframe/> prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> prefix xsd: <http://www.w3.org/2001/XMLSchema#> prefix ns1: <http://id.loc.gov/ontologies/bflc/> prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> select ('"+id+"' as ?identifier) ?title ?url ?description (group_concat(?topic;separator=' ||| ')as ?topics) where {?s ns0:title _:title ._:title ns0:mainTitle ?title .?s ns0:summary _:summary ._:summary rdfs:label ?description .?s ns0:subject _:subject ._:subject rdfs:label ?topic .?s ns0:hasInstance ?t .?t ns0:hasItem ?r .?r ns0:electronicLocator _:url ._:url rdf:value ?url .?s ns0:adminMetadata _:adminData ._:adminData ns0:identifiedBy _:identifiedBy ._:identifiedBy rdf:value '"+id+"' .} group by ?title ?description ?url order by ?title OFFSET 0 limit 10"
	return titleTopicURLDescription
}

// StatusCheck godoc
// @Summary Test whether the API is running
// @Description Test whether the api is running
// @Tags root
// @Produce plain
// @Success 200
// @Router / [get]
func helloResponse(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, you've reached the Go API that lets you talk to the Neptune database. Well done!")
}

// NeptuneSparql godoc
// @Summary Send sparql direct to neptune
// @Description Send sparql direct to neptune
// @Tags Sparql
// @Accept json
// @Produce json
// @Success 200
// @Router /sparql [post]
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

	return c.JSONPretty(http.StatusOK, jsonMap, " ")
}

// Discovery godoc
// @Summary Requests to TNA Discovery API
// @Description Requests to TNA Discovery API
// @Tags Discovery
// @Produce json
// @Param q query string true "string query"
// @Param source query string true "string sourceArchives"
// @Success 200 {object} discoveryAllStruct
// @Router /discovery [get]
func fetchDiscovery(c echo.Context) error {
	keyword := c.Param("q")
	source := c.Param("source")

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

// Moving Images godoc
// @Summary Moving images queries
// @Description Moving images queries
// @Tags MovingImages
// @Param q query string true "string query"
// @Param page query int true "int page"
// @Produce json
// @Success 200 {object} keywordReturnStruct
// @Success 204
// @Failure 400
// @Failure 500
// @Router /movingImages [get]
func movingImages(c echo.Context) error {

	//default params
	keyword := "glasgow"
	pageKeyword := "1"
	pageInt := 1
	numberOfResults := 0

	var jsonToReturn keywordReturnStruct

	userProvidedParams := c.QueryParams()

	//check if we've got both

	if len(userProvidedParams) != 2 {
		return c.String(http.StatusBadRequest, "You need to provide both a keyword and a page number")
	} else {
		keyword = userProvidedParams.Get("q")
		pageKeyword = userProvidedParams.Get("page")
		jsonToReturn.Id.Keyword = keyword
		jsonToReturn.Id.Page = pageKeyword
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
	countParams.Add("format", "json")
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
	jsonToReturn.Total = numberOfResults
	jsonToReturn.First = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/mvingImages?q=" + keyword + "&page=1"
	if off == 1 {
		jsonToReturn.Prev = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/movingImages?q=" + keyword + "&page=1"
	} else {
		jsonToReturn.Prev = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/movingImages?q=" + keyword + "&page=" + strconv.Itoa(off-1)
	}
	if off == maxPages {
		jsonToReturn.Next = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/movingImages?q=" + keyword + "&page=" + strconv.Itoa(maxPages)
	} else {
		jsonToReturn.Next = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/movingImages?q=" + keyword + "&page=" + strconv.Itoa(off+1)
	}
	jsonToReturn.Last = "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/movingImages?q=" + keyword + "&page=" + strconv.Itoa(maxPages)

	defer countResp.Body.Close()

	//Now, we do the actual query

	mainSearchQuery := buildMainSparqlQuery(keyword, offset)

	mainSearchParams := url.Values{}
	mainSearchParams.Add("query", mainSearchQuery)
	mainSearchParams.Add("format", "json")
	mainSearchBody := strings.NewReader(mainSearchParams.Encode())

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

	var mainResultStruct movingImagesTitleTopicUrlDesc

	mainResultData, err := ioutil.ReadAll(mainSearchResp.Body)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong reading the main response.")
	}

	if err := json.Unmarshal(mainResultData, &mainResultStruct); err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong working with the main results.")
	}

	jsonToReturn.Items = mainResultStruct.Results.Bindings

	return c.JSONPretty(http.StatusOK, jsonToReturn, " ")
}

// Moving Images Entity godoc
// @Summary Moving images get specific entity query
// @Description Moving images get specific entity query
// @Tags MovingImages Entity
// @Param id path string true "string id"
// @Produce json
// @Success 200 {object} EntityReturnStruct
// @Failure 500
// @Router /movingImagesEnt/entity [get]
func movingImagesEntity(c echo.Context) error {

	var jsonToReturn EntityReturnStruct

	id := c.Param("id")

	mainSearchQuery := buildEntityMainSparqlQuery(id)

	mainSearchParams := url.Values{}
	mainSearchParams.Add("query", mainSearchQuery)
	mainSearchParams.Add("format", "json")
	mainSearchBody := strings.NewReader(mainSearchParams.Encode())

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

	var mainResultStruct movingImagesTitleTopicUrlDesc

	mainResultData, err := ioutil.ReadAll(mainSearchResp.Body)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong reading the main response.")
	}

	if err := json.Unmarshal(mainResultData, &mainResultStruct); err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong working with the main results.")
	}

	jsonToReturn.Items = mainResultStruct.Results.Bindings

	return c.JSONPretty(http.StatusOK, jsonToReturn, " ")
}

// @title OHOS api
// @version 1.0.1
// @description OHOS api
// @termsOfService http://swagger.io/terms/
// @contact.name The National Archives
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000
// @BasePath /api
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

	e.GET("/movingImages", movingImages)

	e.GET("/movingImagesEnt/entity/:id", movingImagesEntity)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":9000"))

}
