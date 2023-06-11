FROM ubuntu

WORKDIR /srv

COPY ./build /srv/
RUN apt-get update \
    && apt-get install -y \
        wget \
        curl \
        fontconfig \
        libfreetype6 \
        libjpeg-turbo8 \
        libpng16-16 \
        libx11-6 \
        libxcb1 \
        libxext6 \
        libxrender1 \
        xfonts-75dpi \
        xfonts-base \
        gnupg

# install wkhtmltopdf for go-wkhtmltopdf lib usage
RUN curl "https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.jammy_amd64.deb" -L -o "wkhtmltopdf.deb"
RUN dpkg -i wkhtmltopdf.deb

# install chrome for chromedp lib usage
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \ 
    && echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list
RUN apt-get update && apt-get -y install google-chrome-stable

RUN chmod +x pdf-generator
CMD ./pdf-generator start

EXPOSE 8081
EXPOSE 8082
EXPOSE 8083
