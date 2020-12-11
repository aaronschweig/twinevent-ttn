package ditto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
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

func (ds *DittoService) CreateDT(device *ttnsdk.Device) error {
	url := fmt.Sprintf("%s/api/2/things/%s:%s", ds.config.Ditto.Host, ds.config.Ditto.Namespace, device.DevID)
	body := fmt.Sprintf(`{"policyId":"%s:%s"}`, ds.config.Ditto.Namespace, "twin-policy")
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(body))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("ditto", "ditto") // TODO:

	_, err = http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	return nil
}

func (ds *DittoService) CreateTTNConnection() error {

	funcsMap := template.FuncMap{
		"enquote": func(arg string) string {
			return fmt.Sprintf("'%s'", arg)
		},
	}

	// CREATE CONNECTION

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
		return fmt.Errorf("Connection could not be created. Got Status Code %d", res.StatusCode)
	}

	var r map[string]interface{}

	json.NewDecoder(res.Body).Decode(&r)

	// Miserable API-Design from Ditto
	r = r["?"].(map[string]interface{})
	r = r["?"].(map[string]interface{})

	status := int(r["status"].(float64))

	if !(status == http.StatusConflict || status == http.StatusCreated) {
		return fmt.Errorf("Connection could not be created. Got Status Code %d.\n Description: %s", status, r["description"])
	}

	// CREATE POLICY
	policyTmpl := template.Must(template.New("policy.template.json").Funcs(funcsMap).ParseFiles("./ditto/policy.template.json"))

	err = policyTmpl.Execute(&buffer, ds.config)

	if err != nil {
		return err
	}

	req, err = http.NewRequest(http.MethodPut, ds.config.Ditto.Host+"/api/2/policies/"+ds.config.Ditto.Namespace+":twin-policy", &buffer)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("ditto", "ditto") // TODO:

	res, err = http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if !(res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusNoContent) {
		return fmt.Errorf("Policy could not be created. Got Status Code %d", res.StatusCode)
	}

	return nil
}
