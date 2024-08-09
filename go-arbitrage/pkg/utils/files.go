package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ReadFile(logfile string) string {
	chdir()
	file, err := os.Open(logfile)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return fmt.Sprintf("%s", content)
}

func UpdateFile(logfile, content string) {
	chdir()
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777) //读写模式打开，复盖写入
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	f.Write([]byte(content))
}

func AppendFile(logfile, content string) {
	chdir()
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777) //读写模式打开，写入追加
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	f.Write([]byte(content + "\n"))
}

func chdir() {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	execDir := filepath.Dir(execPath)
	os.Chdir(execDir)
	//curpath, _ := os.Getwd()
	//index := strings.Index(curpath, "/goarbitrage") + len("/goarbitrage")
	//setpath := curpath[:index]
	//os.Chdir(setpath)
}
