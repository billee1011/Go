package user

import (
	"testing"

	"github.com/Sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func Test_verifyIDCard(t *testing.T) {
	assert.True(t, verifyIDCard("410322199202152910"))
	assert.True(t, verifyIDCard("522225197612176132"))
	assert.True(t, verifyIDCard("441723199407166317"))
	assert.True(t, verifyIDCard("542128197911042342"))
	assert.True(t, verifyIDCard("469026198312079605"))
	assert.True(t, verifyIDCard("22032319930425218X"))

	assert.False(t, verifyIDCard("41032219920215291X"))
	assert.False(t, verifyIDCard("52222519761217613X"))
	assert.False(t, verifyIDCard("44172319940716631X"))
	assert.False(t, verifyIDCard("54212819791104234X"))
	assert.False(t, verifyIDCard("46902619831207960X"))
	assert.False(t, verifyIDCard("220323199304252181"))

	assert.False(t, verifyIDCard("毛主席万岁"))
	assert.False(t, verifyIDCard("469026198312079大"))
	assert.False(t, verifyIDCard("小026198312079000"))
}

func Test_verifyName(t *testing.T) {
	assert.True(t, verifyName("毛泽东"))
	assert.True(t, verifyName("周恩来"))
	assert.True(t, verifyName("史蒂夫"))
	assert.True(t, verifyName("小朱"))
	assert.True(t, verifyName("玙"))
	assert.True(t, verifyName("玙王"))

	assert.False(t, verifyName("玙/王"))
	assert.False(t, verifyName("22王"))
	assert.False(t, verifyName("ab王"))
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}
