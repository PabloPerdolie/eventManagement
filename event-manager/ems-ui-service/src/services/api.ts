import axios from 'axios';
import Cookies from 'js-cookie';
import { store } from '../store/store';
import { logout, refreshToken } from '../store/slices/authSlice';
import { authService } from './authService';

export const api = axios.create({
  baseURL: 'http://localhost:8081/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = authService.getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Типизация для элементов очереди запросов
interface QueueItem {
  resolve: (value: string | null) => void;
  reject: (error: any) => void;
}

// Флаг для отслеживания, выполняется ли обновление токена
let isRefreshing = false;
// Очередь запросов, ожидающих обновления токена
let failedQueue: QueueItem[] = [];

// Функция для обработки очереди запросов
const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  
  failedQueue = [];
};

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    // Если запрос уже был повторен или это запрос на обновление токена, который вернул ошибку,
    // то прекращаем повторные попытки
    if (originalRequest.url === '/auth/refresh' || originalRequest._retry) {
      // Если это запрос на обновление токена, выходим из системы
      if (originalRequest.url === '/auth/refresh') {
        // Сбрасываем флаг и обрабатываем очередь с ошибкой
        isRefreshing = false;
        processQueue(error, null);
        
        // Выходим из системы при ошибке обновления токена
        store.dispatch(logout());
      }
      
      return Promise.reject(error);
    }

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      // Если обновление токена уже выполняется, добавляем запрос в очередь
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(token => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            return api(originalRequest);
          })
          .catch(err => {
            return Promise.reject(err);
          });
      }
      
      isRefreshing = true;
      
      try {
        // Обновляем токен через authService напрямую
        console.log("Обновляем токен через authService...");
        const result = await authService.refreshToken();
        
        // Если обновление успешно, обновляем Authorization в заголовке
        originalRequest.headers.Authorization = `Bearer ${result.access_token}`;
        
        // Обрабатываем очередь запросов
        processQueue(null, result.access_token);
        isRefreshing = false;
        
        // Повторяем исходный запрос с новым токеном
        return api(originalRequest);
      } catch (refreshError) {
        // Сбрасываем флаг обновления
        isRefreshing = false;
        
        // Обрабатываем очередь с ошибкой
        processQueue(refreshError, null);
        
        // Выходим из системы при ошибке обновления токена
        store.dispatch(logout());
        
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);