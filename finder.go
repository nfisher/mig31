package main

import (
	"github.com/hailocab/mig31/config"
	"github.com/hailocab/mig31/runtime"
	"github.com/hailocab/mig31/set"
)

// FindAppliedSet find the set of migrations that are currently applied.
func FindAppliedSet(rtConfig *runtime.Config, env *config.Environment) (appliedSet set.Set) {
	if rtConfig.Offline {
		appliedSet = set.New()
		return
	}
	return
}
