package main

import (
	"archive/zip"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	fmt.Printf("程序启动\n")

	args := os.Args
	num := len(args)

	var books []string

	if num > 1 {
		for _, arg := range args[1:] {

			info, err := os.Stat(arg)
			if err != nil {
				log.Fatal(err)
			}
			if !info.IsDir() {
				continue
			}

			results, err := search(arg)
			if err != nil {
				log.Fatal(err)
			}

			books = append(books, results...)
		}
	} else {
		fmt.Printf("无输入，程序结束\n")
	}

	bookNum := len(books)
	if bookNum > 0 {
		fmt.Printf("已找到%d个zip文件，尝试寻找同目录下的opf文件并自动重命名\n", bookNum)
		for _, book := range books {
			err := process(book)
			if err != nil {
				log.Fatal(err)
			}
		}

	} else {
		fmt.Printf("未能找到zip文件，程序结束\n")
	}
	fmt.Scanln()
}

func search(arg string) ([]string, error) {
	fmt.Printf("开始搜索目录「%s」下的文件\n", arg)
	var books []string
	err := filepath.WalkDir(arg, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == ".zip" {
				books = append(books, path)
			}
		}

		return err
	})

	if err != nil {
		return books, err
	}

	return books, err
}

func process(book string) error {
	fmt.Printf("开始处理书籍「%s」\n", book)

	// 查找opf文件
	opf := filepath.Join(filepath.Dir(book), "metadata.opf")
	_, err := os.Stat(opf)

	// 检查是否存在，如果没有则打印警告并跳过之后的处理
	if os.IsNotExist(err) {
		log.Warn(err)
		return nil
	}

	b, err := os.ReadFile(opf)
	if err != nil {
		return err
	}

	json, err := gjson.LoadXml(b)
	if err != nil {
		return err
	}
	title := json.GetString("package.metadata.title", "Unknow")
	creator := json.GetString("package.metadata.creator.#text", "Unknow")

	fixedCreator := fixCreator(creator)
	fixedTitle, dir := fixTitle(title, fixedCreator)

	fmt.Printf("输出文件名：%s\n", fixedTitle)
	outputDirPath := filepath.Join("./output", dir)
	outputFilePath := filepath.Join(outputDirPath, fixedTitle)

	_, err = os.Stat(outputDirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(outputDirPath, 0755)
		if err != nil {
			return err
		}
	}

	// 打开zip
	readCloser, err := zip.OpenReader(book)
	if err != nil {
		log.Fatal(err)
	}
	defer readCloser.Close()

	// 创建zip
	file, err := os.Create(outputFilePath + ".zip")
	if err != nil {
		return err
	}
	writer := zip.NewWriter(file)

	for index, file := range readCloser.File {
		file.Name = fmt.Sprintf("%04d", index+1) + filepath.Ext(file.Name)
		writer.Copy(file)
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return err
}

func fixCreator(creator string) string {
	// 去除一切非法字符
	creator = strings.ReplaceAll(creator, "\\", "")
	creator = strings.ReplaceAll(creator, "/", "")
	creator = strings.ReplaceAll(creator, ":", "")
	creator = strings.ReplaceAll(creator, "*", "")
	creator = strings.ReplaceAll(creator, "?", "")
	creator = strings.ReplaceAll(creator, "\"", "")
	creator = strings.ReplaceAll(creator, "<", "")
	creator = strings.ReplaceAll(creator, ">", "")
	creator = strings.ReplaceAll(creator, "|", "")
	// 去除空格
	creator = strings.ReplaceAll(creator, " ", "")
	creator = strings.ReplaceAll(creator, "　", "")
	// 全角数字换半角
	creator = strings.ReplaceAll(creator, "１", "1")
	creator = strings.ReplaceAll(creator, "２", "2")
	creator = strings.ReplaceAll(creator, "３", "3")
	creator = strings.ReplaceAll(creator, "４", "4")
	creator = strings.ReplaceAll(creator, "５", "5")
	creator = strings.ReplaceAll(creator, "６", "6")
	creator = strings.ReplaceAll(creator, "７", "7")
	creator = strings.ReplaceAll(creator, "８", "8")
	creator = strings.ReplaceAll(creator, "９", "9")
	creator = strings.ReplaceAll(creator, "０", "0")
	return creator
}

func fixTitle(title string, creator string) (string, string) {
	// 去除一切非法字符
	title = strings.ReplaceAll(title, "\\", "")
	title = strings.ReplaceAll(title, "/", "")
	title = strings.ReplaceAll(title, ":", "")
	title = strings.ReplaceAll(title, "*", "")
	title = strings.ReplaceAll(title, "?", "")
	title = strings.ReplaceAll(title, "\"", "")
	title = strings.ReplaceAll(title, "<", "")
	title = strings.ReplaceAll(title, ">", "")
	title = strings.ReplaceAll(title, "|", "")

	// 全角数字换半角
	title = strings.ReplaceAll(title, "１", "1")
	title = strings.ReplaceAll(title, "２", "2")
	title = strings.ReplaceAll(title, "３", "3")
	title = strings.ReplaceAll(title, "４", "4")
	title = strings.ReplaceAll(title, "５", "5")
	title = strings.ReplaceAll(title, "６", "6")
	title = strings.ReplaceAll(title, "７", "7")
	title = strings.ReplaceAll(title, "８", "8")
	title = strings.ReplaceAll(title, "９", "9")
	title = strings.ReplaceAll(title, "０", "0")

	// 尝试获取卷数
	vol := ""
	reg := regexp.MustCompile(`[0-9]+`)
	results := reg.FindAllString(title, -1)
	if len(results) > 0 {
		vol = results[len(results)-1]
	}

	// 去除空格
	title = strings.ReplaceAll(title, " ", "")
	title = strings.ReplaceAll(title, "　", "")

	// 去除一切括号和括号内文字
	reg = regexp.MustCompile(`[（,【,(].+?[),】,）]`)
	title = reg.ReplaceAllString(title, "")

	// 去除卷标
	reg = regexp.MustCompile(`[0-9]+巻`)
	title = reg.ReplaceAllString(title, "")
	title = strings.ReplaceAll(title, vol, "")

	// 重组
	var fixedTitle string
	dir := fmt.Sprintf("[%s] %s", creator, title)

	if vol != "" {
		fixedTitle = fmt.Sprintf("[%s] %s 第%s巻", creator, title, vol)
	} else {
		fixedTitle = dir
	}

	return fixedTitle, dir
}

func noTwoByteNumber(str string) string {

	return str
}

func noBracket(str string) string {
	reg := regexp.MustCompile(`[（,【,(].+?[),】,）]`)
	return reg.ReplaceAllString(str, "")
}
