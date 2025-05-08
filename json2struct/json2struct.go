package json2struct

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Config holds configuration options for struct generation
type Config struct {
	StructPrefix string // Prefix for all generated struct names
	IndentSize   int    // Number of spaces for indentation
}

// Generator handles the JSON to Go struct conversion process
type Generator struct {
	structMap map[string]string // Stores generated struct definitions
	config    Config            // Configuration for generation
}

// NewGenerator creates a new Generator instance with the given configuration
func NewGenerator(config Config) *Generator {
	return &Generator{
		structMap: make(map[string]string),
		config:    config,
	}
}

// GenerateFromJSON generates Go structs from JSON data
func (g *Generator) GenerateFromJSON(data []byte) error {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("JSON parse error: %v", err)
	}

	return g.generateStruct(jsonData, "Root")
}

// WriteToFile writes all generated structs to a file with random name
func (g *Generator) WriteToFile() (string, error) {
	filename := "data/" + generateRandomFilename("struct")

	// Get sorted keys
	keys := make([]string, 0, len(g.structMap))
	for k := range g.structMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	output := []byte{}
	for _, key := range keys {
		output = append(output, []byte(g.structMap[key]+"\n")...)
	}

	if err := os.WriteFile(filename, output, 0644); err != nil {
		return "", fmt.Errorf("file write error: %v", err)
	}

	return filename, nil
}

// generateUniqueName creates a unique struct name with prefix
func (g *Generator) generateUniqueName(baseName string) string {
	name := g.config.StructPrefix + baseName
	counter := 1

	for {
		if _, exists := g.structMap[name]; !exists {
			return name
		}
		name = fmt.Sprintf("%s%s%d", g.config.StructPrefix, baseName, counter)
		counter++
	}
}

// generateField creates a struct field definition with proper indentation and JSON tags
func (g *Generator) generateField(fieldName, fieldType, jsonTag string) string {
	indent := strings.Repeat(" ", g.config.IndentSize)
	return fmt.Sprintf("%s%s %s `json:\"%s\"`", indent, fieldName, fieldType, jsonTag)
}

// generateStruct recursively generates Go struct definitions from JSON data
func (g *Generator) generateStruct(data interface{}, structName string) error {
	// Generate unique struct name with prefix
	structName = g.generateUniqueName(structName)

	var fields []string

	// Process map type data
	switch v := data.(type) {
	case map[string]interface{}:
		// Iterate through each field in the map
		for key, value := range v {
			fieldName := toTitle(key)

			// Handle different types of values
			switch value.(type) {
			case map[string]interface{}:
				// Nested struct case
				nestedStructName := fieldName
				if err := g.generateStruct(value, nestedStructName); err != nil {
					return fmt.Errorf("error generating nested struct %s: %v", nestedStructName, err)
				}
				fields = append(fields, g.generateField(fieldName, nestedStructName, key))

			case []interface{}:
				// Array case
				if len(value.([]interface{})) > 0 {
					switch value.([]interface{})[0].(type) {
					case map[string]interface{}:
						// Array of structs case
						nestedStructName := fieldName + "Element"
						if err := g.generateStruct(value.([]interface{})[0], nestedStructName); err != nil {
							return fmt.Errorf("error generating array struct %s: %v", nestedStructName, err)
						}
						fields = append(fields, g.generateField(fieldName, "[]"+nestedStructName, key))
					default:
						// Array of primitive types case
						fields = append(fields, g.generateField(fieldName,
							fmt.Sprintf("[]%T", value.([]interface{})[0]), key))
					}
				} else {
					// Empty array case
					fields = append(fields, g.generateField(fieldName, "[]interface{}", key))
				}

			default:
				// Primitive type case
				fields = append(fields, g.generateField(fieldName, fmt.Sprintf("%T", value), key))
			}
		}

	default:
		return fmt.Errorf("unsupported data type: %T", data)
	}

	// Store the generated struct definition
	g.structMap[structName] = fmt.Sprintf("type %s struct {\n%s\n}\n",
		structName, strings.Join(fields, "\n"))

	return nil
}
