package del_test

import (
	"encoding/json"
	"github.com/Slava02/URL_shortner/internal/http-server/handlers/url/del"
	"github.com/Slava02/URL_shortner/internal/http-server/handlers/url/del/mocks"
	"github.com/Slava02/URL_shortner/internal/http-server/handlers/url/save"
	"github.com/Slava02/URL_shortner/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Parallel()

		urlDeleterMock := mocks.NewURLDeleter(t)

		if tc.respError == "" || tc.mockError != nil {
			urlDeleterMock.On("DeleteUrl", tc.alias).
				Return(tc.mockError).
				Once()
		}

		handler := del.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

		req, err := http.NewRequest(http.MethodDelete, "/test_alias", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusOK)

		body := rr.Body.String()

		var resp save.Response

		require.NoError(t, json.Unmarshal([]byte(body), &resp))

		require.Equal(t, tc.respError, resp.Error)
	}
}
