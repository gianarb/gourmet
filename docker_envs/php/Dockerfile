FROM php:7.0.0
MAINTAINER Gianluca Arbezzano <gianarb92@gmail.com>

RUN apt-get update
RUN apt-get install -y curl git \
    unzip wget \
    libfreetype6-dev \
    libcurl3-dev \
    libjpeg62-turbo-dev \
    libmcrypt-dev 

RUN docker-php-ext-install mcrypt

CMD ["php", "-v"]
