package http

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type clientSuite struct {
	suite.Suite

	url string

	client *Client
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, &clientSuite{})
}

func (s *clientSuite) SetupTest() {
	httpmock.Activate()

	s.url = "http://localhost:8080/test"

	s.client = NewClient(&http.Client{})
}

func (s *clientSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (s *clientSuite) TestPostErr() {
	err := s.client.Post(s.url, []byte(""))

	s.EqualError(
		errors.New("can't send message: Post \"http://localhost:8080/test\": no responder found"),
		err.Error(),
	)
}

func (s *clientSuite) TestPostUnexpectedStatus() {
	httpmock.RegisterResponder(
		http.MethodPost,
		s.url,
		httpmock.NewStringResponder(http.StatusInternalServerError, ""),
	)

	err := s.client.Post(s.url, []byte(""))

	s.Equal(errors.New(fmt.Sprintf("unexpected status: %d", http.StatusInternalServerError)), err)
}

func (s *clientSuite) TestPostOk() {
	httpmock.RegisterResponder(
		http.MethodPost,
		s.url,
		httpmock.NewStringResponder(http.StatusOK, ""),
	)

	err := s.client.Post(s.url, []byte(""))

	s.Nil(err)
}
