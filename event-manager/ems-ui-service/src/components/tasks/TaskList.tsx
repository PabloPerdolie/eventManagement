import React, { useState } from 'react';
import { TaskResponse, TaskStatus } from '../../types/api';
import { HierarchicalTask } from './types';
import TaskItem from './TaskItem';

interface TaskListProps {
  hierarchicalTasks: HierarchicalTask[];
  tasks: TaskResponse[];
  usernames: Map<number, string>;
  onUpdateStatus: (taskId: number, newStatus: TaskStatus) => Promise<void>;
  onDelete?: (taskId: number) => Promise<void>;
  expandedTaskIds: Set<number>;
  setExpandedTaskIds: React.Dispatch<React.SetStateAction<Set<number>>>;
}

const TaskList: React.FC<TaskListProps> = ({
  hierarchicalTasks,
  tasks,
  usernames,
  onUpdateStatus,
  onDelete,
  expandedTaskIds,
  setExpandedTaskIds
}) => {
  // Состояние для обновления анимации
  const [animationKey, setAnimationKey] = useState(0);

  const handleToggleExpand = (taskId: number) => {
    setExpandedTaskIds(prev => {
      const newSet = new Set(prev);
      if (newSet.has(taskId)) {
        newSet.delete(taskId);
      } else {
        newSet.add(taskId);
      }
      return newSet;
    });
    // Увеличиваем ключ анимации при каждом переключении
    setAnimationKey(prev => prev + 1);
  };

  // Преобразуем HierarchicalTask с isExpanded
  const addExpandedState = (tasks: HierarchicalTask[]): HierarchicalTask[] => {
    return tasks.map(task => ({
      ...task,
      isExpanded: expandedTaskIds.has(task.id),
      children: task.children ? addExpandedState(task.children) : []
    }));
  };

  // Добавляем состояние isExpanded к задачам
  const tasksWithExpandedState = addExpandedState(hierarchicalTasks);

  return (
    <div className="task-list">
      {tasksWithExpandedState.length === 0 ? (
        <div className="bg-gray-50 p-8 text-center rounded-lg border border-gray-200">
          <p className="text-gray-500">Нет задач для отображения</p>
        </div>
      ) : (
        tasksWithExpandedState.map(task => (
          <TaskItem
            key={task.id}
            task={task}
            level={0}
            allTasks={tasks}
            usernames={usernames}
            animationKey={animationKey}
            onToggleExpand={handleToggleExpand}
            onUpdateStatus={onUpdateStatus}
            onDelete={onDelete}
          />
        ))
      )}
    </div>
  );
};

export default TaskList; 