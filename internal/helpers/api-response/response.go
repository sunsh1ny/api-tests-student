package api_response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type Error struct {
	Code    string `json:"error"`
	Message string `json:"message"`
}

// CheckStatus проверяет HTTP-статус и показывает тело ответа при ошибке.
func CheckStatus(t *testing.T, response *http.Response, status int) {
	t.Helper()
	require.NotNil(t, response, "HTTP-запрос не вернул ответ")
	if response.StatusCode != status {
		body, err := io.ReadAll(io.LimitReader(response.Body, 1024))
		require.NoError(t, err, "Не удалось прочитать тело ответа")
		require.Equal(t, status, response.StatusCode, "Неожиданный HTTP-статус, тело ответа: %s", body)
	}
}

// Decode проверяет статус и читает JSON в указанную модель.
func Decode[T any](t *testing.T, response *http.Response, status int) T {
	t.Helper()
	require.NotNil(t, response, "HTTP-запрос не вернул ответ")
	defer Close(t, response)
	CheckStatus(t, response, status)
	var result T
	require.NoError(t, json.NewDecoder(response.Body).Decode(&result), "Некорректный JSON в HTTP-ответе")
	return result
}

// Close проверяет, что тело ответа удалось закрыть.
func Close(t *testing.T, response *http.Response) {
	t.Helper()
	require.NotNil(t, response, "HTTP-запрос не вернул ответ")
	require.NoError(t, response.Body.Close(), "Не удалось закрыть тело HTTP-ответа")
}
