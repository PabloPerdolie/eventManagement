import { ExpenseResponse, TaskResponse } from '../../types/api';

export interface TaskExpense {
  expense_id: number;
  description: string;
  amount: number;
  currency: string;
}

export interface HierarchicalTask extends TaskResponse {
  children: HierarchicalTask[];
  isExpanded?: boolean;
  expenses?: TaskExpense[]; // Расходы, связанные с задачей
  totalExpenses?: number; // Общая сумма расходов включая дочерние задачи
}

export interface TaskChildrenStats {
  total: number;
  completed: number;
  percent: number;
} 