FROM node:16-alpine3.12
WORKDIR /www/coup/
COPY . .

# Set environment variables

ENV NODE_ENV production
# RUN yarn config set registry https://registry.npm.taobao.org/
RUN npm install
RUN npm run build
EXPOSE 3000
CMD [ "npm", "run", "start" ]
