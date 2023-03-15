package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

/*
func TestBase(t *testing.T) {
	url := "/"
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

func TestContains(t *testing.T) {
	a := `{"name":"dave","age":3}`
	b := `"name":"dave"`
	assert.Contains(t, a, b)
}
*/

func TestDocker(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "Not connected to docker")
	resource, err := pool.Run("public.ecr.aws/j8z6n5u1/data-go-api", "latest", []string{})
	require.NoError(t, err, "could not start container")
	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failedToRemove")
	})
	var resp *http.Response

	err = pool.Retry(func() error {
		//resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("5000/tcp"), "/api/"))
		resp, err = http.Get("http://localhost:5000/api")
		if err != nil {
			t.Log("container still loading")
			return err
		}
		return nil
	})
	require.NoError(t, err, "http error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "http status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Didn't read the body")
	require.Contains(t, string(body), "done", "not done after all")

}

func TestResponse(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "Not connected to docker")
	resource, err := pool.Run("public.ecr.aws/j8z6n5u1/data-go-api", "latest", []string{})
	require.NoError(t, err, "could not start container")
	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failedToRemove")
	})
	var resp *http.Response

	err = pool.Retry(func() error {
		//resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("5000/tcp"), "/api/"))
		resp, err = http.Get("http://localhost:5000/api/movingImages?keyword=glasgow&page=1")
		if err != nil {
			t.Log("container still loading")
			return err
		}
		return nil
	})
	require.NoError(t, err, "http error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "http status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Didn't read the body")
	require.Contains(t, string(body), "178", "not the right count after all")
}

/*
func TestMovingImages(t *testing.T) {
	baseUrl := "/movingImages"
	t.Run("should return teapot", func(t *testing.T) {
		url := baseUrl
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, movingImages(c)) {
			assert.Equal(t, http.StatusTeapot, rec.Code)
		}
	})

	t.Run("should check it has a count", func(t *testing.T) {
		url := baseUrl + "?keyword=glasgow&page=1"
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		res := rec.Result()
		defer res.Body.Close()
		if assert.NoError(t, movingImages(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, "Hello", "Hello")
		}
	})
}
*/
