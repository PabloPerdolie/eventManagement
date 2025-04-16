FROM node:18-alpine

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

RUN npm run build

EXPOSE 3000

# Используем служебный http-сервер для раздачи статики
RUN npm install -g serve

# Запускаем serve для раздачи собранного Vite приложения (директория dist)
CMD ["serve", "-s", "dist", "-l", "3000"] 