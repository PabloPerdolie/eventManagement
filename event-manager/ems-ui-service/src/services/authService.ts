import axios from 'axios';
import Cookies from 'js-cookie';
import { AuthResponse, UserResponse } from '../types/api';

const API_URL = 'http://localhost:8081/api/v1';

/**
 * Сервис для управления аутентификацией
 */
export const authService = {
  /**
   * Вход в систему
   */
  login: async (username: string, password: string): Promise<AuthResponse> => {
    try {
      const response = await axios.post<AuthResponse>(`${API_URL}/auth/login`, {
        username,
        password
      });
      
      const { access_token, refresh_token, user } = response.data;
      
      authService.setTokens(access_token, refresh_token);
      authService.setUser(user);
      
      return response.data;
    } catch (error) {
      console.error('Ошибка при входе в систему:', error);
      throw error;
    }
  },

  /**
   * Регистрация нового пользователя
   */
  register: async (username: string, email: string, password: string): Promise<AuthResponse> => {
    try {
      const response = await axios.post<AuthResponse>(`${API_URL}/auth/register`, {
        username,
        email,
        password
      });
      
      const { access_token, refresh_token, user } = response.data;
      
      authService.setTokens(access_token, refresh_token);
      authService.setUser(user);
      
      return response.data;
    } catch (error) {
      console.error('Ошибка при регистрации:', error);
      throw error;
    }
  },

  /**
   * Выход из системы
   */
  logout: async (): Promise<void> => {
    try {
      const accessToken = authService.getAccessToken();
      
      if (accessToken) {
        await axios.post(`${API_URL}/auth/logout`, {}, {
          headers: {
            Authorization: `Bearer ${accessToken}`
          }
        });
      }
    } catch (error) {
      console.warn('Ошибка при выходе из системы:', error);
      // Продолжаем выход даже при ошибке
    } finally {
      authService.clearTokens();
    }
  },

  /**
   * Обновление токена доступа
   */
  refreshToken: async (): Promise<{ access_token: string; refresh_token: string }> => {
    try {
      const refreshToken = authService.getRefreshToken();
      
      if (!refreshToken) {
        throw new Error('Refresh token not found');
      }
      
      console.log(`Отправка запроса на обновление токена, refresh token: ${refreshToken.substring(0, 5)}...`);
      
      // Создаем новый экземпляр axios для запроса обновления токена
      // чтобы избежать циклических вызовов через интерцепторы
      const response = await axios.post<{ access_token: string; refresh_token: string; expires_in: number }>
      (`${API_URL}/auth/refresh`, { refresh_token: refreshToken });
      
      const { access_token, refresh_token } = response.data;
      
      console.log('Токен успешно обновлен');
      authService.setTokens(access_token, refresh_token);
      
      return { access_token, refresh_token };
    } catch (error) {
      console.error('Ошибка при обновлении токена:', error);
      authService.clearTokens();
      throw error;
    }
  },

  /**
   * Получение текущего пользователя
   */
  getCurrentUser: async (): Promise<UserResponse> => {
    const userFromCookies = authService.getUserFromCookies();
    
    if (userFromCookies) {
      return userFromCookies;
    }
    
    try {
      const accessToken = authService.getAccessToken();
      
      if (!accessToken) {
        throw new Error('Access token not found');
      }
      
      const response = await axios.get<UserResponse>(`${API_URL}/auth/me`, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      
      const user = response.data;
      authService.setUser(user);
      
      return user;
    } catch (error) {
      console.error('Ошибка при получении данных пользователя:', error);
      throw error;
    }
  },

  // Вспомогательные методы для работы с токенами и данными пользователя
  
  /**
   * Сохранение токенов в куки
   */
  setTokens: (accessToken: string, refreshToken: string): void => {
    Cookies.set('access_token', accessToken);
    Cookies.set('refresh_token', refreshToken);
  },

  /**
   * Очистка токенов из куки
   */
  clearTokens: (): void => {
    Cookies.remove('access_token');
    Cookies.remove('refresh_token');
    Cookies.remove('user_data');
  },

  /**
   * Получение access токена
   */
  getAccessToken: (): string | undefined => {
    return Cookies.get('access_token');
  },

  /**
   * Получение refresh токена
   */
  getRefreshToken: (): string | undefined => {
    return Cookies.get('refresh_token');
  },

  /**
   * Сохранение пользователя в куки
   */
  setUser: (user: UserResponse): void => {
    Cookies.set('user_data', JSON.stringify(user));
  },

  /**
   * Получение пользователя из куки
   */
  getUserFromCookies: (): UserResponse | null => {
    try {
      const userData = Cookies.get('user_data');
      if (!userData) return null;
      return JSON.parse(userData) as UserResponse;
    } catch (error) {
      console.error('Ошибка при получении данных пользователя из куки:', error);
      return null;
    }
  },

  /**
   * Проверка, авторизован ли пользователь
   */
  isAuthenticated: (): boolean => {
    return !!authService.getAccessToken();
  }
}; 