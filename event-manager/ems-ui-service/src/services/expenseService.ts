import { api } from './api';
import { 
  ExpenseCreateRequest, 
  ExpenseResponse, 
  ExpensesResponse, 
  ExpenseUpdateRequest, 
  BalanceReportResponse,
  ExpenseShareUpdateRequest
} from '../types/api';

export const expenseService = {
  // Получить все расходы события
  getExpenses: async (eventId: number): Promise<ExpensesResponse> => {
    try {
      const response = await api.get(`/expenses?event_id=${eventId}`);
      return response.data;
    } catch (error) {
      console.error('Ошибка при получении расходов:', error);
      throw error;
    }
  },

  // Создать новый расход
  createExpense: async (expenseData: ExpenseCreateRequest): Promise<{ expense_id: number }> => {
    try {
      const response = await api.post('/expenses', expenseData);
      return response.data;
    } catch (error) {
      console.error('Ошибка при создании расхода:', error);
      throw error;
    }
  },

  // Обновить существующий расход
  updateExpense: async (expenseId: number, expenseData: ExpenseUpdateRequest): Promise<{ success: boolean }> => {
    try {
      const response = await api.put(`/expenses/${expenseId}`, expenseData);
      return response.data;
    } catch (error) {
      console.error('Ошибка при обновлении расхода:', error);
      throw error;
    }
  },

  // Удалить расход
  deleteExpense: async (expenseId: number): Promise<{ success: boolean }> => {
    try {
      const response = await api.delete(`/expenses/${expenseId}`);
      return response.data;
    } catch (error) {
      console.error('Ошибка при удалении расхода:', error);
      throw error;
    }
  },

  // Получить отчет о балансе события
  getBalanceReport: async (eventId: number): Promise<BalanceReportResponse> => {
    try {
      const response = await api.get(`/expenses/${eventId}/balance`);
      return response.data;
    } catch (error) {
      console.error('Ошибка при получении отчета о балансе:', error);
      throw error;
    }
  },

  // Отметить долю как оплаченную
  markShareAsPaid: async (shareId: number, isPaid: boolean = true): Promise<{ success: boolean }> => {
    try {
      const payload: ExpenseShareUpdateRequest = {
        is_paid: isPaid
      };
      const response = await api.put(`/expenses/shares/${shareId}/paid-status`, payload);
      return response.data;
    } catch (error) {
      console.error('Ошибка при изменении статуса оплаты доли:', error);
      throw error;
    }
  }
}; 