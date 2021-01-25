FROM nginx:1.19.0-alpine
#configuration
#content, comment out the ones you dont need!
#COPY ./*.html /usr/share/nginx/html/
RUN chmod 666 -R /etc/nginx/
COPY . /app
#RUN rm  -R /etc/nginx/conf.d/
#RUN ls
ADD ./html/ /usr/share/nginx/html/
COPY ./nginx/default.conf /etc/nginx/conf.d/default.conf