package favorites

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// Get запрашивает страницу списка избранного.
func Get(t *testing.T, token string, page, size int) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page", strconv.Itoa(page)).
		Query("page_size", strconv.Itoa(size)).
		Expect(t).End().Response
}

// Add отправляет идентификатор объявления в избранное.
func Add(t *testing.T, token, adID string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		JSON(map[string]string{"advertisement_id": adID}).
		Expect(t).End().Response
}

// Delete удаляет объявление из избранного по его идентификатору.
func Delete(t *testing.T, token, adID string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoritesPath + "/" + url.PathEscape(adID)).
		Expect(t).End().Response
}
