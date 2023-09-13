package client_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/jdbann/forestry/pkg/client"
)

func TestClient(t *testing.T) {
	type TestRequest struct{ Name string }
	type TestResponse struct{ Message string }

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/test" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var payload TestRequest
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, _ = fmt.Fprintf(w, `{"message":"well done, %s!"}`, payload.Name)
	}))

	c, err := client.New(srv.URL)
	assert.NilError(t, err)

	gotStatus, gotResponse, err := client.MakeRequest[TestResponse](c, http.MethodGet, "/test", TestRequest{Name: "john"})
	assert.NilError(t, err)
	assert.Equal(t, gotStatus, http.StatusOK)
	assert.DeepEqual(t, gotResponse, TestResponse{Message: "well done, john!"})
}
