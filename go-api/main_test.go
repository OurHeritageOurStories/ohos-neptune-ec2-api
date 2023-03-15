package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	url := "/"
	//url := "http://ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/api/movingImages?keyword=glasgow&page=1"
	// url := "https://www.google.com"
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	res := rec.Result()
	defer res.Body.Close()
	if assert.NoError(t, helloResponse(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

/*
type Controller struct {
}

func (m *Controller) movingImages(c echo.Context) error {
	return nil
}

func testBasic(t *testing.T) {
	t.Run("should return 200", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		controller := Controller{}
		controller.movingImages(c)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
*/
/*
func testBasic(t *assert.TestingT) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, h.movingImages(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	/*e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	controller := Controller{}
	controller.GetAllBooks(c)*/
/*
}

/*
func testBasic(t *testing.T) {
	t.Run("should return 200", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		controller := Controller{}
		controller.GetAllBooks(c)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
*/
/*
func testGet(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler
}
*/
