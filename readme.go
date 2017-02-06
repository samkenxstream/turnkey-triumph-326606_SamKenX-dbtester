package dbtester

import (
	"bytes"
	"fmt"
	"path/filepath"
)

// WriteREADME writes README.
func (cfg *Config) WriteREADME(summary string) error {
	plog.Printf("writing README at %q", cfg.README.OutputPath)

	buf := new(bytes.Buffer)
	buf.WriteString("\n\n")
	buf.WriteString(fmt.Sprintf("<br><br><hr>\n##### %s", cfg.TestTitle))
	buf.WriteString("\n\n")
	buf.WriteString(cfg.TestDescription)
	buf.WriteString("\n\n```\n")
	buf.WriteString(summary)
	buf.WriteString("```\n\n\n")

	for _, img := range cfg.Images {
		switch img.Type {
		case "local":
			imgPath := "./" + filepath.Base(img.Path)
			buf.WriteString(fmt.Sprintf("![%s](%s)\n\n", img.Title, imgPath))
		case "remote":
			buf.WriteString(fmt.Sprintf(`<img src="%s" alt="%s">`, img.Path, img.Title))
			buf.WriteString("\n\n")
		default:
			return fmt.Errorf("%s is not supported", img.Type)
		}
		buf.WriteString("\n\n")
	}

	return toFile(buf.String(), cfg.README.OutputPath)
}
