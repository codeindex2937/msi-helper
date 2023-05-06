package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	model "github.com/codeindex2937/msi-helper/model"
)

var LangSpec = map[string][2]string{
	// lang_id, LUID Decimal, code_page
	"enu": {"1033", "1252"},
	"chs": {"2052", "936"},
	"cht": {"1028", "950"},
	"csy": {"1029", "1250"},
	"dan": {"1030", "1252"},
	"ger": {"1031", "1252"},
	"spn": {"3082", "1252"},
	"fre": {"1036", "1252"},
	"hun": {"1038", "1250"},
	"ita": {"1040", "1252"},
	"jpn": {"1041", "932"},
	"krn": {"1042", "949"},
	"nld": {"1043", "1252"},
	"nor": {"1044", "1252"},
	"plk": {"1045", "1250"},
	"ptb": {"1046", "1252"},
	"ptg": {"2070", "1252"},
	"rus": {"1049", "1251"},
	"sve": {"1053", "1252"},
	"tha": {"1054", "874"},
	"trk": {"1055", "1254"},
}

func (h Helper) SetLanguage(id string) error {
	langCode, codePage := LangSpec[id][0], LangSpec[id][1]

	if err := h.DB.Update(model.Property{
		Property: "ProductLanguage",
		Value:    langCode,
	}, nil, nil); err != nil {
		return err
	}

	summary, err := h.DB.OpenSummaryInformation(20)
	if err != nil {
		return err
	}
	defer summary.Close()

	packageLang, err := summary.GetString(model.PID_TEMPLATE)
	if err != nil {
		return err
	}
	parts := strings.SplitN(packageLang, ";", 2)
	if len(parts) < 2 {
		parts[0] = ""
	}
	if err := summary.SetProperty(model.PID_TEMPLATE, parts[0]+";"+langCode); err != nil {
		return err
	}
	if err := summary.Persist(); err != nil {
		return err
	}

	f, err := os.CreateTemp("", "*")
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("\n\n%v\t_ForceCodepage\n", codePage)); err != nil {
		return err
	}
	f.Close()

	if err := h.DB.Import(filepath.Dir(f.Name()), filepath.Base(f.Name())); err != nil {
		return err
	}

	return nil
}
