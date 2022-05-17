# wkhtmltopdf-golang-docker

This repo demonstrates how to use wkhtmltopdf to convert an HTML page to a PDF file. The code is written in Go and uses a golang wrapper for wkhtmltopdf.

This is a production ready example. Feel free to use it in your own projects.

## Running locally

First, clone this repo

```
git clone https://github.com/theterminalguy/wkhtmltopdf-golang-docker
```

Next, build the image with

```
docker build -t wkhtmltopdf-golang-docker .
```

You can then run the image locally with

```
docker run -p 3000:3000 --env PORT=3000 --rm wkhtmltopdf-golang-docker:latest
```
