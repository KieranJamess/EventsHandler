FROM node:20

WORKDIR /app

COPY package*.json ./
COPY vite.config.ts ./

RUN npm install

# Change the below to pull in source files before build and build / preview from that
COPY dist/ ./build/

EXPOSE 8080

CMD ["npm", "run", "preview"]
