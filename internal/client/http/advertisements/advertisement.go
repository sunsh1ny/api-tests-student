package advertisements

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

type Form struct {
	Fields map[string]string
	Photos [][]byte
}

// Create отправляет поля и PNG-файлы объявления одним multipart-запросом.
func Create(t *testing.T, token string, form Form) *http.Response {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for name, value := range form.Fields {
		require.NoError(t, writer.WriteField(name, value), "Не удалось записать поле %s", name)
	}
	for _, photo := range form.Photos {
		header := textproto.MIMEHeader{}
		header.Set("Content-Disposition", `form-data; name="photos"; filename="fixture.png"`)
		header.Set("Content-Type", "image/png")
		part, err := writer.CreatePart(header)
		require.NoError(t, err, "Не удалось добавить фото в multipart")
		_, err = part.Write(photo)
		require.NoError(t, err, "Не удалось записать фото")
	}
	require.NoError(t, writer.Close(), "Не удалось завершить multipart")
	return apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		ContentType(writer.FormDataContentType()).Body(body.String()).
		Expect(t).End().Response
}

// Get запрашивает карточку объявления по его идентификатору.
func Get(t *testing.T, token, adID string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", adID).Expect(t).End().Response
}

// Delete удаляет карточку объявления по её идентификатору.
func Delete(t *testing.T, token, adID string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", adID).Expect(t).End().Response
}

// Photos запрашивает список фотографий объявления.
func Photos(t *testing.T, token, adID string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath + "/" + url.PathEscape(adID) + "/photos").
		Expect(t).End().Response
}
