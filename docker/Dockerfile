FROM dind
MAINTAINER mihail0zdravkov@gmail.com

# Install Nginx.

RUN apt-get update
RUN apt-get install -qqy python-software-properties software-properties-common

RUN \
  add-apt-repository -y ppa:nginx/stable && \
  apt-get install -y nginx && \
  rm -rf /var/lib/apt/lists/* && \
  echo "\ndaemon off;" >> /etc/nginx/nginx.conf && \
  chown -R www-data:www-data /var/lib/nginx

# Define mountable directories.
VOLUME ["/etc/nginx/sites-enabled", "/etc/nginx/certs", "/etc/nginx/conf.d", "/var/log/nginx"]

# Expose ports.
EXPOSE 80
EXPOSE 443

# Based on the Docker.io's sshd example

RUN apt-get update && apt-get install -y openssh-server
RUN mkdir /var/run/sshd
RUN echo 'root:verysecure' | chpasswd
RUN sed -i 's/PermitRootLogin without-password/PermitRootLogin yes/' /etc/ssh/sshd_config

EXPOSE 22

ADD https://github.com/mzdravkov/kamino/blob/master/kamino?raw=true /kamino
RUN chmod +x /kamino

ADD ./config.cfg /config.cfg

ADD ./dockerstart /dockerstart
RUN chmod +x /dockerstart

RUN mkdir /myr
RUN mkdir /myr/configs
RUN mkdir /myr/plugins
RUN touch /myr/tenant_conf.default.yml

CMD ["/dockerstart"]
