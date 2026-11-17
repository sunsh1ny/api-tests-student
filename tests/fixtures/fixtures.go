package fixtures

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	base "api-tests-template/tests"
)

const (
	TitleField       = "title"
	DescriptionField = "description"
	PriceField       = "price"
	QuantityField    = "quantity"

	Price    = 125.25
	Quantity = 3
)

var setupEnv sync.Once

type Ad struct {
	models.Advertisement
	token   string
	deleted bool
}

// Login авторизует основной или дополнительный тестовый аккаунт по настройкам окружения.
func Login(t *testing.T, secondary bool) authModels.LoginOkResponse {
	t.Helper()
	LoadEnv()
	prefix := "TEST"
	if secondary {
		prefix = "TEST_SECOND"
	}
	login, password := os.Getenv(prefix+"_LOGIN"), os.Getenv(prefix+"_PASSWORD")
	require.NotEmpty(t, login, "Укажите %s_LOGIN в .env.override", prefix)
	require.NotEmpty(t, password, "Укажите %s_PASSWORD в .env.override", prefix)
	t.Logf("Подготовка: авторизуем аккаунт %s", prefix)
	user := auth.Login(t, login, password)
	require.NotEmpty(t, user.Token, "В ответе авторизации нет токена")
	require.NotEmpty(t, user.User.Id, "В ответе авторизации нет ID пользователя")
	return user
}

// LoadEnv загружает настройки стенда, в том числе для тестов без авторизации.
func LoadEnv() {
	setupEnv.Do(base.SetupSuite)
}

// UniqueTitle создаёт короткое уникальное название для каталога.
func UniqueTitle() string {
	return "QA " + strings.ToLower(rand.Text())
}

// AdForm готовит корректную форму с указанным названием и количеством PNG.
func AdForm(t *testing.T, title string, photoCount int) advertisements.Form {
	t.Helper()
	form := advertisements.Form{
		Fields: map[string]string{
			TitleField:       title,
			DescriptionField: "Временное объявление для API-тестов",
			PriceField:       fmt.Sprintf("%.2f", Price),
			QuantityField:    fmt.Sprint(Quantity),
		},
	}
	for i := range photoCount {
		img := image.NewPaletted(image.Rect(0, 0, 32, 32), color.Palette{color.NRGBA{R: 32, G: 140, B: 210, A: 255}, color.White})
		img.SetColorIndex(i, i, 1)
		var photo bytes.Buffer
		require.NoError(t, png.Encode(&photo, img), "Не удалось подготовить PNG")
		form.Photos = append(form.Photos, photo.Bytes())
	}
	return form
}

// CreateAd создаёт карточку из формы и регистрирует её удаление после теста.
func CreateAd(t *testing.T, token string, form advertisements.Form) *Ad {
	t.Helper()
	item := advertisements.Create(t, token, form)
	require.NotEmpty(t, item.ID, "POST /advertisement вернул пустой ID")
	ad := &Ad{Advertisement: item, token: token}
	t.Cleanup(func() {
		if !ad.deleted {
			ad.Delete(t)
		}
	})
	t.Logf("Получен ID созданного объявления: %s", ad.ID)
	return ad
}

// NewAd создаёт временное объявление с уникальным названием и одним фото.
func NewAd(t *testing.T, token string) *Ad {
	t.Helper()
	return CreateAd(t, token, AdForm(t, UniqueTitle(), 1))
}

// Delete удаляет карточку и отмечает, что повторная очистка не нужна.
func (ad *Ad) Delete(t *testing.T) {
	t.Helper()
	t.Logf("Удаляем тестовое объявление %s", ad.ID)
	advertisements.Delete(t, ad.token, ad.ID)
	ad.deleted = true
}
