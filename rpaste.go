package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var rpasteCmd = &cobra.Command{
	Use:   "rpaste [file name]",
	Short: "Paste file (or STDIN) to Rocky Pastebin",
	Long:  `An interface to the Rocky Pastebin.`,
	RunE:  rpaste,
}

func rpaste(cmd *cobra.Command, args []string) error {

	var fileContent string

	if len(args) > 0 {
		fileName := args[0]
		if _, err := os.Stat(fileName); err == nil {
			bytes, err := ioutil.ReadFile(fileName)
			if err != nil {
				fmt.Printf("Could not read file: %s (%s)\n", fileName, err)
				os.Exit(1)
			}
			fileContent = string(bytes)
		}
	} else {
		fmt.Printf("Reading from STDIN\n")
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("Could not read from STDIN: (%s)\n", err)
			os.Exit(1)
		}
		fileContent = string(bytes)

	}

	//fileContent := "Hello World\nThis is a test string\n\n"

	payload := map[string]interface{}{
		"expiry": "1day",
		"files": []map[string]string{
			{
				"name":    "spam",
				"lexer":   "text",
				"content": fileContent,
			},
		},
	}

	payload_json, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Could not marshal JSON\n")
		os.Exit(1)
	}

	resp, err := http.Post("https://rpa.st/api/v1/paste", "application/json", bytes.NewBuffer(payload_json))
	if err != nil {
		fmt.Printf("Could not create HTTP post communicator\n")
		os.Exit(1)
	}

	if resp.StatusCode == 200 {
		var response map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&response)

		fmt.Printf("Link:    %s\n", response["link"])
		fmt.Printf("Removal: %s\n", response["removal"])
	} else {
		fmt.Printf("ERROR: Got HTTP code: %d\n", resp.StatusCode)
	}

	return nil
}

func main() {
	rpasteCmd.Execute()
}
