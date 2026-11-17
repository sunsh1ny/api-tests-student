package advertisementSearch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advList "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	s.loginData = fixtures.Login(s.T(), false)
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestSearchByTitle() {
	ad := fixtures.NewAd(s.T(), s.loginData.Token)
	list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, ad.Title)
	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)
	})
	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, ad.Title)
	})
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	const searchSegmentSize = 10
	unique := strings.TrimPrefix(fixtures.UniqueTitle(), "QA ")
	beginning := unique[:searchSegmentSize]
	middle := "QA" + unique[searchSegmentSize:2*searchSegmentSize]
	end := unique[2*searchSegmentSize:]
	title := beginning + " " + middle + " " + end
	ad := fixtures.CreateAd(s.T(), s.loginData.Token, fixtures.AdForm(s.T(), title, 1))

	cases := []struct {
		name   string
		search string
	}{
		{"начало", beginning},
		{"середина", middle},
		{"конец", end},
		{"нижний регистр", strings.ToLower(middle)},
		{"верхний регистр", strings.ToUpper(middle)},
		{"смешанный регистр", utils.MixCase(middle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, ad.Title)
		})
	}
}
