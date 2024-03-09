FROM node:alpine
WORKDIR /app

COPY package-lock.json package.json ./
RUN npm i

COPY . .
CMD ["npm", "start"]
