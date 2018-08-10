package charge

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	itemData, err := ioutil.ReadFile("./itemlist.json")
	if err != nil {
		panic(err)
	}
	itemListCfg := string(itemData)
	configGetter = func(key, subkey string) (string, error) {
		return itemListCfg, nil
	}
}

// Test_loadItemList test function loadItemList
func Test_loadItemList(t *testing.T) {
	assert.Nil(t, loadItemList())
	assert.Len(t, itemLists, 2)
	android := itemLists["android"]
	assert.Equal(t, 2, len(android))
	defaultCityItems := android["default"]
	assert.Len(t, defaultCityItems, 2)
	firstItem := defaultCityItems[0]
	assert.Equal(t, uint64(1), firstItem.ID)
	assert.Equal(t, firstItem.Coin, 100)
	assert.Equal(t, firstItem.Price, 600)
	assert.Equal(t, firstItem.Name, "金豆 100")
	assert.Equal(t, firstItem.Tag, "热卖")
	assert.Equal(t, firstItem.PresentCoin, 0)
}

func Test_getItemList(t *testing.T) {
	assert.Nil(t, loadItemList())
	items, err := getItemList(0, 2) // iphone, default city
	assert.Nil(t, err)
	assert.Len(t, items, 2)
	firstItem := items[0]
	assert.Equal(t, uint64(1), firstItem.ID)
	assert.Equal(t, firstItem.Coin, 100)
	assert.Equal(t, firstItem.Price, 600)
	assert.Equal(t, firstItem.Name, "金豆 100")
	assert.Equal(t, firstItem.Tag, "热卖")
	assert.Equal(t, firstItem.PresentCoin, 0)
}
