FROM ubuntu:16.04
RUN mkdir /app
RUN apt-get update -qq && apt-get install -y supervisor
COPY ./docker-files/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
ADD PowerController /app
RUN chmod +x /app/PowerController
CMD ["/usr/bin/supervisord"]
