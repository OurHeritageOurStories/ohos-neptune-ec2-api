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

	echoSwagger "github.com/AndrewBewseyTNA/echo-swagger"
	"github.com/AndrewBewseyTNA/echo/v4"
	"github.com/AndrewBewseyTNA/echo/v4/middleware"
	_ "github.com/OurHeritageOurStories/ohos-neptune-ec2-api/docs"
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
	Count DatatypeTypeValue
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
	Identifier  TypeValeStruct
	Title       TypeValeStruct
	Description TypeValeStruct
	URL         DatatypeTypeValue
	Topics      TypeValeStruct
}

type DatatypeTypeValue struct {
	Datatype string `json:"datatype"`
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
	Title  TypeValeStruct `json:"title"`
	Topics TypeValeStruct `json:"topics"`
}

type TypeValeStruct struct {
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
	Keyword string `json:"q"`
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
func buildMainSparqlQuery(keyword string, offset string, env string) string {
	titleTopicURLDescription := "prefix ns0: <http://id.loc.gov/ontologies/bibframe/> prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> prefix xsd: <http://www.w3.org/2001/XMLSchema#> prefix ns1: <http://id.loc.gov/ontologies/bflc/> prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>select (?id as ?identifier) ?title ?url ?description (group_concat(?topic;separator=' ||| ')as ?topics) where {?s ns0:title _:title ._:title ns0:mainTitle ?title filter (regex(str(?title), '" + keyword + "', 'i')) .?s ns0:summary _:summary ._:summary rdfs:label ?description .?s ns0:subject _:subject ._:subject rdfs:label ?topic .?s ns0:hasInstance ?t .?t ns0:hasItem ?r .?r ns0:electronicLocator _:url ._:url rdf:value ?url . ?s ns0:adminMetadata _:adminData ._:adminData ns0:identifiedBy _:identifiedBy ._:identifiedBy rdf:value ?id .} group by ?title ?id ?description ?url order by ?title  OFFSET " + offset + " limit 10"
	return titleTopicURLDescription
}

func buildEntityMainSparqlQuery(id string, env string) string {
	titleTopicURLDescription := "prefix ns0: <http://id.loc.gov/ontologies/bibframe/> prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> prefix xsd: <http://www.w3.org/2001/XMLSchema#> prefix ns1: <http://id.loc.gov/ontologies/bflc/> prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> select ('" + id + "' as ?identifier) ?title ?url ?description (group_concat(?topic;separator=' ||| ')as ?topics) where {?s ns0:title _:title ._:title ns0:mainTitle ?title .?s ns0:summary _:summary ._:summary rdfs:label ?description .?s ns0:subject _:subject ._:subject rdfs:label ?topic .?s ns0:hasInstance ?t .?t ns0:hasItem ?r .?r ns0:electronicLocator _:url ._:url rdf:value ?url .?s ns0:adminMetadata _:adminData ._:adminData ns0:identifiedBy _:identifiedBy ._:identifiedBy rdf:value '" + id + "' .} group by ?title ?description ?url order by ?title OFFSET 0 limit 10"
	return titleTopicURLDescription
}

// StatusCheck godoc
// @Summary Test whether the API is running
// @Description Test whether the api is running
// @Tags root
// @Produce plain
// @Success 200
// @Router / [get]
func helloResponse(welcome string) echo.HandlerFunc {
	fn := func(c echo.Context) error {
		return c.String(http.StatusOK, welcome)
	}
	return echo.HandlerFunc(fn)
}

// NeptuneSparql godoc
// @Summary Send sparql direct to neptune
// @Description Send sparql direct to neptune
// @Tags Sparql
// @Accept json
// @Produce json
// @Success 200
// @Router /sparql [post]
func requestToNeptune(neptuneurl, graph string) echo.HandlerFunc {
	fn := func(c echo.Context) error {
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

		req, err := http.NewRequest("POST", neptuneurl, body)
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
	return echo.HandlerFunc(fn)
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
func fetchDiscovery(discoveryapiurl string) echo.HandlerFunc {
	fn := func(c echo.Context) error {
		userProvidedParams := c.QueryParams()
		keyword := strings.Replace(userProvidedParams.Get("q"), " ", "%20", 1)
		source := strings.ToUpper(userProvidedParams.Get("source"))

		if source == "" {
			source = "ALL"
		}

		response, err := http.Get(discoveryapiurl + "records?sps.heldByCode=" + source + "&sps.searchQuery=" + keyword)

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
	return echo.HandlerFunc(fn)
}

// Moving Images godoc
// @Summary Moving images queries
// @Description Moving images queries
// @Tags MovingImages
// @Param q query string false "string query"
// @Param page query int false "int page"
// @Produce json
// @Success 200 {object} keywordReturnStruct
// @Failure 500
// @Router /moving-images [get]
func movingImages(ec2url, neptuneurl, movingImagesEndpoint, graph string) echo.HandlerFunc {
	fn := func(c echo.Context) error {
		//default params
		keyword := ""
		pageKeyword := "1"
		pageInt := 1
		numberOfResults := 0

		var jsonToReturn keywordReturnStruct

		userProvidedParams := c.QueryParams()

		//check if we've got both

		_, qPresent := userProvidedParams["q"]
		_, pagePresent := userProvidedParams["page"]

		if qPresent {
			keyword = userProvidedParams.Get("q")
		}

		if pagePresent {
			pageKeyword = userProvidedParams.Get("page")
		}

		jsonToReturn.Id.Keyword = keyword
		jsonToReturn.Id.Page = pageKeyword
	
		pageInt, err := strconv.Atoi(pageKeyword)
		if err != nil {
			return c.String(http.StatusBadRequest, "Page needs to be selected as a number")
		}

		if pageInt < 1 {
			pageInt = 1
		}

		off := max(1, pageInt)

		offset := strconv.Itoa((off - 1) * 10)

		//check if there are any actual results

		countQuery := "prefix ns0: <http://id.loc.gov/ontologies/bibframe/> prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> prefix xsd: <http://www.w3.org/2001/XMLSchema#> prefix ns1: <http://id.loc.gov/ontologies/bflc/> prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> select (count(*) as ?count) where {select ?title (group_concat(?topic;separator=' ||| ')as ?topics) where {?s ns0:title _:title ._:title ns0:mainTitle ?title filter (regex(str(?title), '" + keyword + "', 'i')) .?s ns0:subject _:subject ._:subject rdfs:label ?topic .} group by ?title order by desc(?count)}"

		countParams := url.Values{}
		countParams.Add("query", countQuery)
		countParams.Add("format", "json")
		countBody := strings.NewReader(countParams.Encode())

		countReq, err := http.NewRequest("POST", neptuneurl, countBody)

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

		maxPages := int(math.Ceil(float64(numberOfResults) / 10)) //odd way to do it, but that's what the linter did to it

		//Now that we know the number of pages, we can fill in the various page options
		jsonToReturn.Total = numberOfResults
		jsonToReturn.First = ec2url + "/api/" + movingImagesEndpoint + "?q=" + keyword + "&page=1"
		if off == 1 {
			jsonToReturn.Prev = ec2url + "/api/" + movingImagesEndpoint + "?q=" + keyword + "&page=1"
		} else {
			jsonToReturn.Prev = ec2url + "/api/" + movingImagesEndpoint + "?q=" + keyword + "&page=" + strconv.Itoa(off-1)
		}
		if off == maxPages {
			jsonToReturn.Next = ec2url + "/api/" + movingImagesEndpoint + "?q=" + keyword + "&page=" + strconv.Itoa(maxPages)
		} else {
			jsonToReturn.Next = ec2url + "/api/" + movingImagesEndpoint + "?q=" + keyword + "&page=" + strconv.Itoa(off+1)
		}
		jsonToReturn.Last = ec2url + "/api/" + movingImagesEndpoint + "?q=" + keyword + "&page=" + strconv.Itoa(maxPages)

		defer countResp.Body.Close()

		//Now, we do the actual query

		mainSearchQuery := buildMainSparqlQuery(keyword, offset, env)

		mainSearchParams := url.Values{}
		mainSearchParams.Add("query", mainSearchQuery)
		mainSearchParams.Add("format", "json")
		mainSearchBody := strings.NewReader(mainSearchParams.Encode())

		mainSearchReq, err := http.NewRequest("POST", neptuneurl, mainSearchBody)

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

		return c.JSONNonEncodePretty(http.StatusOK, jsonToReturn, " ")
	}
	return echo.HandlerFunc(fn)
}

// Moving Images Entity godoc
// @Summary Moving images get specific entity query
// @Description Moving images get specific entity query
// @Tags MovingImages Entity
// @Param id path string true "string id"
// @Produce json
// @Success 200 {object} EntityReturnStruct
// @Failure 500
// @Router /moving-images-ent/entity/{id} [get]
func movingImagesEntity(neptuneurl, graph string) echo.HandlerFunc {
	fn := func(c echo.Context) error {
		var jsonToReturn EntityReturnStruct

		id := c.Param("id")

		mainSearchQuery := buildEntityMainSparqlQuery(id, env)

		mainSearchParams := url.Values{}
		mainSearchParams.Add("query", mainSearchQuery)
		mainSearchParams.Add("format", "json")
		mainSearchBody := strings.NewReader(mainSearchParams.Encode())

		mainSearchReq, err := http.NewRequest("POST", neptuneurl, mainSearchBody)

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

		return c.JSONNonEncodePretty(http.StatusOK, jsonToReturn, " ")
	}
	return echo.HandlerFunc(fn)
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

	welcomeString := os.Getenv("WELCOME_STRING")
	ec2url := os.Getenv("EC2_URL")
	ec2port := os.Getenv("EC2_PORT")
	neptuneUrl := os.Getenv("NEPTUNE_URL")
	neptunePort := os.Getenv("NEPTUNE_PORT")
	discoveryAPIurl := os.Getenv("DISCOVERY_API")
	movingImagesEndpoint := os.Getenv("MOVING_IMAGES_ENDPOINT")
	graph := os.Getenv("GRAPH")

	neptuneFullSparqlUrl := neptuneUrl + ":" + neptunePort + "/sparql"
	ec2fullurl := ec2url + ":" + ec2port

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.GET("/", helloResponse(welcomeString))

	e.POST("/sparql", requestToNeptune(neptuneFullSparqlUrl, graph)) //to pass requests directly through

	e.GET("/discovery", fetchDiscovery(discoveryAPIurl))

	e.GET("/moving-images", movingImages(ec2fullurl, neptuneFullSparqlUrl, movingImagesEndpoint, graph))

	e.GET("/moving-images-ent/entity/:id", movingImagesEntity(neptuneFullSparqlUrl, graph))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":9000"))

}
