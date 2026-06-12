package agent

import (
	"github.com/mauroociappinaph/ayrton/pkg/agent/roles"
	"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type Factory func(name string, broadcaster shared.BroadcasterInterface) Agent

var factories = map[string]Factory{
	"pm":          func(name string, b shared.BroadcasterInterface) Agent { return roles.NewPM(name, b) },
	"dev":         func(name string, b shared.BroadcasterInterface) Agent { return roles.NewDev(name, b) },
	"marketing":   func(name string, b shared.BroadcasterInterface) Agent { return roles.NewMarketing(name, b) },
	"ops":         func(name string, b shared.BroadcasterInterface) Agent { return roles.NewOps(name, b) },
	"prospeccion": func(name string, b shared.BroadcasterInterface) Agent { return roles.NewProspeccion(name, b) },
	"auditor":     func(name string, b shared.BroadcasterInterface) Agent { return nil },
	"learning":    func(name string, b shared.BroadcasterInterface) Agent { return nil },
}

func CreateAgent(role, name string, broadcaster shared.BroadcasterInterface) (Agent, bool) {
	if factory, ok := factories[role]; ok {
		if a := factory(name, broadcaster); a != nil {
			return a, true
		}
	}
	return nil, false
}

func RegisterFactory(role string, factory Factory) {
	factories[role] = factory
}