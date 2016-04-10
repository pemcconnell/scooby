package main

import (
    "github.com/pemcconnell/scooby/docker"
    "github.com/pemcconnell/scooby/gcloud"
    "github.com/pemcconnell/scooby/kubectl"

    log "github.com/Sirupsen/logrus"
    "github.com/codegangsta/cli"

    "os"
    "path/filepath"
    "time"
)

type Config struct {
    dir        string // app directory
    sub        string // subdomain
    port       string // container port
    hubbase    string // hub-base
    version    string // docker tag version
    image      string // image used for dockerfile
    dockerpath string // dest path in container
    localpath  string // relative path to add
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
        project   string
    )

    isEnvironmentReady()

    app := cli.NewApp()
    app.Name = "scooby"
    app.Usage = "a super simple, opinionated way to deploy " +
        "web apps to kubernetes"
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:        "to, t",
            Usage:       "the desired subdomain to deploy to",
            Destination: &subdomain,
        },
        cli.StringFlag{
            Name:        "project, p",
            Usage:       "gcloud project name",
            EnvVar:      "SCOOBY_PROJECT",
            Destination: &project,
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
                    log.Fatal("You must specify a subdomain. Use -t=")
                }
                if project == "" {
                    log.Fatal("You must specify a project. Use -p=, or the " +
                        "SCOOBY_PROJECT environment variable")
                }
                config.hubbase = "gcr.io/" + project + "/"
                config.dir = filepath.Clean(c.Args().First())
                if config.dir != "/" {
                    config.dir = config.dir + "/"
                }
                tag := config.hubbase + subdomain + config.sub + ":" + config.version

                docker.Dockerfile(config.image, config.dir, config.localpath, config.dockerpath)
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
        port:       "80",
        image:      "nginx:stable-alpine",
        version:    time.Now().Format("20060102150405"),
        sub:        subdomain,
        dockerpath: "/usr/share/nginx/html",
        localpath:  ".",
    }
    return config
}

func generateScoobyFiles() {
    log.Debug("generating scooby files, if needed")
}

func presentData() {
    log.Debug("presenting data back to user")
}
