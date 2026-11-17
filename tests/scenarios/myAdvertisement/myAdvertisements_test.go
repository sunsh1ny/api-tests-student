package myAdvertisement

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"

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

// Positive: Проверяем получение собственных объявлений в профиле
func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	ad := fixtures.NewAd(s.T(), s.loginData.Token)
	s.T().Log("Получаем список собственных объявлений и проверяем созданную карточку")
	body := myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK, 50)
	item := gjson.Get(body, `items.#(id=="`+ad.ID+`")`)
	require.True(s.T(), item.Exists(), "Объявление %s не найдено среди собственных", ad.ID)
	require.Equal(s.T(), s.loginData.User.Id, item.Get("user_id").String(), "user_id не совпадает")
	require.Equal(s.T(), ad.Title, item.Get("title").String(), "Название не совпадает")
	require.Equal(s.T(), ad.Description, item.Get("description").String(), "Описание не совпадает")
	require.NotEmpty(s.T(), item.Get("created_at").String(), "created_at пустое")
	require.NotEmpty(s.T(), item.Get("updated_at").String(), "updated_at пустое")
	require.NotEmpty(s.T(), item.Get("photos").Array(), "Нет фотографий")
}

// Negative: Проверяем 401 ошибку при попытке получить собственные объявления без токена
func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений с неправильным авторизационным токеном", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), "incorrect_token", http.StatusUnauthorized, 50)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		require.Equal(s.T(), "unauthorized", gjson.Get(advertisementsBody, "error").String(),
			"Некорректное значение error параметра")
		require.Equal(s.T(), "Invalid or expired token", gjson.Get(advertisementsBody, "message").String(),
			"Некорректное значение message параметра")
	})
}
