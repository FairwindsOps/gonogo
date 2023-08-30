package bundle

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2" // Import the YAML library
)

type StructB struct {
    // Define Struct B fields here
    Field1B string
    Field2B int
    // ...
}

type StructA struct {
    Field1 []StructB
    // Other Struct A fields, if any
}

func main() {
    var combinedStructA StructA

    // List of YAML file names
    yamlFiles := []string{"file1.yaml", "file2.yaml", "file3.yaml"}

    for _, fileName := range yamlFiles {
        // Read the content of the YAML file
        yamlData, err := ioutil.ReadFile(fileName)
        if err != nil {
            fmt.Printf("Error reading file %s: %v\n", fileName, err)
            continue
        }

        // Temporary Struct A for unmarshaling the Field1 data
        var tempStructA struct {
            Field1 []StructB `yaml:"Field1"`
        }

        // Unmarshal the YAML data into the temporary Struct A
        if err := yaml.Unmarshal(yamlData, &tempStructA); err != nil {
            fmt.Printf("Error unmarshaling file %s: %v\n", fileName, err)
            continue
        }

        // Append the Field1 data from the temporary Struct A to the combined Struct A
        combinedStructA.Field1 = append(combinedStructA.Field1, tempStructA.Field1...)
    }

    // Now, combinedStructA contains the merged data
    fmt.Printf("Merged Struct A: %+v\n", combinedStructA)
}
