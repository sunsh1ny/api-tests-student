package favorites

import (
	"testing"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesAddAndRemove(t *testing.T) {
	owner, actor := loginAccounts(t)
	ad := fixtures.NewAd(t, owner)

	addFavorite(t, actor, ad.ID)
	checkFavorite(t, actor, ad.ID, true)

	t.Log("Удаляем объявление из избранного через DELETE /favorites/{ad_id}")
	favorites.Delete(t, actor, ad.ID)
	checkFavorite(t, actor, ad.ID, false)
}

func TestFavoritesRemoveDeletedAdvertisement(t *testing.T) {
	owner, actor := loginAccounts(t)
	ad := fixtures.NewAd(t, owner)

	addFavorite(t, actor, ad.ID)
	checkFavorite(t, actor, ad.ID, true)

	t.Log("Владелец удаляет само объявление через DELETE /advertisement")
	ad.Delete(t)
	checkFavorite(t, actor, ad.ID, false)
}
