FROM node:16-alpine3.12
WORKDIR /www/coup/
COPY . .

# Set environment variables

ENV NODE_ENV production

RUN yarn config set registry 'https://registry.npm.taobao.org'
RUN yarn install
RUN yarn build
EXPOSE 3000
CMD [ "yarn", "start" ]