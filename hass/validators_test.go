package hass

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidatorsTestSuite struct {
	suite.Suite
}

func (suite *ValidatorsTestSuite) TestValidateEntityName() {
	suite.Nil(ValidateEntityName("shell.my_command"))
	suite.Nil(ValidateEntityName("1shell.1_my_command"))
	suite.NotNil(ValidateEntityName("shell.a.my_command"))
	suite.NotNil(ValidateEntityName("my_command"))
}

func (suite *ValidatorsTestSuite) TestValidateHostAndBearer() {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		suite.Equal(
			"GET",
			req.Method,
		)

		suite.Equal(
			"application/json",
			req.Header.Get("Content-Type"),
		)

		suite.Equal(
			"/api/",
			req.URL.Path,
		)

		if req.Header.Get("Authorization") != "Bearer VALID" {
			res.WriteHeader(403)
			res.Write([]byte(`{"message": "invalid"}`))
			return
		}

		res.WriteHeader(200)
		res.Write([]byte(`{"message": "API running."}`))
	}))

	defer func() { testServer.Close() }()

	suite.Nil(ValidateHostAndBearer(testServer.URL, "VALID"))
	suite.NotNil(ValidateHostAndBearer(testServer.URL, "INVALID"))
	suite.NotNil(ValidateHostAndBearer("invalid_addr", "VALID"))
}

func TestValidatorsTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorsTestSuite))

}
