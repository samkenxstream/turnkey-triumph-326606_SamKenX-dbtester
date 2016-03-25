package analyze

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	ReadmeCommand = &cobra.Command{
		Use:   "readme",
		Short: "Generates README markdown.",
		RunE:  ReadmeCommandFunc,
	}

	readmeDir         string
	readmePrefacePath string
)

func init() {
	ReadmeCommand.PersistentFlags().StringVarP(&readmeDir, "readme-dir", "d", "", "Directory path to generate README.")
	ReadmeCommand.PersistentFlags().StringVarP(&readmePrefacePath, "readme-preface", "t", "", "README template file to preface.")
}

func ReadmeCommandFunc(cmd *cobra.Command, args []string) error {
	bts, err := ioutil.ReadFile(readmePrefacePath)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.WriteString(string(bts))
	buf.WriteString("\n\n\n")

	buf.WriteString("<br><br><hr>\n##### Results")
	buf.WriteString("\n\n")

	paths, err := walkExt(readmeDir, ".svg")
	if err != nil {
		return err
	}
	for _, path := range paths {
		buf.WriteString(fmt.Sprintf("![%s](./%s)\n\n", filepath.Base(path), filepath.Base(path)))
	}

	return toFile(buf.String(), filepath.Join(readmeDir, "README.md"))
}

// walkExt returns all FileInfos with specific extension.
func walkExt(targetDir, ext string) ([]string, error) {
	rmap := make(map[string]struct{})
	visit := func(path string, f os.FileInfo, err error) error {
		if f != nil {
			if !f.IsDir() {
				if filepath.Ext(path) == ext {
					if !filepath.HasPrefix(path, ".") && !strings.Contains(path, "/.") {
						wd, err := os.Getwd()
						if err != nil {
							return err
						}
						thepath := filepath.Join(wd, strings.Replace(path, wd, "", -1))
						rmap[thepath] = struct{}{}
					}
				}
			}
		}
		return nil
	}
	err := filepath.Walk(targetDir, visit)
	if err != nil {
		return nil, err
	}
	rs := []string{}
	for k := range rmap {
		rs = append(rs, k)
	}
	sort.Strings(rs)
	return rs, nil
}

func toFile(txt, fpath string) error {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
	}
	defer f.Close()
	if _, err := f.WriteString(txt); err != nil {
		return err
	}
	return nil
}
