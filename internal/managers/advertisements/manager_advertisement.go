package advertisements

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	apiResponse "api-tests-template/internal/helpers/api-response"
	"api-tests-template/internal/managers/advertisements/models"
)

type Form = advertisements.Form

// Create сохраняет объявление и возвращает его карточку.
func Create(t *testing.T, token string, form Form) models.Advertisement {
	t.Helper()
	return apiResponse.Decode[models.Advertisement](t, advertisements.Create(t, token, form), http.StatusCreated)
}

// CreateError возвращает ошибку создания и удаляет карточку, если сервер её ошибочно принял.
func CreateError(t *testing.T, token string, form Form, status int) apiResponse.Error {
	t.Helper()
	response := advertisements.Create(t, token, form)
	if response != nil && response.StatusCode == http.StatusCreated && status == http.StatusBadRequest {
		ad := apiResponse.Decode[models.Advertisement](t, response, http.StatusCreated)
		t.Cleanup(func() { Delete(t, token, ad.ID) })
		require.FailNowf(t, "Некорректное объявление принято", "POST /advertisement вернул %d вместо %d, id=%s", http.StatusCreated, status, ad.ID)
	}
	return apiResponse.Decode[apiResponse.Error](t, response, status)
}

// Get возвращает полную карточку объявления.
func Get(t *testing.T, token, adID string) models.Advertisement {
	t.Helper()
	return apiResponse.Decode[models.Advertisement](t, advertisements.Get(t, token, adID), http.StatusOK)
}

// Photos возвращает фотографии объявления из отдельной ручки.
func Photos(t *testing.T, token, adID string) []models.Photo {
	t.Helper()
	return apiResponse.Decode[[]models.Photo](t, advertisements.Photos(t, token, adID), http.StatusOK)
}

// Delete удаляет объявление и проверяет успешный ответ.
func Delete(t *testing.T, token, adID string) {
	t.Helper()
	response := advertisements.Delete(t, token, adID)
	defer apiResponse.Close(t, response)
	apiResponse.CheckStatus(t, response, http.StatusNoContent)
}
