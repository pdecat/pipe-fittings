package credential

import (
	"context"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/zclconf/go-cty/cty"
)

type AbuseIPDBCredential struct {
	CredentialImpl

	APIKey *string `json:"api_key,omitempty" cty:"api_key" hcl:"api_key,optional"`
}

func (c *AbuseIPDBCredential) getEnv() map[string]cty.Value {
	// There is no environment variable listed in the AbuseIPDB official API docs
	// https://www.abuseipdb.com/api.html
	return nil
}

func (c *AbuseIPDBCredential) CtyValue() (cty.Value, error) {
	return ctyValueForCredential(c)
}

func (c *AbuseIPDBCredential) Equals(otherCredential Credential) bool {
	// If both pointers are nil, they are considered equal
	if c == nil && helpers.IsNil(otherCredential) {
		return true
	}

	if (c == nil && !helpers.IsNil(otherCredential)) || (c != nil && helpers.IsNil(otherCredential)) {
		return false
	}

	other, ok := otherCredential.(*AbuseIPDBCredential)
	if !ok {
		return false
	}

	if !utils.PtrEqual(c.APIKey, other.APIKey) {
		return false
	}

	return true
}

func (c *AbuseIPDBCredential) Resolve(ctx context.Context) (Credential, error) {
	if c.APIKey == nil {
		abuseIPDBAPIKeyEnvVar := os.Getenv("ABUSEIPDB_API_KEY")

		// Don't modify existing credential, resolve to a new one
		newCreds := &AbuseIPDBCredential{
			CredentialImpl: c.CredentialImpl,
			APIKey:         &abuseIPDBAPIKeyEnvVar,
		}

		return newCreds, nil
	}
	return c, nil
}

func (c *AbuseIPDBCredential) GetTtl() int {
	return -1
}

func (c *AbuseIPDBCredential) Validate() hcl.Diagnostics {
	return hcl.Diagnostics{}
}

type AbuseIPDBConnectionConfig struct {
	APIKey *string `cty:"api_key" hcl:"api_key"`
}

func (c *AbuseIPDBConnectionConfig) GetCredential(name string, shortName string) Credential {

	abuseIPDBCred := &AbuseIPDBCredential{
		CredentialImpl: CredentialImpl{
			HclResourceImpl: modconfig.HclResourceImpl{
				FullName:        name,
				ShortName:       shortName,
				UnqualifiedName: name,
			},
			Type: "abuseipdb",
		},

		APIKey: c.APIKey,
	}

	return abuseIPDBCred
}
