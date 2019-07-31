package table

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/segmentio/terraform-docs/internal/pkg/doc"
	"github.com/segmentio/terraform-docs/internal/pkg/print/markdown"
	"github.com/segmentio/terraform-docs/internal/pkg/settings"
)

// Print prints a document as Markdown tables.
func Print(document *doc.Doc, printSettings settings.Settings) (string, error) {
	var buffer bytes.Buffer

	if printSettings.Has(settings.WithProviders) {
		printProviders(&buffer, document.Providers)
	}

	printVariables(&buffer, document.Variables, printSettings)
	printOutputs(&buffer, document.Outputs, printSettings)

	return markdown.Sanitize(buffer.String()), nil
}

func printProviders(buffer *bytes.Buffer, providers []doc.Provider) {
	buffer.WriteString("## Providers\n\n")

	if len(providers) == 0 {
		buffer.WriteString("None\n\n")
	} else {

		buffer.WriteString("| Name | Alias | Version |\n")
		buffer.WriteString("|------|-------|---------|\n")

		for _, provider := range providers {
			buffer.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				strings.ReplaceAll(provider.Name, "_", "\\_"),
				strings.ReplaceAll(provider.Alias, "_", "\\_"),
				provider.Version))
		}

		buffer.WriteString("\n")
	}
}

func getVariableDefaultValue(variable *doc.Variable) string {
	var result = "n/a"

	if variable.HasDefault() {
		result = markdown.PrintCode(variable.Default, "json")
	}

	return result
}

func printVariables(buffer *bytes.Buffer, variables []doc.Variable, printSettings settings.Settings) {
	buffer.WriteString("## Variables\n\n")

	if len(variables) == 0 {
		buffer.WriteString("None\n\n")
	} else {



		if printSettings.Has(settings.WithRequired) {
			buffer.WriteString("| Name | Description | Type | Default | Required |\n")
			buffer.WriteString("|------|-------------|------|---------|----------|\n")

		} else {
			buffer.WriteString("| Name | Description | Type | Default |\n")
			buffer.WriteString("|------|-------------|------|---------|\n")
		}

		for _, variable := range variables {
			buffer.WriteString("|")
			buffer.WriteString(fmt.Sprintf("%s | ", variable.Name))
			buffer.WriteString(fmt.Sprintf("%s | ", markdown.ConvertMultiLineText(variable.Description)))
			buffer.WriteString(fmt.Sprintf("%s | ", markdown.PrintCode(variable.Type, "hcl")))
			buffer.WriteString(fmt.Sprintf("%s |", getVariableDefaultValue(&variable)))
			if printSettings.Has(settings.WithRequired) {
				buffer.WriteString(fmt.Sprintf(" %s|", printIsVariableRequired(&variable)))
			}
			buffer.WriteString("\n")
		}
	}
}

func printIsVariableRequired(variable *doc.Variable) string {
	if !variable.HasDefault() {
		return "yes"
	}

	return "no"
}

func printOutputs(buffer *bytes.Buffer, outputs []doc.Output, settings settings.Settings) {
	buffer.WriteString("## Outputs\n\n")

	if len(outputs) == 0 {
		buffer.WriteString("None\n\n")
	} else {
		buffer.WriteString("| Name | Description |\n")
		buffer.WriteString("|------|-------------|\n")

		for _, output := range outputs {
			buffer.WriteString(
				fmt.Sprintf("| %s | %s |\n",
					strings.Replace(output.Name, "_", "\\_", -1),
					markdown.ConvertMultiLineText(output.Description)))
		}
		buffer.WriteString("\n")
	}
}
