package fixtures

import (
	"testing"

	"github.com/stretchr/testify/require"

	apiResponse "api-tests-template/internal/helpers/api-response"
)

type UnauthorizedCase struct {
	Name    string
	Token   string
	Message string
}

// UnauthorizedCases задаёт общие варианты проверки отсутствующего и неверного токена.
func UnauthorizedCases() []UnauthorizedCase {
	return []UnauthorizedCase{
		{Name: "без токена", Message: "Authorization header required"},
		{Name: "с неверным токеном", Token: "incorrect_token", Message: "Invalid or expired token"},
	}
}

// AssertError проверяет машинный код и текст ошибки, которые вернул API.
func AssertError(t *testing.T, actual apiResponse.Error, code, message string) {
	t.Helper()
	require.Equal(t, apiResponse.Error{Code: code, Message: message}, actual, "Код или сообщение об ошибке не совпали")
	t.Logf("Проверен ответ об ошибке: %s: %s", code, message)
}
