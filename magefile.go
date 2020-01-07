// +build mage

package main

import (
	"fmt"
	"os"
	"strings"

	// "github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	s, err := sh.Output("go", "list", "./...")
	if err != nil {
		return err
	}
	pkgs := strings.Split(s, "\n")
	for i := range pkgs {
		pkgs[i] = strings.TrimPrefix(pkgs[i], "_")
		BuildGo(pkgs[i])
	}
	return nil
}

func BuildGo(filepath string) error {
	os.Chdir(filepath)
	os.Remove("main")
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return err
	}
	if err := sh.Run("gofmt", "-w", "-s", "."); err != nil {
		return err
	}
	ldflags := "-s -w"
	if err := sh.Run("go", "build", "-ldflags", ldflags, "-o", "main", "."); err != nil {
		return err
	}
	fmt.Printf("built %v\n", filepath)
	return nil
}

func BuildFrontend() error {
	os.Chdir("services/golib")
	return sh.Run("npm", "run", "build")
}
