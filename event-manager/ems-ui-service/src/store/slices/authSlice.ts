import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { AuthResponse, UserResponse } from '../../types/api';
import Cookies from 'js-cookie';
import { authService } from '../../services/authService';

interface AuthState {
  user: UserResponse | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

// Обновляем начальное состояние
const initialState: AuthState = {
  user: authService.getUserFromCookies(),
  isAuthenticated: authService.isAuthenticated(),
  loading: false,
  error: null,
};

export const getCurrentUser = createAsyncThunk(
  'auth/getCurrentUser',
  async (_, { rejectWithValue }) => {
    try {
      return await authService.getCurrentUser();
    } catch (error) {
      console.error('Ошибка при получении данных текущего пользователя:', error);
      return rejectWithValue('Не удалось получить данные пользователя');
    }
  }
);

export const refreshToken = createAsyncThunk(
  'auth/refreshToken',
  async (_, { rejectWithValue }) => {
    try {
      return await authService.refreshToken();
    } catch (error) {
      console.error('Ошибка при обновлении токена:', error);
      return rejectWithValue('Не удалось обновить токен');
    }
  }
);

export const login = createAsyncThunk(
  'auth/login',
  async (credentials: { username: string; password: string }, { rejectWithValue }) => {
    try {
      const response = await authService.login(credentials.username, credentials.password);
      return response.user;
    } catch (error) {
      console.error('Ошибка при входе в систему:', error);
      return rejectWithValue('Не удалось войти в систему. Проверьте имя пользователя и пароль.');
    }
  }
);

export const register = createAsyncThunk(
  'auth/register',
  async (data: { username: string; email: string; password: string }, { rejectWithValue }) => {
    try {
      const response = await authService.register(data.username, data.email, data.password);
      return response.user;
    } catch (error) {
      console.error('Ошибка при регистрации:', error);
      return rejectWithValue('Не удалось зарегистрироваться. Возможно, это имя пользователя или email уже используются.');
    }
  }
);

export const logout = createAsyncThunk(
  'auth/logout', 
  async (_, { rejectWithValue }) => {
    try {
      await authService.logout();
      return null;
    } catch (error) {
      console.error('Ошибка при выходе из системы:', error);
      return rejectWithValue('Ошибка при выходе из системы');
    }
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearAuthError: (state) => {
      state.error = null;
    }
  },
  extraReducers: (builder) => {
    builder
      // Login actions
      .addCase(login.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.user = action.payload;
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string || 'Ошибка входа в систему';
      })
      
      // Register actions
      .addCase(register.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(register.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.user = action.payload;
      })
      .addCase(register.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string || 'Ошибка регистрации';
      })
      
      // Get current user actions
      .addCase(getCurrentUser.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(getCurrentUser.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.user = action.payload;
      })
      .addCase(getCurrentUser.rejected, (state) => {
        state.loading = false;
        state.isAuthenticated = false;
        state.user = null;
      })
      
      // Refresh token actions
      .addCase(refreshToken.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(refreshToken.fulfilled, (state) => {
        state.loading = false;
        state.isAuthenticated = true;
      })
      .addCase(refreshToken.rejected, (state, action) => {
        state.loading = false;
        state.isAuthenticated = false;
        state.user = null;
        state.error = action.payload as string || 'Не удалось обновить токен';
      })
      
      // Logout actions
      .addCase(logout.fulfilled, (state) => {
        state.user = null;
        state.isAuthenticated = false;
      });
  },
});

export const { clearAuthError } = authSlice.actions;
export default authSlice.reducer;