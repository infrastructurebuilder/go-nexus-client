package nexus3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityAnonymous(t *testing.T) {
	client := getTestClient()

	oldAnonymous, err := client.Security.Anonymous.Read()
	assert.Nil(t, err)
	assert.NotNil(t, oldAnonymous)

	newAnonymous := SecurityAnonymousAccessSettings{
		Enabled:   true,
		UserID:    "anonymous",
		RealmName: "NexusAuthorizingRealm",
	}
	err = client.Security.Anonymous.Update(newAnonymous)
	assert.Nil(t, err)

	anonymous, err := client.Security.Anonymous.Read()
	assert.Nil(t, err)
	assert.NotNil(t, anonymous)
	assert.Equal(t, *anonymous, newAnonymous)

	err = client.Security.Anonymous.Update(*oldAnonymous)
	assert.Nil(t, err)
}
