package flag

import (
	"github.com/machtwatch/catalystdk/go/flag"
	"github.com/machtwatch/catalystdk/go/flag/config"
	flagprovider "github.com/machtwatch/catalystdk/go/flag/provider"
)

func NewFlag(config config.Config) *flagprovider.Flag {
	return flag.New(config)
}
