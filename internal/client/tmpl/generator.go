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

//go:embed client.tmpl
var clientTemplate string

type clientTParams struct {
	RecordType common.RecordType
}

func generateClient(clientType common.RecordType) ([]byte, error) {
	buf := new(bytes.Buffer)
	t, err := template.New("client").Parse(clientTemplate)
	if err != nil {
		return nil, err
	}
	err = t.Execute(buf, clientTParams{RecordType: clientType})
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
	return fmt.Sprintf("usage: %s RECORD_TYPE\n",
		filepath.Base(os.Args[0]),
	)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(usage())
	}

	recordType := common.RecordType(os.Args[1])
	p, err := generateClient(recordType)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(string(recordType)+"_client.go", p, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
