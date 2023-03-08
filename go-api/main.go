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

	return c.JSON(http.StatusOK, jsonMap)
}

func getEntities(c echo.Context) error {
	//	keyword := c.Request().URL.Query().Get("keyword")
	//	//page := strings.ToUpper(c.Request().URL.Query().Get("page"))
	//	page := c.Request().URL.Query().Get("page")

	//keyword := c.QueryParam("keyword\n")
	//page := c.QueryParam("page\n")

	//queryParams := c.Request().URL.Query()
	//fmt.Println(len(queryParams))
	//for i, s := range queryParams {
	//	fmt.Println(i, s)
	//}
	//keyword := queryParams.Get("keyword")
	//page := queryParams.Get("page")

	//keyword := queryParams["keyword"]
	//pagething := queryParams["thingiy"]
	//pagething := queryParams.Get("thingiy")

	//fmt.Print("keyword" + strings.Join(keyword, " "))
	//fmt.Print("page" + strings.Join(pagething, " "))
	//pagethingalt := strings.Join(pagething, " ")
	//return c.String(http.StatusOK, pagething)

	keyword := c.QueryParam("keyword")
	page := c.QueryParam("page")
	fmt.Println(keyword)
	fmt.Println(page)
	return c.String(http.StatusOK, "keyword -> "+keyword+"   | page -> "+page)
	/*
		off, err := strconv.Atoi(page)

		if err != nil {
			log.Fatal(err)
		}

		off = max(1, off)

		o := strconv.Itoa((off - 1) * 10)

		constructedQuery := "prefix tanc: <http://tanc.manchester.ac.uk/> SELECT DISTINCT ?o (count(?text) as ?count) WHERE { ?s <http://tanc.manchester.ac.uk/text> ?text. ?s tanc:mentions ?o FILTER (regex(str(?o), '" + keyword + "', 'i'))} GROUP BY ?o ORDER BY DESC(?count) OFFSET " + o + " LIMIT 10"

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

		constructedQueryTwo := "prefix tanc: <http://tanc.manchester.ac.uk/> SELECT (count(*) as ?count) WHERE {SELECT DISTINCT ?o (count(?text) as ?count) WHERE { ?s <http://tanc.manchester.ac.uk/text> ?text. ?s tanc:mentions ?o FILTER (regex(str(?o), '" + keyword + "', 'i'))} GROUP BY ?o ORDER BY DESC(?count)}"

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

		return c.JSON(http.StatusOK, jsonMap)*/
}

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

	e.GET("/entities", getEntities)

	e.GET("/testentities", entitytest)

	//e.GET("/entities/{entity}", getEntity)

	e.Logger.Fatal(e.Start(":9000"))

}
