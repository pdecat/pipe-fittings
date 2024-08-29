package connection

import (
	"context"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

type SlackConnection struct {
	ConnectionImpl

	Token *string `json:"token,omitempty" cty:"token" hcl:"token,optional"`
}

func (c *SlackConnection) getEnv() map[string]cty.Value {
	env := map[string]cty.Value{}
	if c.Token != nil {
		env["SLACK_TOKEN"] = cty.StringVal(*c.Token)
	}
	return env
}

func (c *SlackConnection) CtyValue() (cty.Value, error) {
	ctyValue, err := modconfig.GetCtyValue(c)
	if err != nil {
		return cty.NilVal, err
	}

	valueMap := ctyValue.AsValueMap()
	valueMap["env"] = cty.ObjectVal(c.getEnv())

	return cty.ObjectVal(valueMap), nil
}

func (c *SlackConnection) Equals(otherConnection PipelingConnection) bool {
	// If both pointers are nil, they are considered equal
	if c == nil && helpers.IsNil(otherConnection) {
		return true
	}

	if (c == nil && !helpers.IsNil(otherConnection)) || (c != nil && helpers.IsNil(otherConnection)) {
		return false
	}

	other, ok := otherConnection.(*SlackConnection)
	if !ok {
		return false
	}

	if !utils.PtrEqual(c.Token, other.Token) {
		return false
	}

	return true
}

func (c *SlackConnection) Resolve(ctx context.Context) (PipelingConnection, error) {
	if c.Token == nil {
		slackTokenEnvVar := os.Getenv("SLACK_TOKEN")

		// Don't modify existing credential, resolve to a new one
		newCreds := &SlackConnection{
			ConnectionImpl: c.ConnectionImpl,
			Token:          &slackTokenEnvVar,
		}

		return newCreds, nil
	}
	return c, nil
}

func (c *SlackConnection) GetTtl() int {
	return -1
}

func (c *SlackConnection) Validate() hcl.Diagnostics {
	return hcl.Diagnostics{}
}

type SlackConnectionConfig struct {
	Token *string `cty:"token" hcl:"token"`
}

func (c *SlackConnectionConfig) GetConnection(name string, shortName string) PipelingConnection {

	slackCred := &SlackConnection{
		ConnectionImpl: ConnectionImpl{
			HclResourceImpl: modconfig.HclResourceImpl{
				FullName:        name,
				ShortName:       shortName,
				UnqualifiedName: name,
			},
			Type: "slack",
		},

		Token: c.Token,
	}

	return slackCred
}
