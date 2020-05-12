package printing

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/alexbrainman/printer"
)

var (
	networkNames map[string]string
)

func init() {
	networkNames = map[string]string{}

	output, _ := exec.Command("powershell", "Get-Printer", "|", "select", "-exp", "name").Output()
	// log.Println(string(output))
	names := strings.Fields(strings.ToLower(string(output)))
	for _, v := range names {
		if len(v) > 2 && v[0:2] == "\\\\" {

			parts := strings.Split(v, "\\")
			networkNames[parts[len(parts)-1]] = v
			networkNames[strings.ToUpper(parts[len(parts)-1])] = v
		}

	}
}

func printOneDocument(printerName, documentName string, output []byte) error {
	p, err := printer.Open(printerName)
	if err != nil {
		return err
	}
	defer p.Close()

	err = p.StartRawDocument(documentName)
	if err != nil {
		return err
	}
	defer p.EndDocument()

	err = p.StartPage()
	if err != nil {
		return err
	}

	p.Write(output)

	return p.EndPage()
}

func PrintDocument(printerName string, path string, avoidNetwork bool) error {
	if !avoidNetwork {
		n := networkNames[printerName]
		if n != "" {
			printerName = n
		}
	}
	output, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = printOneDocument(printerName, path, output)
	if err != nil {
		return err
	}

	return nil
}
