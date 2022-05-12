package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/alexey-mavrin/graduate-2/internal/common"
)

//go:embed action.tmpl
var actionTemplate string

type actionParams struct {
	RecordType common.RecordType
}

func generateCode(tmpl string, actionType common.RecordType) ([]byte, error) {
	buf := new(bytes.Buffer)
	t, err := template.New("code").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	err = t.Execute(buf, actionParams{RecordType: actionType})
	if err != nil {
		return nil, err
	}
	p, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return p, nil
}

func usage() string {
	return fmt.Sprintf("usage: %s ACTION_TYPE\n",
		filepath.Base(os.Args[0]),
	)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(usage())
	}

	actionType := common.RecordType(os.Args[1])

	p, err := generateCode(actionTemplate, actionType)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(string(actionType)+"_action.go", p, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
