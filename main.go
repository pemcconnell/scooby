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
}

func main() {
    log.Debug("running")

    var (
        subdomain string
        namespace string
        project   string
        port      string
    )

    isEnvironmentReady()

    app := cli.NewApp()
    app.Name = "scooby"
    app.Usage = "a super simple, opinionated way to deploy " +
        "web apps to kubernetes"
    app.Version = "0.0.1"
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:        "project, pr",
            Usage:       "gcloud project name",
            EnvVar:      "SCOOBY_PROJECT",
            Destination: &project,
        },
        cli.StringFlag{
            Name:        "namespace, n",
            Usage:       "set the kubernetes namespace to use",
            EnvVar:      "SCOOBY_NAMESPACE",
            Destination: &namespace,
        },
    }
    config := config()
    app.Commands = []cli.Command{
        {
            Name:  "deploy",
            Usage: "deploy an app",
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:        "to, t",
                    Usage:       "the desired subdomain to deploy to",
                    Destination: &subdomain,
                },
                cli.StringFlag{
                    Name:        "port, p",
                    Usage:       "container port",
                    Value:       "80",
                    Destination: &port,
                },
            },
            Action: func(c *cli.Context) {
                log.Debug("deploy")
                if subdomain == "" {
                    log.Fatal("You must specify a subdomain. Use -t=")

                }
                if project == "" {
                    log.Fatal("You must specify a project. Use -p=, or the " +
                        "SCOOBY_PROJECT environment variable")
                }
                config.port = port
                config.hubbase = "gcr.io/" + project + "/"
                config.dir = filepath.Clean(c.Args().First())
                if config.dir != "/" {
                    config.dir = config.dir + "/"
                }
                tag := config.hubbase + subdomain + config.sub + ":" + config.version

                docker.Dockerfile(config.image, config.dir, config.localpath, config.dockerpath)
                docker.BuildAndTagContainer(tag, config.dir)
                gcloud.PushContainer(tag)
                kubectl.Deploy(namespace, subdomain, tag, config.port, config.dir)

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

func config() Config {
    log.Debug("checking that the config is ok")
    config := Config{
        image:      "nginx:stable-alpine",
        version:    time.Now().Format("20060102150405"),
        dockerpath: "/usr/share/nginx/html",
        localpath:  ".",
    }
    return config
}

func presentData() {
    log.Debug("presenting data back to user")
}
