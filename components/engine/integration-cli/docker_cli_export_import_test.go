package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// export an image and try to import it into a new one
func TestExportContainerAndImportImage(t *testing.T) {
	runCmd := exec.Command(dockerBinary, "run", "-d", "busybox", "true")
	out, _, err := runCommandWithOutput(runCmd)
	if err != nil {
		t.Fatal("failed to create a container", out, err)
	}

	cleanedContainerID := stripTrailingCharacters(out)

	inspectCmd := exec.Command(dockerBinary, "inspect", cleanedContainerID)
	out, _, err = runCommandWithOutput(inspectCmd)
	if err != nil {
		t.Fatalf("output should've been a container id: %s %s ", cleanedContainerID, err)
	}

	exportCmd := exec.Command(dockerBinary, "export", cleanedContainerID)
	if out, _, err = runCommandWithOutput(exportCmd); err != nil {
		t.Fatalf("failed to export container: %s, %v", out, err)
	}

	importCmd := exec.Command(dockerBinary, "import", "-", "repo/testexp:v1")
	importCmd.Stdin = strings.NewReader(out)
	out, _, err = runCommandWithOutput(importCmd)
	if err != nil {
		t.Fatalf("failed to import image: %s, %v", out, err)
	}

	cleanedImageID := stripTrailingCharacters(out)

	inspectCmd = exec.Command(dockerBinary, "inspect", cleanedImageID)
	if out, _, err = runCommandWithOutput(inspectCmd); err != nil {
		t.Fatalf("output should've been an image id: %s, %v", out, err)
	}

	deleteContainer(cleanedContainerID)
	deleteImages("repo/testexp:v1")

	logDone("export - export a container")
	logDone("import - import an image")
}

// Used to test output flag in the export command
func TestExportContainerWithOutputAndImportImage(t *testing.T) {
	runCmd := exec.Command(dockerBinary, "run", "-d", "busybox", "true")
	out, _, err := runCommandWithOutput(runCmd)
	if err != nil {
		t.Fatal("failed to create a container", out, err)
	}

	cleanedContainerID := stripTrailingCharacters(out)

	inspectCmd := exec.Command(dockerBinary, "inspect", cleanedContainerID)
	out, _, err = runCommandWithOutput(inspectCmd)
	if err != nil {
		t.Fatalf("output should've been a container id: %s %s ", cleanedContainerID, err)
	}

	exportCmdTemplate := `%v export --output=/tmp/testexp.tar %v`
	exportCmdFinal := fmt.Sprintf(exportCmdTemplate, dockerBinary, cleanedContainerID)
	exportCmd := exec.Command(exportCmdFinal)
	if out, _, err = runCommandWithOutput(exportCmd); err != nil {
		t.Fatalf("failed to export container: %s, %v", out, err)
	}

	importCmdFinal := `cat /tmp/testexp.tar | docker import - repo/testexp:v1`
	importCmd := exec.Command(importCmdFinal)
	out, _, err = runCommandWithOutput(importCmd)
	if err != nil {
		t.Fatalf("failed to import image: %s, %v", out, err)
	}

	cleanedImageID := stripTrailingCharacters(out)

	inspectCmd = exec.Command(dockerBinary, "inspect", cleanedImageID)
	if out, _, err = runCommandWithOutput(inspectCmd); err != nil {
		t.Fatalf("output should've been an image id: %s, %v", out, err)
	}

	deleteContainer(cleanedContainerID)
	deleteImages("repo/testexp:v1")

	os.Remove("/tmp/testexp.tar")

	logDone("export - export a container with output flag")
	logDone("import - import an image with output flag")
}
