package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AndrewBewseyTNA/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	welcomeString := os.Getenv("WELCOME_STRING")
	url := "/"
	e := echo.New()
	req, err := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("The request could not be created because of: %v", err)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := rec.Result()
	defer res.Body.Close()

	t.Run("Should get the welcome string back.", func(t *testing.T) {
		if assert.NoError(t, helloResponse(welcomeString)(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, welcomeString, rec.Body.String())
		}
	})

}

func TestEntity(t *testing.T) {
	neptuneUrl := os.Getenv("NEPTUNE_URL")
	neptunePort := os.Getenv("NEPTUNE_PORT")
	movingImagesEntityEndpoint := os.Getenv("MOVING_IMAGES_ENTITY_ENDPOINT")
	neptuneFullSparqlUrl := neptuneUrl + ":" + neptunePort + "/sparql"
	url := "/" + movingImagesEntityEndpoint + "/entity/(filmRef)5385"
	e := echo.New()
	req, err := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("The request could not be created because of: %v", err)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("(filmRef)5385")

	res := rec.Result()
	defer res.Body.Close()

	t.Run("should get identifier", func(t *testing.T) {
		if assert.NoError(t, movingImagesEntity(neptuneFullSparqlUrl)(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "(filmRef)5385", "identifier not correct")
		}
	})

	t.Run("should get title", func(t *testing.T) {
		if assert.NoError(t, movingImagesEntity(neptuneFullSparqlUrl)(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "('BOOTS' ON CORNER OF RENFIELD STREET AND ARGYLE STREET GLASGOW )", "title not correct")
		}
	})

	t.Run("should get description", func(t *testing.T) {
		if assert.NoError(t, movingImagesEntity(neptuneFullSparqlUrl)(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Shot of the corner of Renfield Street and Argyle Street showing 'Boots' the chemist shop.", "description not correct")
		}
	})

	t.Run("should get url", func(t *testing.T) {
		if assert.NoError(t, movingImagesEntity(neptuneFullSparqlUrl)(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "http://movingimage.nls.uk/film/5385", "url not correct")
		}
	})

	t.Run("should get topics", func(t *testing.T) {
		if assert.NoError(t, movingImagesEntity(neptuneFullSparqlUrl)(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Architecture and Buildings ||| Glasgow", "topics not correct")
		}
	})

	t.Run("should get error for missing id", func(t *testing.T) {
		url := movingImagesEntityEndpoint + "/"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("")
		c.SetParamNames("")
		c.SetParamValues("")

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, movingImagesEntity(neptuneFullSparqlUrl)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), `{
 "items": []
}`, "Not throwing error")
		}
	})

}

func TestDiscovery(t *testing.T) {
	discoveryAPIurl := os.Getenv("DISCOVERY_API")
	discoveryEndpoint := os.Getenv("DISCOVERY_ENDPOINT")

	t.Run("should get records from all archives", func(t *testing.T) {
		url := discoveryEndpoint
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, fetchDiscovery(discoveryAPIurl)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), `"catalogueLevels": [],
 "closureStatuses": [],`, "Getting records only from TNA")
		}
	})

	t.Run("should get records from all archives", func(t *testing.T) {
		url := discoveryEndpoint + "?source=ALL"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, fetchDiscovery(discoveryAPIurl)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), `"catalogueLevels": [],
 "closureStatuses": [],`, "Getting records only from TNA")
		}
	})

	t.Run("should get records from other archives", func(t *testing.T) {
		url := discoveryEndpoint + "?source=OTH"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, fetchDiscovery(discoveryAPIurl)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), `"catalogueLevels": [],
 "closureStatuses": [],`, "Getting records only from TNA")
		}
	})

	t.Run("should get records from TNA", func(t *testing.T) {
		url := discoveryEndpoint + "?source=TNA"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, fetchDiscovery(discoveryAPIurl)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), `"closureStatuses": [
  {
   "code": "O",`, "Not getting records from TNA")
		}
	})

	t.Run("should get records from all archives for edinburgh tram", func(t *testing.T) {
		url := discoveryEndpoint + "?q=edinburgh%20tram"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, fetchDiscovery(discoveryAPIurl)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), `{
 "catalogueLevels": [],
 "closureStatuses": [],
 "count": 6,
 "departments": [],
 "heldByReps": [
  {
   "code": "OTH",
   "count": 4
  },
  {
   "code": "TNA",
   "count": 2
  }`, "Not getting all records from all archives")
		}
	})

	t.Run("should get records from all archives for edinburgh tram", func(t *testing.T) {
		url := discoveryEndpoint + "?q=edinburgh%20tram&source=ALL"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, fetchDiscovery(discoveryAPIurl)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), `{
 "catalogueLevels": [],
 "closureStatuses": [],
 "count": 6,
 "departments": [],
 "heldByReps": [
  {
   "code": "OTH",
   "count": 4
  },
  {
   "code": "TNA",
   "count": 2
  }`, "Not getting all records from all archives")
		}
	})

	t.Run("should get records from other archives for edinburgh tram", func(t *testing.T) {
		url := discoveryEndpoint + "?q=edinburgh%20tram&source=OTH"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, fetchDiscovery(discoveryAPIurl)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), `"catalogueLevels": [],
 "closureStatuses": [],
 "count": 4,
 "departments": [],
 "heldByReps": [
  {
   "code": "OTH",
   "count": 4
  }`, "Not getting all records from other archives")
		}
	})

	t.Run("should get records from TNA for edinburgh tram", func(t *testing.T) {
		url := discoveryEndpoint + "?q=edinburgh%20tram&source=TNA"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, fetchDiscovery(discoveryAPIurl)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), `{
   "code": "6",
   "count": 1
  },
  {
   "code": "7",
   "count": 1
  }`, "Not getting all records for edinburgh tram from TNA")
		}
	})

}

