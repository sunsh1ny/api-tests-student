package advertisementCreation

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/fixtures"
)

const photoCount = 3

func TestCreateAdvertisementWithAllFields(t *testing.T) {
	owner := fixtures.Login(t, false)
	title := fixtures.UniqueTitle()
	form := fixtures.AdForm(t, title, photoCount)

	t.Log("Отправляем POST /advertisement со всеми полями и тремя фотографиями")
	created := fixtures.CreateAd(t, owner.Token, form)

	t.Log("Читаем карточку по ID и сравниваем поля с формой")
	record := advertisements.Get(t, owner.Token, created.ID)
	expected := models.Advertisement{
		ID: created.ID, Title: title, Description: form.Fields[fixtures.DescriptionField],
		Price: fixtures.Price, Quantity: fixtures.Quantity, UserID: owner.User.Id,
	}
	require.Equal(t, expected, models.Advertisement{
		ID: record.ID, Title: record.Title, Description: record.Description,
		Price: record.Price, Quantity: record.Quantity, UserID: record.UserID,
	}, "Карточка объявления отличается от отправленных полей")

	t.Logf("Ищем объявление по полному названию %q", title)
	results := advertisements.GetAdvertisementsList(t, owner.Token, title)
	titles := foundTitlesByID(results.Items)
	require.Contains(t, titles, created.ID, "Поиск не вернул только что созданное объявление")
	require.Equal(t, title, titles[created.ID], "Поиск вернул другое название для созданного объявления")

	t.Log("Запрашиваем фотографии и сверяем их с результатом загрузки")
	photos := advertisements.Photos(t, owner.Token, created.ID)
	require.Len(t, photos, len(form.Photos), "GET /advertisements/{id}/photos вернул неверное число фотографий")
	require.ElementsMatch(t, created.Photos, photos, "Список фотографий отличается от ответа на создание")
	require.ElementsMatch(t, record.Photos, photos, "Список фотографий отличается от карточки")
	ids := map[string]bool{}
	var orders []int
	for _, photo := range photos {
		require.NotEmpty(t, photo.ID, "В ответе отсутствует ID фото")
		require.False(t, ids[photo.ID], "В ответе повторяется фото %s", photo.ID)
		ids[photo.ID] = true
		link, err := url.ParseRequestURI(photo.URL)
		require.NoError(t, err, "В ответе некорректный URL фото")
		require.Equal(t, "https", link.Scheme, "Фото должно быть доступно по HTTPS")
		require.NotEmpty(t, link.Host, "У ссылки на фото отсутствует хост")
		orders = append(orders, photo.SortOrder)
	}
	require.ElementsMatch(t, []int{0, 1, 2}, orders, "Неверные позиции загруженных фото")
	t.Log("Поля карточки, результат поиска и фотографии совпали")
}

// foundTitlesByID индексирует результаты поиска по идентификатору объявления.
func foundTitlesByID(items []models.Advertisement) map[string]string {
	titles := make(map[string]string, len(items))
	for _, item := range items {
		titles[item.ID] = item.Title
	}
	return titles
}
