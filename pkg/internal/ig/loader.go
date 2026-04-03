package ig

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5"
)

// Loader handles loading FHIR resources from the IG
type Loader struct {
	CodeSystems map[string]*r5.CodeSystem
	ValueSets   map[string]*r5.ValueSet
}

func NewLoader() *Loader {
	return &Loader{
		CodeSystems: make(map[string]*r5.CodeSystem),
		ValueSets:   make(map[string]*r5.ValueSet),
	}
}

// LoadFromIG loads terminologies from the given IG path
func (l *Loader) LoadFromIG(igPath string) error {
	fshPath := filepath.Join(igPath, "input", "fsh")

	// Load CodeSystems
	csDir := filepath.Join(fshPath, "codeSystems")
	err := filepath.Walk(csDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".fsh") {
			return nil
		}
		return l.parseFSHFile(path)
	})
	if err != nil {
		return fmt.Errorf("walk codeSystems: %w", err)
	}

	// Load ValueSets
	vsDir := filepath.Join(fshPath, "valueSets")
	err = filepath.Walk(vsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".fsh") {
			return nil
		}
		return l.parseFSHFile(path)
	})
	if err != nil {
		return fmt.Errorf("walk valueSets: %w", err)
	}

	return nil
}

func (l *Loader) parseFSHFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentCS *r5.CodeSystem
	var currentVS *r5.ValueSet

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			// Handle code entries starting with * #
			if strings.HasPrefix(line, "* #") && currentCS != nil {
				codeLine := strings.TrimPrefix(line, "* #")
				codeParts := strings.SplitN(codeLine, "\"", 2)
				code := strings.TrimSpace(codeParts[0])
				display := ""
				if len(codeParts) > 1 {
					display = strings.Trim(codeParts[1], "\" ")
				}
				currentCS.Concept = append(currentCS.Concept, r5.CodeSystemConcept{
					Code:    code,
					Display: &display,
				})
			}
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "CodeSystem":
			currentCS = &r5.CodeSystem{
				Name:    &value,
				Status:  "active",
				Content: "complete",
			}
			currentVS = nil
		case "ValueSet":
			currentVS = &r5.ValueSet{
				Name:   &value,
				Status: "active",
			}
			currentCS = nil
		case "Id":
			if currentCS != nil {
				currentCS.ID = &value
			} else if currentVS != nil {
				currentVS.ID = &value
			}
		case "Title":
			val := strings.Trim(value, "\"")
			if currentCS != nil {
				currentCS.Title = &val
			} else if currentVS != nil {
				currentVS.Title = &val
			}
		case "Description":
			val := strings.Trim(value, "\"")
			if currentCS != nil {
				currentCS.Description = &val
			} else if currentVS != nil {
				currentVS.Description = &val
			}
		case "* ^url":
			val := strings.Trim(value, "=\" ")
			if currentCS != nil {
				currentCS.URL = &val
				l.CodeSystems[val] = currentCS
			} else if currentVS != nil {
				currentVS.URL = &val
				l.ValueSets[val] = currentVS
			}
		}
	}

	return scanner.Err()
}

func pointer[T any](v T) *T {
	return &v
}
