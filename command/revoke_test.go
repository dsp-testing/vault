package command

import (
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestRevoke(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RevokeCommand{
		Meta: Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	client := testClient(t, addr, token)
	err := client.Logical().Write("secret/foo", map[string]interface{}{
		"key":   "value",
		"lease": "1m",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	secret, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	args := []string{
		"-address", addr,
		secret.VaultId,
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
