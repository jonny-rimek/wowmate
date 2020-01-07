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

func Build() error {
	// mg.SerialDeps(Go, BuildFrontend)
	if err := Go(); err != nil {
		return err
	}
	if err := Frontend(); err != nil {
		return err
	}
	if err := CDK(); err != nil {
		return err
	}

	return nil
}

func Go() error {
	s, err := sh.Output("go", "list", "./...")
	if err != nil {
		return err
	}
	pkgs := strings.Split(s, "\n")
	for i := range pkgs {
		filepath := strings.TrimPrefix(pkgs[i], "_")
		os.Chdir(filepath)
		os.Remove("main")
		
		//for some reason go list ./... finds files in cdk.out
		if strings.Contains(filepath, "cdk.out") == true {
			continue;
		}
		if err := GoTidy(); err != nil {
			return err
		}
		if err := Gofmt(); err != nil {
			return err
		}
		if err := GoBuild(); err != nil {
			return err
		}
		fmt.Printf("built %v\n", filepath)
	}
	return nil
}

func GoTidy() error {
	return sh.Run("go", "mod", "tidy")
}

func Gofmt() error {
	return sh.Run("gofmt", "-w", "-s", ".")
}

func GoBuild() error {
	return sh.Run("go", "build", "-ldflags", "-s -w", "-o", "main", ".")
}

func Frontend() error {
	os.Chdir("/home/jimbo/dev/wowmate/services/frontend")
	if err := sh.Run("npm", "install"); err != nil {
		return err
	}

	return sh.Run("npm", "run", "build")
}

func CDK() error {
	os.Chdir("/home/jimbo/dev/wowmate")
	if err := sh.Run("npm", "install"); err != nil {
		return err
	}
	if err := sh.Run("tsc"); err != nil {
		return err
	}

	return sh.Run("cdk", "deploy", "--require-approval=never")
}