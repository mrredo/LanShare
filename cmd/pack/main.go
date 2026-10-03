package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

/*
This packs the frontend, builds it so the app can be packaged into a standalone executable file.
*/
func main() {
	root, err := projectRoot()
	if err != nil {
		panic(err)
	}

	frontendDir := filepath.Join(root, "lanshare-fe")
	frontendDist := filepath.Join(frontendDir, "dist")
	embedDir := filepath.Join(root, "cmd", "lanshare", "web", "dist")
	buildDir := filepath.Join(root, "build")

	if err := run(frontendDir, "npm", "run", "build"); err != nil {
		panic(err)
	}

	if err := os.RemoveAll(embedDir); err != nil {
		panic(err)
	}

	if err := copyDir(frontendDist, embedDir); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		panic(err)
	}

	output := filepath.Join(buildDir, "lanshare")
	if runtime.GOOS == "windows" {
		output += ".exe"
	}

	if err := run(root, "go", "build", "-o", output, "./cmd/lanshare"); err != nil {
		panic(err)
	}

	fmt.Printf("Built LanShare: %s\n", output)
}

func projectRoot() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	root, err := filepath.Abs(filepath.Join(workingDir, "..", ".."))
	if err != nil {
		return "", err
	}

	return root, nil
}

func run(dir string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func copyDir(source, destination string) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		target := filepath.Join(destination, relative)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		input, err := os.Open(path)
		if err != nil {
			return err
		}
		defer input.Close()

		output, err := os.Create(target)
		if err != nil {
			return err
		}
		defer output.Close()

		_, err = io.Copy(output, input)
		return err
	})
}
