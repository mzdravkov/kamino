Kamino
======

A tool for deploying web applications using [Docker](http://docker.io) and [nginx](http://nginx.org).

**NOTE:** This branch is currently in a non-functional state, because I'm working on making kamino able to work in a distributed manner. You can find the old functional version in branch ``release-0.0.1``.

#About

Kamino is a command line tool, written in [Go language](http://golang.org), for making deployment of many instances of one web application easy. It's small and simple.

Currently, Kamino is able to be given a docker image and then deploy many docker containers from that image. For example, you have built some kick-ass blog application in your favourite language (Kamino is language agnostic), for example Cobol (just kidding :P). Then, you [make a docker image](http://docs.docker.io/en/latest/use/builder/) with your application and all of it's dependencies and give the name of the image in the config.cfg file (and make few other configurations there). And voila, Kamino can now deploy as many "clones" of this web application as you'd like.

Kamino is created, with the very idea to be used with [beyond](http://github.com/mzdravkov/beyond). Beyond is a Rails engine that, after plugged to rails application, will provide to that application the ability create and configure new tenants (clones of the template application) and manage their plugins.

Kamino will take care of configuring and reloading your nginx. After deploying tenant called "llama", you will have llama.yourhost.com path that points to your application.

#Installation
Assuming that [Go language](http://golang.org), [Docker](http://docker.io) and [nginx](http://nginx.org) are installed:

``git clone https://github.com/mzdravkov/kamino``

``cd kamino``

``go build``

Note: In order to work, the user running kamino should have rights to reconfigure/reload nginx and use the Docker.

#Usage
You can see example configuration at config.cfg. The non-obvious configurations are:

`tenants_configs_dir` Path to the directory where all user configurations resides. Each configuration is mounted into the file system of it's corresponding Docker container.

`tenants_port` The port, on which the template web application is listening inside each container. (Each tenant is a web application running in separate Docker container. The port, on which it's listening (inside the container) is exposed and mapped to a host system port. All web apps listen on the same port (inside the container), but each one of them has it's port mapped to a different host system port. When a request to tenant.domain.tld is made, nginx forwards it to the host system port for the tenant. The host system port is mapped to the `tenants_port` port inside the tenant container, so the request is processed by the web application.)

`tenants_config_path` The path, to which the configuration file of the tenant will be mounted inside it's file system.

`tenants_default_config` Path to the default config for all tenants. On each tenant creation, the system crates it's own configuration file, by copying the default one.

`tenants_plugins_dir` Similar to `tenants_configs_dir`. The path to the directory with plugins for tenants.

`tenants_plugins_path` Similar to `tenants_config_path`. The path, to which the plugins of the tenant (they resides in `tenants_configs_dir/tenant_name`) will be mounted inside the tenant container's file system.


To deploy new tenant:

`kamino deploy -name='llama'`


You can stop tenants by just stopping their container (`docker stop name`) and start them (after they've been stopped) with `docker start name`
