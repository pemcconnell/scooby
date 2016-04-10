package main

import (
    "github.com/pemcconnell/scooby/docker"
    "github.com/pemcconnell/scooby/gcloud"
    "github.com/pemcconnell/scooby/kubectl"

    log "github.com/Sirupsen/logrus"
    "github.com/codegangsta/cli"

    "os"
    "path/filepath"
    //    "time"
)

type Config struct {
    dir        string // app directory
    tld        string // top-level domain
    sub        string // subdomain
    port       string // container port
    hubbase    string // hub-base
    version    string // docker tag version
    image      string // image used for dockerfile
    dockerpath string // dest path in container
}

/*
usage:
scooby deploy ./app -to=thissubdomain
*/

func init() {
    log.SetOutput(os.Stderr)
    log.SetLevel(log.DebugLevel)
}

func main() {
    log.Debug("running")

    var (
        subdomain string
    )

    isEnvironmentReady()

    app := cli.NewApp()
    app.Name = "scooby"
    app.Usage = "a super simple, opinionated way to deploy " +
        "web apps to kubernetes"
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:        "to",
            Usage:       "the desired subdomain to deploy to",
            Destination: &subdomain,
        },
    }
    config := config(subdomain)
    app.Commands = []cli.Command{
        {
            Name:  "deploy",
            Usage: "deploy an app",
            Action: func(c *cli.Context) {
                log.Debug("deploy")
                if subdomain == "" {
                    log.Fatal("You must specify a subdomain. Use -to=")
                }
                config.dir = filepath.Clean(c.Args().First())
                if config.dir != "/" {
                    config.dir = config.dir + "/"
                }
                tag := config.hubbase + config.sub + ":" + config.version

                docker.Dockerfile(config.image, config.dir, config.dockerpath)
                docker.BuildAndTagContainer(tag, config.dir)
                gcloud.PushContainer(tag)
                kubectl.Deploy(subdomain, tag, config.port, config.dir)

                generateScoobyFiles()
                presentData()
            },
        },
    }
    app.Run(os.Args)
}

func isEnvironmentReady() {
    log.Debug("checking that your environment is ready [general]")
    docker.IsEnvReady()
    gcloud.IsEnvReady()
    kubectl.IsEnvReady()
}

func config(subdomain string) Config {
    log.Debug("checking that the config is ok")
    config := Config{
        tld:        "myprototypes.com",
        port:       "80",
        hubbase:    "gcr.io/something",
        image:      "nginx:stable-alpine",
        version:    "latest", //time.Now().Format("20060102150405"),
        sub:        subdomain,
        dockerpath: "/usr/share/nginx/html",
    }
    return config
}

func generateScoobyFiles() {
    log.Debug("generating scooby files, if needed")
}

func presentData() {
    log.Debug("presenting data back to user")
}
