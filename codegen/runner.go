package codegen

import (
	"fmt"
	"os"
	"os/exec"
)

func RunTerraformHack(fileName string) error {
	cmd := exec.Command("go", "run", fileName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error starting hack: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Error waiting for hack: %s", err)
	}

	return nil
}