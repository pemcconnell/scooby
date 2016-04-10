package gcloud

/*
@todo: further investigation on the gcloud go api to see if
it's possible to replace cli calls
*/

import (
	log "github.com/Sirupsen/logrus"

	"fmt"
	"os"
	"os/exec"
)

func IsEnvReady() {
	log.Debug("checking that your environment is ready [gcloud]")
}

func PushContainer(name string) {
	log.Debug("pushing the container onto gcr.io")
	cmd := "gcloud"
	args := []string{"docker", "push", name}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
