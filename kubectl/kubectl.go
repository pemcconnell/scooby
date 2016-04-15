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
	"os/exec"
)

func IsEnvReady() {
	log.Debug("checking that your environment is ready [kubectl]")
}

func Deploy(subdomain, tag, port, dir string) {
	log.Debug("running deploy from kubectl")

	// create .kube directory
	log.Debug("creating .kube folder")
	os.Mkdir(dir+".kube", 0755)

	if _, err := os.Stat(dir + ".kube/service.yml"); os.IsNotExist(err) {
		// create service yaml
		kubeService(subdomain, tag, port, dir)
		kubeDeployService(dir)
	}

	kubeDeployment(subdomain, tag, port, dir)
	if _, err := os.Stat(dir + ".kube/deployment.yml"); os.IsNotExist(err) {
		// create deployment yaml
		kubeDeployDeployment(dir)
	} else {
		kubeUpdate(dir)
	}

	if _, err := os.Stat(dir + ".kube/config"); os.IsNotExist(err) {
		// create .kube/config file
		kubeConfig(subdomain, tag, dir)
	}

}

func kubeDeployService(dir string) {
	kubedir := dir + ".kube/"
	// service
	log.Debug("kubectl create service")
	cmd := "kubectl"
	args := []string{"create", "-f", kubedir + "service.yml", "--namespace=scooby"}
	if out, err := exec.Command(cmd, args...).Output(); err != nil {
		log.Debug(out)
		log.Fatal(err)
	}
}

func kubeDeployDeployment(dir string) {
	kubedir := dir + ".kube/"
	// deployment
	log.Debug("kubectl create deploy")
	cmd := "kubectl"
	args := []string{"create", "-f", kubedir + "deployment.yml", "--namespace=scooby"}
	if out, err := exec.Command(cmd, args...).Output(); err != nil {
		log.Debug(out)
		log.Fatal(err)
	}
}

func kubeUpdate(dir string) {
	log.Debug("kubectl apply")
	cmd := "kubectl"
	args := []string{"apply", "-f", dir + ".kube/deployment.yml", "--namespace=scooby"}
	if out, err := exec.Command(cmd, args...).Output(); err != nil {
		log.Debug(out)
		log.Fatal(err)
	}
}

func kubeService(subdomain, tag, port, dir string) {
	log.Debug("kubeservice")
	filepath := dir + ".kube/service.yml"
	content := []byte(`apiVersion: v1
kind: Service
metadata:
  name: ` + subdomain + `
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: ` + port + `
  selector:
    app: ` + subdomain)
	writeFile(filepath, content)
}

func kubeDeployment(subdomain, tag, port, dir string) {
	log.Debug("kubedeployment")
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
        subdomain: ` + subdomain + `
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
        - hostPort: 80
          containerPort: ` + port)
	writeFile(filepath, content)
}

func kubeConfig(subdomain, tag, dir string) {
	log.Debug("kube config")
	kubeconfig := dir + ".kube/config"
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
