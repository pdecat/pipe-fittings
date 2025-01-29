package modconfig

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/options"
)

type OptionsBlockFactory = func(*hcl.Block) (options.Options, hcl.Diagnostics)
