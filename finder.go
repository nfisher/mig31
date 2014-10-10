package main

import (
	"github.com/hailocab/mig31/config"
	"github.com/hailocab/mig31/set"
)

// FindAppliedSet find the set of migrations that are currently applied.
func FindAppliedSet(env *config.Environment) (appliedSet set.Set) {
	// host is blank and is therefore an offline request
	if env.Host == "" {
		appliedSet = set.New()
		return
	}
	return
}
