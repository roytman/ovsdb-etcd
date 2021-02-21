package codegenerator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	SchemaFile     	string
	OutputDir		string
	PkgName			string
)

func Run() {
	// data, err := ioutil.ReadFile("./json/_server.ovsschema")
	data, err := ioutil.ReadFile(SchemaFile)
	if err != nil {
		fmt.Printf("ERROR %v\n", err)
		return
	}
	schemaInst := map[string]interface{}{}
	err = json.Unmarshal(data, &schemaInst)
	if err != nil {
		fmt.Printf("ERROR %v\n", err)
		return
	}
	if len(PkgName) ==0 {
		PkgName = schemaInst["name"].(string)
	}
	tables, ok := schemaInst["tables"]
	if !ok {
		fmt.Println("No tables")
		return
	}

	tablesMap := tables.(map[string]interface{})
	fmt.Printf("package %s\n\n", PkgName)

	fmt.Printf( "import \"github.com/roytman/ovsdb-etcd/pkg/ovsdb\"\n\n")

	for tableName, table := range tablesMap {
		structColumns := []string{}
		tabMap := table.(map[string]interface{})
		columns, ok := tabMap["columns"]
		if !ok {
			fmt.Println("No columns")
			continue
		}
		columnsMap := columns.(map[string]interface{})
		for k, v := range columnsMap {
			fieldName := toUppercase(k)
			typeValue := v.(map[string]interface{})
			if vt, ok := typeValue["type"]; ok {
				switch vt.(type) {
				case string:
					l := fmt.Sprintf(" %s\t%s `json:\"%s\"`\n", fieldName, typeConvert(vt), k)
					structColumns = append(structColumns, l)
				case map[string]interface{}:
					k2Map := vt.(map[string]interface{})
					max, maxOK := k2Map["max"]
					t2, ok := k2Map["key"]
					if !ok {
						fmt.Printf("NO KEY\n")
						continue
					}
					// check if it is a map
					t2v, valueOK := k2Map["value"]
					var mapValue string
					if valueOK {
						mapValue = getValueType(t2v)
						/*l := fmt.Sprintf(" %s \tmap[%s]%s `json:\"%s\"`\n", fieldName, t2, t2v, k)
						structColumns = append(structColumns, l)
						continue*/
					}
					switch t2.(type) {
					case map[string]interface{}:
						k3Map := t2.(map[string]interface{})

						type3, ok := k3Map["type"]
						if !ok {
							fmt.Printf("NO TYPE 3\n")
							continue
						}
						var l string
						if valueOK {
							l = fmt.Sprintf(" %s \tmap[%s]%s `json:\"%s\"`\n", fieldName, typeConvert(type3), typeConvert(mapValue), k)
						} else if maxOK {
							maxStr := fmt.Sprintf("%v", max)
							if maxStr != "1" {
								l = fmt.Sprintf(" %s \t[]%s `json:\"%s\"`\n", fieldName, typeConvert(type3), k)
							} else {
								l = fmt.Sprintf(" %s \t%s `json:\"%s\"`\n", fieldName, typeConvert(type3), k)
							}
						} else {
							l = fmt.Sprintf(" %s \t%s `json:\"%s\"`\n", fieldName, typeConvert(type3), k)
						}
						structColumns = append(structColumns, l)

					case string:
						var l string
						if valueOK {
							l = fmt.Sprintf(" %s \tmap[%s]%s `json:\"%s\"`\n", fieldName, typeConvert(t2), typeConvert(mapValue), k)
						} else if maxOK {
							maxStr := fmt.Sprintf("%v", max)
							if maxStr != "1" {
								l = fmt.Sprintf(" %s \t[]%s `json:\"%s\"`\n", fieldName, typeConvert(t2), k)
							} else {
								l = fmt.Sprintf(" %s \t%s `json:\"%s\"`\n", fieldName, typeConvert(t2), k)
							}
						} else {
							l = fmt.Sprintf(" %s \t%s `json:\"%s\"`\n", fieldName, typeConvert(t2), k)
						}
						structColumns = append(structColumns, l)
					default:
						fmt.Printf("111 Default table %s column %s  %T %v\n", tableName, k, t2, t2)
					}

				}
			}


		}
		printStruct(tableName, structColumns)
	}
}
func printStruct(tableName string, columns []string) {
	fmt.Printf("type %s struct { \n", tableName)
	for _, line := range columns {
		fmt.Printf("\t%s", line)
	}
	fmt.Println("}\n\n")
}

func getValueType(value interface{}) string {
	switch value.(type) {
	case string:
		return value.(string)
	case map[string]interface{}:
		m:= value.(map[string]interface{})
		s, ok := m["type"]
		if ok {
			return s.(string)
		}
	}
	// TODO
	return ""
}

func typeConvert(typeName interface{}) string {
	s := fmt.Sprintf("%s", typeName)
	switch s {
	case "string":
		return s
	case "boolean":
		return "bool"
	case "integer":
		return "int64"
	case "real":
		return "float64"
	case "uuid":
		return "ovsdb.Uuid"
	}
	return ""
}

func toUppercase( name string) string {
	startLetter := name[:1]
	startLetter = strings.ToUpper(startLetter)
	return startLetter + name[1:]
}


