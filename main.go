package main

import (
	"github.com/appscode/go/runtime"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sort"
)

/*
                <type-mappings>intstr.IntOrString=IntOrString,resource.Quantity=Quantity</type-mappings>
                <import-mappings>IntOrString=io.kubernetes.client.custom.IntOrString,Quantity=io.kubernetes.client.custom.Quantity</
*/
func main() {
	typeMappings := map[string]string{
		"intstr.IntOrString":"IntOrString",
		"resource.Quantity":"Quantity",
	}
	importMappings := map[string]string{
		"IntOrString":"io.kubernetes.client.custom.IntOrString",
		"Quantity":"io.kubernetes.client.custom.Quantity",
	}

	filename := runtime.GOPath() + "/src/github.com/kubernetes-client/java/kubernetes/src/main/java/io/kubernetes/client/models"
	fmt.Println(filename)
	files, err := ioutil.ReadDir(filename)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		hasV1Prefix := strings.HasPrefix(f.Name(), "V1") && !strings.HasPrefix(f.Name(), "V1alpha") && !strings.HasPrefix(f.Name(), "V1beta")
		if !hasV1Prefix {
			continue
		}

		name := f.Name()
		name = name[:len(name)-len(".java")]

		// io.kubernetes.client.models.V1Service.java

		swaggerDefName := "v1." + name[2:]
		className := name
		fullName := "io.kubernetes.client.models." + name

		if _, found := typeMappings[swaggerDefName]; !found {
			typeMappings[swaggerDefName] = className
		}
		if _, found := importMappings[className]; !found {
			importMappings[className] = fullName
		}
	}

	tm := make([]string, 0, len(typeMappings))
	for k, v := range typeMappings {
		tm = append(tm, k+"="+v)
	}
	sort.Strings(tm)

	im := make([]string, 0, len(importMappings))
	for k, v := range importMappings {
		im = append(im, k+"="+v)
	}
	sort.Strings(im)

	fmt.Println()
	fmt.Printf("<type-mappings>%s</type-mappings>\n", strings.Join(tm,","))
	fmt.Println()
	fmt.Printf("<import-mappings>%s</import-mappings>\n", strings.Join(im,","))
	fmt.Println()
}
