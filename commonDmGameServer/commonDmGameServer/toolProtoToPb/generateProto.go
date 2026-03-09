package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	targetFolder := filepath.Join(currentDir, "pb")       // 目标文件夹
	srcDir := filepath.Join(currentDir, "protocol")       // 源目录
	destDir1 := filepath.Join(currentDir, "toolpb\\temp") // 目录1
	sourceFolder := destDir1

	// 删除目标文件夹下的所有文件
	err = removeAllFiles(targetFolder)
	if err != nil {
		fmt.Println("Error removing target folder:", err)
		return
	}
	// 复制proto文件到目录1
	err = copyProtoFiles(srcDir, destDir1)
	if err != nil {
		fmt.Println("Error copying proto files:", err)
		return
	}
	p := filepath.Join(currentDir, "toolpb\\bin\\protoc.exe")
	// 递归遍历源文件夹下的所有文件
	err = filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 如果是文件并且文件名以 ".proto" 结尾
		if !info.IsDir() && filepath.Ext(info.Name()) == ".proto" {
			// 构建目标文件路径
			relativePath, err := filepath.Rel(sourceFolder, path)
			if err != nil {
				return err
			}
			targetPath := filepath.Join(targetFolder, fmt.Sprintf("%s.pb.go", relativePath[:len(relativePath)-len(".proto")]))

			// 调用 protoc 编译 protobuf 文件
			cmd := exec.Command(
				p,
				"--go_out="+targetFolder,
				"--proto_path="+sourceFolder, // 指定 --proto_path 参数
				path,
			)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				// return fmt.Errorf("failed to compile %s: %v", path, err)
				fmt.Println("failed to compile %s: %v", path, err)
			}
			// 修改生成的 .pb.go 文件的 package
			err = modifyPackage(targetPath)
			if err != nil {
				fmt.Println("failed to compile %s: %v", path, err)
			}
			fmt.Printf("Compiled %s to %s\n", path, targetPath)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	// 删除目录1下的所有proto文件
	err = deleteProtoFiles(destDir1)
	if err != nil {
		fmt.Println("Error deleting proto files:", err)
		return
	}

	fmt.Println("Process completed successfully.")
}
func copyProtoFiles(srcDir, destDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".proto") {
			destFile := filepath.Join(destDir, info.Name())
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(destFile, data, info.Mode())
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func compileProtoFiles(srcDir, destDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".proto") {
			destFile := filepath.Join(destDir, strings.TrimSuffix(info.Name(), ".proto")+".pb.go")
			cmd := exec.Command("protoc", "--go_out="+destDir, path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
			fmt.Printf("Compiled: %s\n", destFile)
		}
		return nil
	})
	return err
}

func deleteProtoFiles(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".proto") {
			err := os.Remove(filepath.Join(dir, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// removeAllFiles 删除目录下的所有文件
func removeAllFiles(dirPath string) error {
	files, err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}

	return nil
}

// modifyPackage 读取文件内容，将 package 替换为指定值，然后写回文件
func modifyPackage(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 替换 package 行
	newContent := strings.Replace(string(content), "package __", "package pb", 1)

	// 写回文件
	err = ioutil.WriteFile(filePath, []byte(newContent), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
