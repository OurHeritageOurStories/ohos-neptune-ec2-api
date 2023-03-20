package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

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

func TestMovingImagesGlasgow(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "Not connected to docker")
	resource, err := pool.Run("public.ecr.aws/j8z6n5u1/data-go-api", "latest", []string{})
	require.NoError(t, err, "could not start container")
	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failedToRemove")
	})
	var resp *http.Response

	t.Run("Test Glasgow to determine the data is present", func(t *testing.T) {
		err = pool.Retry(func() error {
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
		t.Run("should get result count", func(t *testing.T) {
			require.Contains(t, string(body), "178", "not the right count after all")
		})
	})

	t.Run("Test missing parameter works fine - no page count", func(t *testing.T) {
		err = pool.Retry(func() error {
			resp, err = http.Get("http://localhost:5000/api/movingImages?keyword=glasgow")
			if err != nil {
				t.Log("container still loading")
				return err
			}
			return nil
		})
		require.NoError(t, err, "http error")
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "http status code")
		/*
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "Didn't read the body")
			t.Run("Should get the missing param notification in the response", func(t *testing.T) {
				require.Contains(t, string(body), "Missing params, using default of keyword=glasgow, page=1", "not the right count after all")
			})
			t.Run("Should get the Glasgow result count", func(t *testing.T) {
				require.Contains(t, string(body), "178", "not the right count after all")
			})*/
	})

	t.Run("Test page count supplied as a word", func(t *testing.T) {
		err = pool.Retry(func() error {
			resp, err = http.Get("http://localhost:5000/api/movingImages?keyword=glasgow&page=saussage")
			if err != nil {
				t.Log("still loading")
				return err
			}
			return nil
		})
		require.NoError(t, err, "http error")
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Http status code")
	})

	t.Run("Test valid params that don't exist - huge page number", func(t *testing.T) {
		err = pool.Retry(func() error {
			resp, err = http.Get("http://localhost:5000/api/movingImages?keyword=glasgow&page=10000")
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
		t.Run("Should get all the results", func(t *testing.T) {
			require.Contains(t, string(body), "5021", "not the right count after all")
		})
	})

	t.Run("Test valid params that don't exist - weird keyword", func(t *testing.T) {
		err = pool.Retry(func() error {
			resp, err = http.Get("http://localhost:5000/api/movingImages?keyword=adsj&page=1")
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
		t.Run("Should get all the results", func(t *testing.T) {
			require.Contains(t, string(body), "5021", "not the right count after all")
		})
		t.Run("Should get all the results", func(t *testing.T) {
			require.Contains(t, string(body), "Water and Waterways ||| Celebrations, Traditions and Customs ||| Transport ||| Ships and Shipping ||| Glasgow", "Didn't get an expected result")
		})
	})

}
