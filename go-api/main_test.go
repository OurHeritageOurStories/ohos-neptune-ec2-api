package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
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

    t.Run("should get result count", func(t *testing.T) {
		if assert.NoError(t, helloResponse(c)) {
	        assert.Equal(t, http.StatusOK, rec.Code)
	        assert.Equal(t, `Hello, you've reached the Go API that lets you talk to the Neptune database. Well done!`, rec.Body.String())
	    }
	})

}

func TestEntity(t *testing.T) {
    url := "/movingImagesEnt/entity/(filmRef)5385"
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
		if assert.NoError(t, movingImagesEntity(c)) {
	        assert.Equal(t, http.StatusOK, rec.Code)
	        assert.Contains(t, rec.Body.String(), "(filmRef)5385", "identifier not correct")
	    }
	})

	t.Run("should get title", func(t *testing.T) {
		if assert.NoError(t, movingImagesEntity(c)) {
	        assert.Equal(t, http.StatusOK, rec.Code)
	        assert.Contains(t, rec.Body.String(), "('BOOTS' ON CORNER OF RENFIELD STREET AND ARGYLE STREET GLASGOW )", "title not correct")
	    }
	})

	t.Run("should get description", func(t *testing.T) {
		if assert.NoError(t, movingImagesEntity(c)) {
	        assert.Equal(t, http.StatusOK, rec.Code)
	        assert.Contains(t, rec.Body.String(), "Shot of the corner of Renfield Street and Argyle Street showing 'Boots' the chemist shop.", "description not correct")
	    }
	})

	t.Run("should get url", func(t *testing.T) {
		if assert.NoError(t, movingImagesEntity(c)) {
	        assert.Equal(t, http.StatusOK, rec.Code)
	        assert.Contains(t, rec.Body.String(), "http://movingimage.nls.uk/film/5385", "url not correct")
	    }
	})

	t.Run("should get topics", func(t *testing.T) {
		if assert.NoError(t, movingImagesEntity(c)) {
	        assert.Equal(t, http.StatusOK, rec.Code)
	        assert.Contains(t, rec.Body.String(), "Architecture and Buildings ||| Glasgow", "topics not correct")
	    }
	})

	t.Run("should get error for missing id", func(t *testing.T) {
		url := "movingImagesEnt/"
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
		if assert.NoError(t, movingImagesEntity(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
	        assert.Contains(t, rec.Body.String(), `{
 "items": []
}`, "Not throwing error")
	    }
	})

}

func TestDiscovery(t *testing.T) {
    
    t.Run("should get records from all archives", func(t *testing.T) {
		url := "discovery"
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
		if assert.NoError(t, fetchDiscovery(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
	        assert.Contains(t, rec.Body.String(), `"catalogueLevels": [],
 "closureStatuses": [],`, "Getting records only from TNA")
	    }
	})

	t.Run("should get records from all archives", func(t *testing.T) {
		url := "discovery?source=ALL"
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
		if assert.NoError(t, fetchDiscovery(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
	        assert.Contains(t, rec.Body.String(), `"catalogueLevels": [],
 "closureStatuses": [],`, "Getting records only from TNA")
	    }
	})

	t.Run("should get records from other archives", func(t *testing.T) {
		url := "discovery?source=OTH"
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
		if assert.NoError(t, fetchDiscovery(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
	        assert.Contains(t, rec.Body.String(), `"catalogueLevels": [],
 "closureStatuses": [],`, "Getting records only from TNA")
	    }
	})

	t.Run("should get records from TNA", func(t *testing.T) {
		url := "discovery?source=TNA"
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
		if assert.NoError(t, fetchDiscovery(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
	        assert.Contains(t, rec.Body.String(), `"closureStatuses": [
  {
   "code": "O",`, "Not getting records from TNA")
	    }
	})


	t.Run("should get records from all archives for edinburgh tram", func(t *testing.T) {
		url := "discovery?q=edinburgh%20tram"
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
		if assert.NoError(t, fetchDiscovery(c)) {
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
		url := "discovery?q=edinburgh%20tram&source=ALL"
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
		if assert.NoError(t, fetchDiscovery(c)) {
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
		url := "discovery?q=edinburgh%20tram&source=OTH"
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
		if assert.NoError(t, fetchDiscovery(c)) {
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
		url := "discovery?q=edinburgh%20tram&source=TNA"
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
		if assert.NoError(t, fetchDiscovery(c)) {
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
    url := "movingImages?q=glasgow&page=1"
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
		if assert.NoError(t, movingImages(c)) {
	        assert.Equal(t, http.StatusOK, rec.Code)
	        assert.Contains(t, rec.Body.String(), "178", "not getting the right data")
	    }
	})

	t.Run("should get error for missing page param", func(t *testing.T) {
		url := "movingImages?q=glasgow"
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
		if assert.NoError(t, movingImages(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
	        assert.Contains(t, rec.Body.String(), "You need to provide both a keyword and a page number", "Not throwing error")
	    }
	})

	t.Run("should get error for missing query param", func(t *testing.T) {
		url := "movingImages?page=1"
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
		if assert.NoError(t, movingImages(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
	        assert.Contains(t, rec.Body.String(), "You need to provide both a keyword and a page number", "Not throwing error")
	    }
	})

	t.Run("should get error for page provided as string", func(t *testing.T) {
		url := "movingImages?q=glasgow&page=string"
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
		if assert.NoError(t, movingImages(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
	        assert.Contains(t, rec.Body.String(), "Page needs to be selected as a number", "Not throwing error")
	    }
	})

	t.Run("should not get data for page greater than max page", func(t *testing.T) {
		url := "movingImages?q=glasgow&page=2500"
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
		if assert.NoError(t, movingImages(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "http status code")
	        assert.Contains(t, rec.Body.String(), `"items": []`, "Showing data for more than max page")
	    }
	})

	t.Run("should not get data for non-existing keyword", func(t *testing.T) {
		url := "movingImages?q=harshad&page=1"
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
		if assert.NoError(t, movingImages(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "Getting result for non existing keyword")
	        assert.Contains(t, rec.Body.String(), "The search worked, there just aren't any results", "Showing data for non-existing keyword")
	    }
	})

}


func TestFailOnPurpose(t *testing.T) {
	t.Run("should not get data for non-existing keyword", func(t *testing.T) {
		url := "movingImages?q=harshad&page=1"
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
		if assert.NoError(t, movingImages(c)) {
	        assert.Equal(t, http.StatusOK, res.StatusCode, "Getting result for non existing keyword")
	        assert.Contains(t, rec.Body.String(), "'items': []", "Not showing error")
	    }
	})
}