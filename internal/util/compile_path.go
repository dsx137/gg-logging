package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var compileRoot string

func getCompileRoot() (string, error) {
	// 对于编译后的二进制文件，go.mod 的位置是固定的。
	// 也就是说，这个函数查找的是编译时 go.mod 所在的目录。
	//
	// 注意：这个函数在编译后，其 os.Executable() 仍然指向打包后的二进制文件。
	// 但如果 go.mod 是在编译目录的某个父级，这个回溯逻辑依然有效。
	//
	// 关键点在于：无论你在哪里运行编译好的程序，
	// 当它第一次调用这个函数时，它会从 os.Executable() 的位置开始向上找 go.mod。
	// 但是，如果 go.mod 并不在可执行文件的当前路径的父级，而是位于编译时的一个完全不同的路径，
	// 例如你在 /build/temp 编译，而 go.mod 在 /src/my-project，那这个方法就不适用。
	//
	// 更好的方法是使用 runtime.Caller(0).File 结合 GOPATH/GOROOT
	// 但对于标准模块项目，最可靠的是查找 go.mod
	// 让我们稍微修改一下 getProjectRoot，让它更明确地寻找编译时go.mod的位置
	//
	// 考虑在开发环境下，go.mod 总是位于项目的根目录。
	// 当你运行 `go run main.go` 时，`main.go` 所在目录就是当前工作目录。
	// 当你运行 `go build` 时，`main.go` 所在目录也是编译上下文的目录。
	// 所以，我们可以用一种更直接的方式来获取“编译时”的项目根目录，即从 `runtime.Caller(0)` 开始向上找 `go.mod`。

	// 从当前文件（common.go 或 main.go）的编译时路径开始向上查找 go.mod
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get caller info for project root determination")
	}

	currentDir := filepath.Dir(currentFile) // 当前文件所在的目录

	// 向上回溯查找 go.mod 文件
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir { // 已经到达文件系统根目录
			break
		}
		currentDir = parentDir
	}

	return "", nil
}

func GetRelativePath(fullPath string) string {
	if compileRoot == "" {
		return fullPath
	}
	relPath, err := filepath.Rel(compileRoot, fullPath)
	if err != nil {
		return fullPath
	}
	return relPath
}

func init() {
	var err error
	compileRoot, err = getCompileRoot()
	if err != nil {
		panic(err)
	}
}
