FROM node:latest as build

COPY package.json package-lock.json /

RUN npm install

COPY . .
RUN npm run build

FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build ./build /usr/share/nginx/html
