package docker

import (
	log "github.com/Sirupsen/logrus"
	docker "github.com/fsouza/go-dockerclient"
	"os"
)

func IsEnvReady() {
	log.Debug("checking that your environment is ready [docker]")
}

func Dockerfile(image, path, localpath, containerpath string) {
	log.Debug("getting / setting Dockerfile")
	dockerfilepath := path + "Dockerfile"

	if _, err := os.Stat(dockerfilepath); os.IsNotExist(err) {
		// file doesn't exist. should we create one?
		// well. we *shouldnt*, but you know, laziness
		createDockerfile(image, path, localpath, containerpath)
	}
}

func BuildAndTagContainer(name string, path string) {
	log.Debug("build and tag container")
	log.Debug(name)

	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to get docker client. %s", err)
	}

	opts := docker.BuildImageOptions{
		Name:         name,
		ContextDir:   path,
		OutputStream: os.Stdout,
	}
	if err := client.BuildImage(opts); err != nil {
		log.Fatal(err)
	}
}

func createDockerfile(image, path, localpath, containerpath string) {
	log.Debug("createDockerfile")
	log.Debug(localpath + " to " + containerpath)
	dockerfilepath := path + "Dockerfile"
	f, err := os.Create(dockerfilepath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	content := []byte("FROM " + image + "\nADD " + localpath + " " + containerpath)
	if _, err := f.Write(content); err != nil {
		log.Fatal(content)
	}
}
