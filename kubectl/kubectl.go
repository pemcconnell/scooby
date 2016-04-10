package kubectl

/*
@todo set fire to this.

This is VERY far from ideal - there's a kubernetes API
which I should try and avail of instead of shelling out
to the cli. It's been written this way as I want to get
a working concept together asap.
*/

import (
	log "github.com/Sirupsen/logrus"

	"encoding/json"
	"os"
)

func IsEnvReady() {
	log.Debug("checking that your environment is ready [kubectl]")
}

func Deploy(subdomain, tag, port, dir string) {
	log.Debug("running deploy from kubectl")

	kubeConfigs(subdomain, tag, port, dir)
}

func kubeConfigs(subdomain, tag, port, dir string) {
	// create .kube directory
	os.Mkdir(dir+".kube", 0755)

	// create service yaml
	kubeService(subdomain, tag, port, dir)

	// create deployment yaml
	kubeDeployment(subdomain, tag, port, dir)

	// create .kube/config file
	kubeConfig(subdomain, tag, dir)

}

func kubeService(subdomain, tag, port, dir string) {
	filepath := dir + ".kube/service.yml"
	content := []byte(`apiVersion: v1
kind: Service
metadata:
  name: ` + subdomain + `
spec:
  type: LoadBalancer
  ports:
  - port: ` + port + `
  selector:
    app: ` + subdomain + `
    tier: ` + subdomain)
	writeFile(filepath, content)
}

func kubeDeployment(subdomain, tag, port, dir string) {
	filepath := dir + ".kube/deployment.yml"
	content := []byte(`apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: ` + subdomain + `
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: ` + subdomain + `
        tier: ` + subdomain + `
    spec:
      containers:
      - name: ` + subdomain + `
        image: ` + tag + `
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
        env:
        - name: GET_HOSTS_FROM
          value: dns
        ports:
        - containerPort: ` + port)
	writeFile(filepath, content)
}

func kubeConfig(subdomain, tag, dir string) {
	kubeconfig := dir + ".kube/config"
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		content := &KubeConfig{
			container: KubeConfigContainer{
				subdomain: subdomain,
				tag:       tag,
			},
		}
		j, jerr := json.MarshalIndent(content, "", "  ")
		if jerr != nil {
			log.Fatal(jerr)
		}
		writeFile(kubeconfig, j)
	}
}

func writeFile(path string, content []byte) {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(content); err != nil {
		log.Fatal(err)
	}
}

type KubeConfig struct {
	container KubeConfigContainer
}

type KubeConfigContainer struct {
	tag       string
	subdomain string
}
