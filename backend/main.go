// SPDX-License-Identifier: MIT
package main

import (
	"github.com/ENIACSystems/FileENIAC/backend/cmd"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
)

func main() {
	defer log.Sync()
	cmd.Execute()
}
