FROM node:18-alpine

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN npm install --include=dev
RUN npm run build

CMD ["node", ".output/server/index.mjs"]



