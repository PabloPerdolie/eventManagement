import React from 'react';
import { TaskResponse, TaskStatus } from '../../types/api';
import { HierarchicalTask } from './types';
import { CheckCircle, Clock, AlertCircle, XCircle } from 'lucide-react';

// Получение иконки статуса
export const getStatusIcon = (status: TaskStatus): React.ReactNode => {
  switch (status) {
    case TaskStatus.Pending:
      return <Clock className="h-5 w-5 text-yellow-500" />;
    case TaskStatus.InProgress:
      return <AlertCircle className="h-5 w-5 text-blue-500" />;
    case TaskStatus.Completed:
      return <CheckCircle className="h-5 w-5 text-green-500" />;
    case TaskStatus.Cancelled:
      return <XCircle className="h-5 w-5 text-red-500" />;
    default:
      return <Clock className="h-5 w-5 text-gray-500" />;
  }
};

// Получение CSS-класса для цвета статуса
export const getStatusColorClass = (status: TaskStatus): string => {
  switch (status) {
    case TaskStatus.Pending:
      return 'bg-yellow-100 text-yellow-800';
    case TaskStatus.InProgress:
      return 'bg-blue-100 text-blue-800';
    case TaskStatus.Completed:
      return 'bg-green-100 text-green-800';
    case TaskStatus.Cancelled:
      return 'bg-red-100 text-red-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

// Форматирование даты
export const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
};

// Проверка, можно ли изменить статус задачи на определенный
export const canChangeTaskStatus = (task: TaskResponse, newStatus: TaskStatus, allTasks: TaskResponse[]): boolean => {
  // Если пытаемся пометить задачу как выполненную
  if (newStatus === TaskStatus.Completed) {
    // Проверяем, есть ли у задачи невыполненные дочерние задачи
    const childTasks = allTasks.filter(t => t.parent_id === task.id);
    
    // Если у задачи есть дочерние задачи
    if (childTasks.length > 0) {
      // Проверяем, все ли дочерние задачи уже выполнены
      const allChildrenCompleted = childTasks.every(t => t.status === TaskStatus.Completed);
      
      // Если не все дочерние задачи выполнены, нельзя выполнить родительскую
      return allChildrenCompleted;
    }
  }
  
  // Для остальных статусов ограничений нет
  return true;
};

// Преобразование плоского списка задач в иерархическую структуру
export const buildHierarchicalTasks = (tasksList: TaskResponse[]): HierarchicalTask[] => {
  // Нормализация parent_id
  const normalizeParentId = (task: TaskResponse): number | null => {
    if (task.parent_id === null || task.parent_id === undefined) return null;
    
    if (typeof task.parent_id === 'number') {
      return task.parent_id === 0 ? null : task.parent_id;
    }
    
    if (typeof task.parent_id === 'string') {
      const parentIdStr = task.parent_id as string;
      if (parentIdStr.trim() === '') return null;
      
      const numericId = parseInt(parentIdStr, 10);
      return isNaN(numericId) ? null : numericId;
    }
    
    return null;
  };
  
  // Предотвращение циклических зависимостей
  const checkForCycles = (tasks: TaskResponse[]): TaskResponse[] => {
    const result: TaskResponse[] = [];
    const visited = new Set<number>();
    
    tasks.forEach(task => {
      visited.clear();
      
      let currentTask = task;
      let currentParentId = normalizeParentId(currentTask);
      let hasCycle = false;
      
      visited.add(currentTask.id);
      
      while (currentParentId !== null && !hasCycle) {
        if (visited.has(currentParentId)) {
          hasCycle = true;
          break;
        }
        
        const parentTask = tasks.find(t => t.id === currentParentId);
        
        if (!parentTask) break;
        
        visited.add(currentParentId);
        
        currentTask = parentTask;
        currentParentId = normalizeParentId(currentTask);
      }
      
      if (hasCycle) {
        result.push({ ...task, parent_id: null });
      } else {
        result.push({ ...task, parent_id: normalizeParentId(task) });
      }
    });
    
    return result;
  };
  
  // Проверяем циклические зависимости и нормализуем parent_id
  const normalizedTasks = checkForCycles(tasksList).map(task => ({
    ...task,
    parent_id: normalizeParentId(task)
  }));
  
  // Создаем карту задач для быстрого доступа
  const taskMap = new Map<number, HierarchicalTask>();
  
  // Инициализируем иерархическую структуру для всех задач
  normalizedTasks.forEach(task => {
    taskMap.set(task.id, { 
      ...task, 
      children: [], 
      isExpanded: false
    });
  });
  
  // Разделяем корневые задачи и дочерние задачи
  const rootTasks: HierarchicalTask[] = [];
  const orphanedTasks: HierarchicalTask[] = [];
  
  // Распределяем задачи на корневые и дочерние
  normalizedTasks.forEach(task => {
    const hierarchicalTask = taskMap.get(task.id);
    if (!hierarchicalTask) return;
    
    if (task.parent_id === null || task.parent_id === undefined || task.parent_id === 0) {
      rootTasks.push(hierarchicalTask);
    } else {
      if (task.parent_id === task.id) {
        rootTasks.push(hierarchicalTask);
        return;
      }
      
      const parentTask = taskMap.get(task.parent_id);
      
      if (parentTask) {
        parentTask.children.push(hierarchicalTask);
      } else {
        orphanedTasks.push(hierarchicalTask);
      }
    }
  });
  
  // Порядок сортировки по статусу
  const statusOrder = {
    [TaskStatus.Pending]: 0,
    [TaskStatus.InProgress]: 1,
    [TaskStatus.Completed]: 2,
    [TaskStatus.Cancelled]: 3
  };
  
  // Функция сортировки по статусу и заголовку
  const sortByStatusAndTitle = (a: HierarchicalTask, b: HierarchicalTask) => {
    if (a.status === b.status) {
      return a.title.localeCompare(b.title);
    }
    return statusOrder[a.status] - statusOrder[b.status];
  };
  
  // Рекурсивная сортировка дочерних задач
  const sortChildrenRecursive = (task: HierarchicalTask): HierarchicalTask => {
    if (task.children.length > 0) {
      const sortedChildren = task.children
        .map(sortChildrenRecursive)
        .sort(sortByStatusAndTitle);
      
      return { ...task, children: sortedChildren };
    }
    return task;
  };
  
  // Автоматическое раскрытие задач с активными подзадачами
  const updateTaskExpansionState = (task: HierarchicalTask): HierarchicalTask => {
    if (task.children.length > 0) {
      const updatedChildren = task.children.map(updateTaskExpansionState);
      
      const hasActiveChildren = updatedChildren.some(
        child => child.status === TaskStatus.Completed || child.status === TaskStatus.InProgress
      );
      
      return {
        ...task,
        children: updatedChildren,
        isExpanded: task.status === TaskStatus.Completed || hasActiveChildren
      };
    }
    
    return task;
  };
  
  // Сортируем корневые задачи
  rootTasks.sort(sortByStatusAndTitle);
  
  // Обрабатываем корневые задачи и устанавливаем автоматическое раскрытие
  const processedRootTasks = rootTasks
    .map(sortChildrenRecursive)
    .map(updateTaskExpansionState);
  
  // Сортируем "осиротевшие" задачи
  orphanedTasks.sort(sortByStatusAndTitle);
  
  // Комбинируем все задачи
  return [...processedRootTasks, ...orphanedTasks];
}; 