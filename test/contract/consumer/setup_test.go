package consumer

import (
	"bytes"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"io/ioutil"
	"text/template"
)

func _loadContractTestConfig(pactPort int) {
	templateName := "config_contracttests"

	configYamlTemplateBytes, _ := ioutil.ReadFile("../../resources/config-contracttests.yaml")
	tmpl, _ := template.New(templateName).Parse(string(configYamlTemplateBytes))

	var buf bytes.Buffer
	parameters := map[string]interface{} {
		"pactPort": pactPort,
	}

	_ = tmpl.ExecuteTemplate(&buf, templateName, parameters)

	_ = config.ParseAndOverwriteConfig(buf.Bytes())
}
