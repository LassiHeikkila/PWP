package main

import (
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func makeTask(task *types.Task) func() {
	return func() {}
}
