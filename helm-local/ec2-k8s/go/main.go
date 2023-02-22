package main

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/*
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
*/

func requestToNeptune(c echo.Context) error {
	// Get name
	name := c.FormValue("name")
	// Get avatar
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<b>Thank you! "+name+"</b>")
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
