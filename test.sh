#!/bin/bash

docker run -t --rm -v "$PWD":/usr/src/app -p "4000:4000" starefossen/github-pages

#docker run -p 3000:3000 -e github='https://github.com/pmcdowell-okta/dockertest.git' -it oktaadmin/dockertest