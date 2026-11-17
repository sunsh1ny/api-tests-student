package favorites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

// loginAccounts возвращает токены владельца объявления и второго пользователя.
func loginAccounts(t *testing.T) (owner, actor string) {
	t.Helper()
	first := fixtures.Login(t, false)
	second := fixtures.Login(t, true)
	require.NotEqual(t, first.User.Id, second.User.Id, "Для избранного нужны два разных аккаунта TestBoard")
	return first.Token, second.Token
}

// addFavorite добавляет запись и регистрирует её очистку, даже если следующая проверка упадёт.
func addFavorite(t *testing.T, token, adID string) {
	t.Helper()
	t.Logf("Добавляем %s в избранное через POST /favorites", adID)
	t.Cleanup(func() { favorites.DeleteIfPresent(t, token, adID) })
	favorites.Add(t, token, adID)
}

// checkFavorite проверяет присутствие объявления по ID во всех страницах списка.
func checkFavorite(t *testing.T, token, adID string, expected bool) {
	t.Helper()
	t.Logf("Читаем GET /favorites, ожидаем присутствие %s: %t", adID, expected)
	for page := 1; ; page++ {
		result := favorites.List(t, token, page)
		for _, item := range result.Items {
			if item.Advertisement.ID == adID {
				require.True(t, expected, "Объявление %s осталось на странице %d избранного", adID, page)
				require.NotEmpty(t, item.ID, "У объявления %s в списке отсутствует favorite_id", adID)
				t.Logf("Объявление найдено на странице %d", page)
				return
			}
		}
		if page >= result.TotalPages {
			break
		}
	}
	require.False(t, expected, "Объявление %s отсутствует в избранном", adID)
	t.Log("Объявление отсутствует в избранном")
}
