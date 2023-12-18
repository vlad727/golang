// 
package main

import (
	"encoding/base64"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// struct for kubeconfig file
type KubeStruct struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster,omitempty"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
	Kind           string `yaml:"kind"`
	Preferences    struct {
	} `yaml:"preferences"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

// func for decode certs and key form base64
func base64d(x string) string {

	decoded, err := base64.StdEncoding.DecodeString(x)
	if err != nil {
		log.Println("Error to decode base64 data")
	}
	return string(decoded)
}

func main() {
	// create an object for struct
	var c KubeStruct

	// read the file kubeconfig
	fileKube, err := os.ReadFile("kubeconfig")
	if err != nil {
		log.Println("File not found")
	}

	// unmarshal yaml file
	yaml.Unmarshal([]byte(fileKube), &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// vars for ca.crt tls.crt and tls.key
	CACRT := c.Clusters[0].Cluster.CertificateAuthorityData
	TLSCRT := c.Users[0].User.ClientCertificateData
	TLSKEY := c.Users[0].User.ClientKeyData

	// create map 
	m := map[string]string{
		"ca.crt":  CACRT,
		"tlc.crt": TLSCRT,
		"tls.key": TLSKEY,
	}
  // loop for map 
	for k, v := range m {

		decodedDate := base64d(v)
    // create separate files ca.crt, tlc.crt and tls.key with decoded data
		err := os.WriteFile(k, []byte(decodedDate), 0644)
		if err != nil {

			log.Println("Can't create file")
		}

	}

}
