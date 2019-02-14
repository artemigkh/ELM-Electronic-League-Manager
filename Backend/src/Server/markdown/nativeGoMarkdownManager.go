package markdown

import (
	"Server/config"
	"github.com/Pallinder/go-randomdata"
	"io/ioutil"
	"os"
	"path/filepath"
)

type GoMdManager struct {
	OutPath string
}

func CreateGoMarkdownManager(conf config.Config) *GoMdManager {
	path := conf.GetMarkdownDir()
	println(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	return &GoMdManager{
		OutPath: path,
	}
}

func (m *GoMdManager) StoreMarkdown(leagueId int, markdown, oldFile string) (string, error) {
	fileName := randomdata.RandStringRunes(10) + ".md"
	f, err := os.Create(filepath.Join(m.OutPath, fileName))
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString(markdown)
	if err != nil {
		return "", err
	}

	if oldFile != "" {
		os.Remove(filepath.Join(m.OutPath, oldFile))
	}

	return fileName, nil
}

func (m *GoMdManager) GetMarkdown(fileName string) (string, error) {
	md, err := ioutil.ReadFile(filepath.Join(m.OutPath, fileName))
	if err != nil {
		return "", err
	}

	return string(md), nil
}
