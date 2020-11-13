package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/snichols/election2020/pkg/states"
)

func download(url string, out string) error {
	var response *http.Response
	var err error

	if response, err = http.Get(url); err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}

	var data []byte
	if data, err = ioutil.ReadAll(response.Body); err != nil {
		return err
	}

	var j map[string]interface{}
	if err = json.Unmarshal(data, &j); err != nil {
		return err
	}

	if data, err = json.MarshalIndent(j, "", "    "); err != nil {
		return err
	}

	if err = ioutil.WriteFile(out, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func main() {
	for _, n := range states.Name {
		url := fmt.Sprintf("https://static01.nyt.com/elections-assets/2020/data/api/2020-11-03/race-page/%s/president.json", n)
		out := fmt.Sprintf("%s.json", n)
		fmt.Println("downloading:", out)
		if err := download(url, out); err != nil {
			panic(err)
		}
	}
}
