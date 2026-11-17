package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/client/http/favorites"
	apiResponse "api-tests-template/internal/helpers/api-response"
	"api-tests-template/internal/managers/advertisements/models"
)

const PageSize = 50

type Favorite struct {
	ID            string               `json:"favorite_id"`
	Advertisement models.Advertisement `json:"advertisement"`
}

type Page struct {
	Items      []Favorite `json:"items"`
	TotalPages int        `json:"total_pages"`
	Total      int        `json:"total"`
}

// List возвращает одну страницу избранного.
func List(t *testing.T, token string, page int) Page {
	t.Helper()
	return apiResponse.Decode[Page](t, favorites.Get(t, token, page, PageSize), http.StatusOK)
}

// ListError возвращает ошибку чтения избранного.
func ListError(t *testing.T, token string, status int) apiResponse.Error {
	t.Helper()
	return apiResponse.Decode[apiResponse.Error](t, favorites.Get(t, token, 1, PageSize), status)
}

// Add добавляет объявление в избранное и проверяет HTTP 201.
func Add(t *testing.T, token, adID string) {
	t.Helper()
	response := favorites.Add(t, token, adID)
	defer apiResponse.Close(t, response)
	apiResponse.CheckStatus(t, response, http.StatusCreated)
}

// AddError возвращает ошибку добавления в избранное.
func AddError(t *testing.T, token, adID string, status int) apiResponse.Error {
	t.Helper()
	return apiResponse.Decode[apiResponse.Error](t, favorites.Add(t, token, adID), status)
}

// Delete удаляет объявление из списка избранного и проверяет HTTP 204.
func Delete(t *testing.T, token, adID string) {
	t.Helper()
	response := favorites.Delete(t, token, adID)
	defer apiResponse.Close(t, response)
	apiResponse.CheckStatus(t, response, http.StatusNoContent)
}

// DeleteIfPresent удаляет тестовую запись, которая могла исчезнуть ранее.
func DeleteIfPresent(t *testing.T, token, adID string) {
	t.Helper()
	response := favorites.Delete(t, token, adID)
	defer apiResponse.Close(t, response)
	if response != nil && response.StatusCode == http.StatusNotFound {
		return
	}
	apiResponse.CheckStatus(t, response, http.StatusNoContent)
}

// DeleteError возвращает ошибку удаления из избранного.
func DeleteError(t *testing.T, token, adID string, status int) apiResponse.Error {
	t.Helper()
	return apiResponse.Decode[apiResponse.Error](t, favorites.Delete(t, token, adID), status)
}
