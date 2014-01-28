Kamino
======

A tool for deploying web applications using Docker and Nginx.

#About

Kamino is a command line tool, written in Go, for making deployment of many instances of one web application easy. It's small and simple.

Currently, Kamino is able to be given a docker image and then deploy many docker containers from that image. For example, you have built some kick-ass blog application in your favourite language (Kamino can deploy applications, written in any language), for example Cobol (just kidding :P). Then, you [make a docker image](http://docs.docker.io/en/latest/use/builder/) with your application and all of it's dependencies and give the name of the image in the config.cfg file (and make few other configurations there). And voila, Kamino can now deploy as much "clones" of this web application as you'd like.

Kamino is created, with the very idea to be used with [Myrmidon](http://github.com/mzdravkov/myrmidon) or any other similar application. Myrmidon is a Rails application created to serve as a frontend for Kamino. Myrmidon provides the ability to it's registered user to create tenants, configure them, manage their plugins, etc.

Kamino will take care of configuring and reloading your nginx. After deploying tenant called "llama", you will have yourhost.com/llama path that points to your application.

#How does it work