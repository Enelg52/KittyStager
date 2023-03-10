package test

import (
	"KittyStager/internal/task/priv"
	"github.com/go-playground/assert/v2"
	"testing"
)

var (
	privName        string
	privDescription string
	privEnabled     bool
	privIntegrity   string
	privList        []*priv.Privilege
)

func init() {
	privName = "test"
	privDescription = "descriptions"
	privIntegrity = "medium"
	privEnabled = true
}

func TestGetPriv(t *testing.T) {
	t.Parallel()
	p := priv.NewPrivilege(privName, privDescription, privEnabled)
	assert.Equal(t, p.GetName(), privName)
	assert.Equal(t, p.GetDescription(), privDescription)
	assert.Equal(t, p.GetEnable(), privEnabled)
}

func TestMarshallUnmashall(t *testing.T) {
	t.Parallel()
	p := priv.NewPrivilege(privName, privDescription, privEnabled)
	privList = append(privList, p)
	pr := priv.NewPrivileges(privList, privIntegrity)
	// when
	b, err := pr.MarshallPrivileges()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = pr.UnmarshallPrivileges(b)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}
