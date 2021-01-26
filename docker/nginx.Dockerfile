FROM nginx:1.18.0-alpine
#configuration
#content, comment out the ones you dont need!
RUN chmod 666 -R /etc/nginx/

COPY ./nginx/ /etc/nginx/conf.d/

COPY . /app
