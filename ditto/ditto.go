package ditto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/aaronschweig/twinevent-ttn/config"
)

type DittoService struct {
	config *config.Config
}

func NewDittoService(cfg *config.Config) *DittoService {
	return &DittoService{
		cfg,
	}
}

func (ds *DittoService) CreateTTNConnection() error {

	funcsMap := template.FuncMap{
		"enquote": func(arg string) string {
			return fmt.Sprintf("'%s'", arg)
		},
	}

	configTmpl := template.Must(template.New("connection.template.json").Funcs(funcsMap).ParseFiles("./ditto/connection.template.json"))

	var buffer bytes.Buffer

	err := configTmpl.Execute(&buffer, ds.config)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, ds.config.Ditto.Host+"/devops/piggyback/connectivity", &buffer)

	if err != nil {
		return err
	}

	req.SetBasicAuth("devops", "foobar") // TODO:

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Ressource could not be created. Got Status Code %d", res.StatusCode)
	}

	var r map[string]interface{}

	json.NewDecoder(res.Body).Decode(&r)

	// Miserable API-Design from Ditto
	r = r["?"].(map[string]interface{})
	r = r["?"].(map[string]interface{})

	status := int(r["status"].(float64))

	if !(status == http.StatusConflict || status == http.StatusCreated) {
		return fmt.Errorf("Ressource could not be created. Got Status Code %d.\n Description: %s", status, r["description"])
	}

	return nil
}
