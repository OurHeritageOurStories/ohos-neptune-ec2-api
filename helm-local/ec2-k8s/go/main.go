package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func requestToNeptune(c *gin.Context) {
	sparqlQuery := "query=select ?s ?p ?o where {?s ?p ?o} limit 10"
	sparqlBuffer := bytes.NewBuffer([]byte(sparqlQuery))
	resp, err := http.Post("https://ohos-live-data-neptune.cluster-ro-c7ehmaoz3lrl.eu-west-2.neptune.amazonaws.com:8182/sparql", "data-binary", sparqlBuffer)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatal(err)
	}
	sb := string(body)
	log.Print(sb)
}

func testActive(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "hi, you've reached the goapi that lets you talk to neptune. Well done!")
}

func main() {
	router := gin.Default()
	router.POST("/other", requestToNeptune)
	router.GET("/helloTest", testActive)
	router.Run("localhost:9000")
}
