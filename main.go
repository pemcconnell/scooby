package main
/*
Scooby is a very simple, very opinionated tool designed to help people deploy apps to kubernetes.

It assumes:
- the user has gcloud installed and configured
- the user wishes to create a SUBDOMAIN, on a pre-defined domain
- a single port exposed
- a pre-defined top-level domain

Specifically, this application has been built for a certain use-case:
Allow a team of developers to "push" their local prototypes online, with the minimum configuration and devops support possible

@todo:
- deploy
- update
- remove? // might be a bad idea
- auth options
*/
func main() {
	println("scooby dooby doo")
}
