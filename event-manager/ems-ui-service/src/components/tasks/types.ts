import { TaskResponse } from '../../types/api';

export interface HierarchicalTask extends TaskResponse {
  children: HierarchicalTask[];
  isExpanded?: boolean;
}

export interface TaskChildrenStats {
  total: number;
  completed: number;
  percent: number;
} 