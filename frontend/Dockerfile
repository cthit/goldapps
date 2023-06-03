FROM node:latest as build

COPY . .
ENV NODE_ENV=development
RUN npm install
ENV NODE_ENV=production
RUN npm run build
RUN rm -rf node_modules
RUN npm install

EXPOSE 3000
CMD npm run start
