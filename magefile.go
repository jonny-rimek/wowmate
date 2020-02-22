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

var Aliases = map[string]interface{} {
	"d": Deploy,
	"di": Diff,
	"b": Build,
	"g": Go,
	"f": Frontend,
}

func rootDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	os.Chdir(home)
	os.Chdir("dev/wowmate")
	return nil
}

func Build() error {
	mg.SerialDeps(Go, Frontend)

	return nil
}
/* 
TODO: add a scaffold lambda function
		- [ ] run npm install if dir is missing
		- [ ] add command to update
		- [ ] create folder and go file with the same name
		- [ ] go mod init
		- [ ] go mod edit -replace=github.com/alexedwards/argon2id=/home/alex/code/argon2i
		- [ ] add biler plate go code to mail file, including golib.InitLogging()
 */

func Go() error {
	if err := rootDir(); err != nil {
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
	os.Remove("main")
	return sh.Run("go", "build", "-ldflags", "-s -w", "-o", "main", ".")
}

func Frontend() error {
	if err := rootDir(); err != nil {
		return err
	}
	os.Chdir("services/frontend")

	if _, err := os.Stat("node_modules"); err != nil {
		if os.IsNotExist(err) {
			sh.Run("npm", "install")
		} else {
			return err
		}
	}
	return sh.Run("npm", "run", "build")
}

func Deploy() error {
	if err := rootDir(); err != nil {
		return err
	}
	// if err := sh.Run("npm", "install"); err != nil {
	// 	return err
	// }
	if err := sh.Run("tsc"); err != nil {
		return err
	}

	return sh.Run("cdk", "deploy", "--require-approval=never")
}

func Diff() error {
	if err := rootDir(); err != nil {
		return err
	}
	if _, err := os.Stat("node_modules"); err != nil {
		if os.IsNotExist(err) {
			sh.Run("npm", "install")
		} else {
			return err
		}
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
	if err := rootDir(); err != nil {
		return err
	}
	os.Chdir("services/frontend")
	sh.Run("rm", "-rf", "node_modules");
	return sh.Run("rm", "-rf", "dist");
}

func clearCDK() error {
	if err := rootDir(); err != nil {
		return err
	}
	return sh.Run("rm", "-rf", "node_modules");
}

func clearGo() error {
	if err := rootDir(); err != nil {
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
		fmt.Println(filepath)

		sh.Run("rm", "main");
	}

	return nil
}
