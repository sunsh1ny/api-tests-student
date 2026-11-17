package advertisementCreation

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

const (
	titleRequired       = "title is required"
	descriptionRequired = "description is required"
	invalidPrice        = "price must be between 0 and 1000000 inclusive"
	invalidQuantity     = "quantity must be between 1 and 100 inclusive"
	titleLimit          = 50
	descriptionLimit    = 500
)

func TestCreateAdvertisementRejectsInvalidData(t *testing.T) {
	owner := fixtures.Login(t, false)
	cases := []struct {
		name, field, value, message string
		omit                        bool
	}{
		{name: "без названия", field: fixtures.TitleField, omit: true, message: titleRequired},
		{name: "пустое название", field: fixtures.TitleField, message: titleRequired},
		{name: "длинное название", field: fixtures.TitleField, value: strings.Repeat("x", titleLimit+1), message: fmt.Sprintf("title must not exceed %d characters provided %d", titleLimit, titleLimit+1)},
		{name: "без описания", field: fixtures.DescriptionField, omit: true, message: descriptionRequired},
		{name: "пустое описание", field: fixtures.DescriptionField, message: descriptionRequired},
		{name: "длинное описание", field: fixtures.DescriptionField, value: strings.Repeat("x", descriptionLimit+1), message: "description must be between 1 and 500 characters"},
		{name: "без фото", message: "at least 1 photo is required"},
		{name: "отрицательная цена", field: fixtures.PriceField, value: "-1", message: invalidPrice},
		{name: "цена выше максимума", field: fixtures.PriceField, value: "1000001", message: invalidPrice},
		{name: "нулевое количество", field: fixtures.QuantityField, value: "0", message: invalidQuantity},
		{name: "количество выше максимума", field: fixtures.QuantityField, value: "101", message: invalidQuantity},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			form := fixtures.AdForm(t, fixtures.UniqueTitle(), 1)
			switch {
			case tc.field == "":
				form.Photos = nil
			case tc.omit:
				delete(form.Fields, tc.field)
			default:
				form.Fields[tc.field] = tc.value
			}
			t.Logf("Отправляем POST /advertisement: %s", tc.name)
			actual := advertisements.CreateError(t, owner.Token, form, http.StatusBadRequest)
			fixtures.AssertError(t, actual, "invalid_request", tc.message)
		})
	}
}

func TestCreateAdvertisementRequiresAuthorization(t *testing.T) {
	for _, tc := range fixtures.UnauthorizedCases() {
		t.Run(tc.Name, func(t *testing.T) {
			fixtures.LoadEnv()
			form := fixtures.AdForm(t, fixtures.UniqueTitle(), 1)
			t.Logf("Отправляем корректный POST /advertisement %s", tc.Name)
			actual := advertisements.CreateError(t, tc.Token, form, http.StatusUnauthorized)
			fixtures.AssertError(t, actual, "unauthorized", tc.Message)
		})
	}
}
