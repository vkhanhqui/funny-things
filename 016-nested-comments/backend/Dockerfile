FROM node:22.16.0-alpine
WORKDIR /app

COPY . .
COPY .env.dist .env

RUN npm install
RUN npx prisma generate

CMD ["npm", "start"]