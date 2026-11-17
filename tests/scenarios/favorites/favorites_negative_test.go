package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	apiResponse "api-tests-template/internal/helpers/api-response"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

const favoritesLimit = 50

func TestFavoritesRejectNonexistentAdvertisement(t *testing.T) {
	user := fixtures.Login(t, false)
	ad := fixtures.NewAd(t, user.Token)
	ad.Delete(t)

	t.Run("добавление", func(t *testing.T) {
		t.Logf("Отправляем POST /favorites для удалённого объявления %s", ad.ID)
		actual := favorites.AddError(t, user.Token, ad.ID, http.StatusNotFound)
		fixtures.AssertError(t, actual, "not_found", "Advertisement not found")
	})
	t.Run("удаление", func(t *testing.T) {
		t.Logf("Отправляем DELETE /favorites/{ad_id} для удалённого объявления %s", ad.ID)
		actual := favorites.DeleteError(t, user.Token, ad.ID, http.StatusNotFound)
		fixtures.AssertError(t, actual, "not_found", "Favorite entry not found")
	})
}

func TestFavoritesRequiresAuthorization(t *testing.T) {
	owner := fixtures.Login(t, false)
	ad := fixtures.NewAd(t, owner.Token)
	operations := []struct {
		name string
		call func(*testing.T, string, int) apiResponse.Error
	}{
		{"GET /favorites", favorites.ListError},
		{"POST /favorites", func(t *testing.T, token string, code int) apiResponse.Error {
			t.Helper()
			return favorites.AddError(t, token, ad.ID, code)
		}},
		{"DELETE /favorites/{ad_id}", func(t *testing.T, token string, code int) apiResponse.Error {
			t.Helper()
			return favorites.DeleteError(t, token, ad.ID, code)
		}},
	}
	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			for _, tc := range fixtures.UnauthorizedCases() {
				t.Run(tc.Name, func(t *testing.T) {
					t.Logf("Отправляем %s %s", operation.name, tc.Name)
					actual := operation.call(t, tc.Token, http.StatusUnauthorized)
					fixtures.AssertError(t, actual, "unauthorized", tc.Message)
				})
			}
		})
	}
}

func TestFavoritesRejectDuplicate(t *testing.T) {
	owner, actor := loginAccounts(t)
	ad := fixtures.NewAd(t, owner)
	addFavorite(t, actor, ad.ID)

	t.Log("Повторно добавляем то же объявление")
	actual := favorites.AddError(t, actor, ad.ID, http.StatusConflict)
	fixtures.AssertError(t, actual, "already_exists", "Advertisement is already in favorites")
}

func TestFavoritesRejectOwnAdvertisement(t *testing.T) {
	owner := fixtures.Login(t, false)
	ad := fixtures.NewAd(t, owner.Token)

	t.Log("Владелец пытается добавить собственное объявление в избранное")
	actual := favorites.AddError(t, owner.Token, ad.ID, http.StatusForbidden)
	fixtures.AssertError(t, actual, "forbidden", "Cannot add own advertisement to favorites")
}

func TestFavoritesRejectMalformedAdvertisementID(t *testing.T) {
	user := fixtures.Login(t, false)
	const invalidID = "not-a-uuid"
	t.Log("Добавляем в избранное объявление с некорректным UUID")
	actual := favorites.AddError(t, user.Token, invalidID, http.StatusBadRequest)
	fixtures.AssertError(t, actual, "invalid_request", "invalid UUID length: 10")
}

func TestFavoritesRejectOverLimit(t *testing.T) {
	owner, actor := loginAccounts(t)
	initial := favorites.List(t, actor, 1).Total
	require.LessOrEqual(t, initial, favoritesLimit, "До теста превышен максимальный размер избранного")

	t.Logf("Доводим число избранных с %d до лимита %d", initial, favoritesLimit)
	for range favoritesLimit - initial {
		ad := fixtures.NewAd(t, owner)
		addFavorite(t, actor, ad.ID)
	}
	require.Equal(t, favoritesLimit, favorites.List(t, actor, 1).Total, "Не удалось заполнить избранное до лимита")

	ad := fixtures.NewAd(t, owner)
	t.Logf("Пытаемся добавить объявление №%d", favoritesLimit+1)
	actual := favorites.AddError(t, actor, ad.ID, http.StatusBadRequest)
	fixtures.AssertError(t, actual, "limit_reached", fmt.Sprintf("Favorites limit of %d reached", favoritesLimit))
}
