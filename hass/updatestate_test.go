package hass

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UpdateStateTestSuite struct {
	suite.Suite
}

func (suite *UpdateStateTestSuite) TestOnStateCreation() {
	suite.updateWithStatus(http.StatusCreated)
}

func (suite *UpdateStateTestSuite) TestOnStateUpdate() {
	suite.updateWithStatus(http.StatusOK)
}

func (suite *UpdateStateTestSuite) updateWithStatus(status int) {
	callReceived := false
	payload := "state"
	bearer := "ABC"
	entity := "a.b.c"

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		suite.False(callReceived)
		callReceived = true

		suite.Equal(
			"POST",
			req.Method,
		)

		suite.Equal(
			"Bearer ABC",
			req.Header.Get("Authorization"),
		)

		suite.Equal(
			"application/json",
			req.Header.Get("Content-Type"),
		)

		suite.Equal(
			"/api/states/"+entity,
			req.URL.Path,
		)

		body, err := io.ReadAll(req.Body)
		suite.Nil(err)
		suite.Equal(payload, string(body))

		res.WriteHeader(status)
	}))

	defer func() { testServer.Close() }()

	hass := NewHass(
		bearer,
		testServer.URL,
		entity,
	)

	err := hass.UpdateState(string(payload))

	suite.Nil(err)
	suite.True(callReceived)
}

func (suite *UpdateStateTestSuite) TestOnHassFailure() {
	callReceived := false

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		suite.False(callReceived)
		callReceived = true

		res.WriteHeader(http.StatusInternalServerError)
	}))

	defer func() { testServer.Close() }()

	hass := NewHass(
		"bearer",
		testServer.URL,
		"entity",
	)

	err := hass.UpdateState("state")

	suite.NotNil(err)
}

func TestUpdateStateTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateStateTestSuite))
}