func TestMovingImages(t *testing.T) {
	ec2url := os.Getenv("EC2_URL")
	ec2port := os.Getenv("EC2_PORT")
	neptuneUrl := os.Getenv("NEPTUNE_URL")
	neptunePort := os.Getenv("NEPTUNE_PORT")
	movingImagesEndpoint := os.Getenv("MOVING_IMAGES_ENDPOINT")
	neptuneFullSparqlUrl := neptuneUrl + ":" + neptunePort + "/sparql"
	ec2fullurl := ec2url + ":" + ec2port
	url := movingImagesEndpoint + "?q=glasgow&page=1"
	e := echo.New()
	req, err := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("The request could not be created because of: %v", err)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := rec.Result()
	defer res.Body.Close()

	t.Run("should get data for Glasgow", func(t *testing.T) {
		if assert.NoError(t, movingImages(ec2fullurl, neptuneFullSparqlUrl, movingImagesendpoint)(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "178", "not getting the right data")
		}
	})

	t.Run("should get error for missing page param", func(t *testing.T) {
		url := movingImagesEndpoint + "?q=glasgow"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, movingImages(ec2fullurl, neptuneFullSparqlUrl, movingImagesendpoint)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), "You need to provide both a keyword (as q) and a page number", "Not throwing error")
		}
	})

	t.Run("should get error for missing query param", func(t *testing.T) {
		url := movingImagesEndpoint + "?page=1"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, movingImages(ec2fullurl, neptuneFullSparqlUrl, movingImagesendpoint)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), "You need to provide both a keyword (as q) and a page number", "Not throwing error")
		}
	})

	t.Run("should get error for page provided as string", func(t *testing.T) {
		url := movingImagesEndpoint + "?q=glasgow&page=string"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, movingImages(ec2fullurl, neptuneFullSparqlUrl, movingImagesendpoint)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), "Page needs to be selected as a number", "Not throwing error")
		}
	})

	t.Run("should not get data for page greater than max page", func(t *testing.T) {
		url := movingImagesEndpoint + "?q=glasgow&page=2500"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, movingImages(ec2fullurl, neptuneFullSparqlUrl, movingImagesendpoint)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
			assert.Contains(t, rec.Body.String(), `"items": []`, "Showing data for more than max page")
		}
	})

	t.Run("should not get data for non-existing keyword", func(t *testing.T) {
		url := movingImagesEndpoint + "?q=harshad&page=1"
		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("The request could not be created because of: %v", err)
		}
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, movingImages(ec2fullurl, neptuneFullSparqlUrl, movingImagesendpoint)(c)) {
			assert.Equal(t, http.StatusOK, res.StatusCode, "Getting result for non existing keyword")
			assert.Contains(t, rec.Body.String(), `"items": []`, "Showing data for non-existing keyword")
		}
	})

}
