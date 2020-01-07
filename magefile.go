// +build mage

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() error {
	mg.SerialDeps(Go, Frontend)

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
		
		//for some reason go list ./... finds files in cdk.out
		if strings.Contains(filepath, "cdk.out") == true {
			continue;
		}
		if err := goTidy(); err != nil {
			return err
		}
		if err := gofmt(); err != nil {
			return err
		}
		if err := goBuild(); err != nil {
			return err
		}
		fmt.Printf("built %v\n", filepath)
	}
	return nil
}

func goTidy() error {
	return sh.Run("go", "mod", "tidy")
}

func gofmt() error {
	return sh.Run("gofmt", "-w", "-s", ".")
}

func goBuild() error {
	os.Remove("main")
	return sh.Run("go", "build", "-ldflags", "-s -w", "-o", "main", ".")
}

func Frontend() error {
	os.Chdir("/home/jimbo/dev/wowmate/services/frontend")
	if err := sh.Run("npm", "install"); err != nil {
		return err
	}

	return sh.Run("npm", "run", "build")
}

func Deploy() error {
	mg.SerialDeps(Go, Frontend)
	os.Chdir("/home/jimbo/dev/wowmate")
	if err := sh.Run("npm", "install"); err != nil {
		return err
	}
	if err := sh.Run("tsc"); err != nil {
		return err
	}

	return sh.Run("cdk", "deploy", "--require-approval=never")
}