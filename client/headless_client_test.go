package client

import (
	"testing"

	"github.com/cention-mujibur-rahman/msgoraph/scopes"
)

func TestHeadlessClientInitialization(t *testing.T) {
	applicationID := ""
	applicationSecret := ""
	c := NewHeadless(applicationID, applicationSecret, scopes.All(scopes.PermissionTypeApplication))
	err := c.InitializeCredentials()
	if err != nil {
		t.Fatalf(err.Error())
	}
}
