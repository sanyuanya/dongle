package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {

	platforms := []string{"linux/amd64", "linux/arm64", "darwin/amd64", "darwin/arm64", "windows/amd64"}

	// 程序名称
	program := "dongle-program"

	// 获取当前环境变量
	env := os.Environ()

	// 添加 GOPATH 和 GOMODCACHE 到环境变量
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = os.Getenv("HOME") + "/go"
	}

	gomodcache := os.Getenv("GOMODCACHE")
	if gomodcache == "" {
		gomodcache = gopath + "/pkg/mod"
	}

	env = append(env, "GOPATH="+gopath, "GOMODCACHE="+gomodcache)

	for _, platform := range platforms {
		split := strings.Split(platform, "/")
		GOOS := split[0]
		GOARCH := split[1]
		outputName := fmt.Sprintf("build/%s-%s-%s", program, GOOS, GOARCH)
		if GOOS == "windows" {
			outputName += ".exe"
		}

		fmt.Printf("Building for %s/%s...\n", GOOS, GOARCH)
		fmt.Printf("GOOS=%s GOARCH=%s go build -o %s\n", GOOS, GOARCH, outputName)
		cmd := exec.Command("go", "build", "-o", outputName, "/Users/sanyuanya/hjworkspace/go_dev/dongle_new/main.go")
		cmd.Env = append(env, fmt.Sprintf("GOOS=%s", GOOS), fmt.Sprintf("GOARCH=%s", GOARCH))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("Error building for %s/%s: %v\n", GOOS, GOARCH, err)
			os.Exit(1)
		}
	}
}
