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

const ProjectPath = "/dev/wowmate/"

var Aliases = map[string]interface{} {
	"d": Deploy,
	"di": Diff,
	"b": Build,
	"g": Go,
	"f": Frontend,
}

/* 
TODO: add a scaffold lambda function
		- [x] run npm install if dir is missing
		- [ ] add command to update //go get -u ./...
		- [ ] command to delete all cache
			- gp/pkg
			- ~/.cache/go-build
			- node_modules
			- services/frontend/node_modules

IMPROVE:
		- [ ] create folder and go file with the same name (?)
		- [ ] go mod init
		- [ ] go mod edit -replace=github.com/alexedwards/argon2id=/home/alex/code/argon2i
		- [ ] add boiler plate go code to mail file, including golib.InitLogging()
 */

func Go() error {
	if err := projectDir(); err != nil {
		return err
	}

	s, err := sh.Output("go", "list", "./...")
	if err != nil {
		return err
	}
	pkgs := strings.Split(s, "\n")
	for i := range pkgs {
		filepath := strings.TrimPrefix(pkgs[i], "_")
		
		//for some reason go list ./... finds files in cdk.out
		if strings.Contains(filepath, "cdk.out") == true {
			continue;
		}
		os.Chdir(filepath)

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
	return sh.Run("go", "build", "-ldflags", "-s -w", "-o", "main", ".")
}

func Frontend() error {
	if err := projectSubDir("services/frontend"); err != nil {
		return err
	}
	if err := yarnInstall(); err != nil {
		return err
	}
	return sh.Run("yarn", "generate")
}

func Deploy() error {
	if err := projectDir(); err != nil {
		return err
	}
	if err := npmInstall(); err != nil {
		return err
	}
	if err := sh.Run("tsc"); err != nil {
		return err
	}
	return sh.Run("cdk", "deploy", "--require-approval=never")
}

func yarnInstall() error {
	if _, err := os.Stat("node_modules"); err != nil {
		if os.IsNotExist(err) {
			sh.Run("yarn", "install")
		} else {
			return err
		}
	}
	return nil
}

func npmInstall() error {
	if _, err := os.Stat("node_modules"); err != nil {
		if os.IsNotExist(err) {
			sh.Run("npm", "install")
		} else {
			return err
		}
	}
	return nil
}

func Diff() error {
	if err := projectDir(); err != nil {
		return err
	}
	if err := npmInstall(); err != nil {
		return err
	}
	if err := sh.Run("tsc"); err != nil {
		return err
	}

	return sh.Run("cdk", "diff")
}

func Clear() error {
	mg.SerialDeps(clearGo, clearFrontend, clearCDK)

	return nil
}

func clearFrontend() error {
	if err := projectSubDir("services/frontend"); err != nil {
		return err
	}
	if err := sh.Run("rm", "-rf", "node_modules"); err != nil {
		return err
	}
	return sh.Run("rm", "-rf", "dist");
}

func clearCDK() error {
	if err := projectDir(); err != nil {
		return err
	}
	return sh.Run("rm", "-rf", "node_modules");
}

func clearGo() error {
	if err := projectDir(); err != nil {
		return err
	}

	s, err := sh.Output("go", "list", "./...")
	if err != nil {
		return err
	}
	pkgs := strings.Split(s, "\n")
	for i := range pkgs {
		filepath := strings.TrimPrefix(pkgs[i], "_")
		
		//for some reason go list ./... finds files in cdk.out
		if strings.Contains(filepath, "cdk.out") == true {
			continue;
		}
		if err = os.Chdir(filepath); err != nil {
			return err
		}
		if err = sh.Run("rm", "main"); err != nil {
			return err
		}
	}
	return nil
}

func projectSubDir(subPath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if err = os.Chdir(home + ProjectPath + subPath); err != nil {
		return err
	}

	return nil
}

func projectDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if err = os.Chdir(home + ProjectPath); err != nil {
		return err
	}

	return nil
}

func Build() error {
	mg.SerialDeps(Go, Frontend)

	return nil
}