#!/bin/bash
web_path='/var/www/massgrid.cn/'
echo "start development"
cd ~/web/massgrid/nodeweb && \
git reset --hard origin/develop && \
git clean -f && \
git pull && \
npm update --registry=https://registry.npm.taobao.org && \
npm run build:prod && \
cp -r dist/* $web_path && \
echo "end development"