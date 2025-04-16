import React, { useMemo } from 'react';
import { CheckCircle, Clock, AlertCircle, XCircle, ChevronDown, ChevronRight, Trash2 } from 'lucide-react';
import { TaskResponse, TaskStatus } from '../../types/api';
import { HierarchicalTask } from './types';
import { getStatusIcon, getStatusColorClass, formatDate } from './taskUtils';
import { toast } from 'react-toastify';

interface TaskItemProps {
  task: HierarchicalTask;
  level: number;
  isLastChild?: boolean;
  allTasks: TaskResponse[];
  usernames: Map<number, string>;
  animationKey: number;
  onToggleExpand: (taskId: number) => void;
  onUpdateStatus: (taskId: number, newStatus: TaskStatus) => void;
  onDelete?: (taskId: number) => Promise<void>;
}

export const TaskItem: React.FC<TaskItemProps> = ({
  task,
  level,
  isLastChild = false,
  allTasks,
  usernames,
  animationKey,
  onToggleExpand,
  onUpdateStatus,
  onDelete
}) => {
  const hasChildren = task.children && task.children.length > 0;
  
  // Вычисляем статистику по дочерним задачам один раз
  const childrenStats = useMemo(() => {
    if (!hasChildren) return { total: 0, completed: 0, percent: 0 };
    
    const total = task.children.length;
    const completed = task.children.filter(child => child.status === TaskStatus.Completed).length;
    const percent = total > 0 ? Math.round((completed / total) * 100) : 0;
    
    return { total, completed, percent };
  }, [task.children, hasChildren]);
  
  // Проверяем, можно ли пометить задачу как выполненную
  const canMarkCompleted = hasChildren ? childrenStats.completed === childrenStats.total : true;
  
  // Вычисляем следующий статус при щелчке
  const getNextStatus = (currentStatus: TaskStatus): TaskStatus => {
    switch (currentStatus) {
      case TaskStatus.Pending:
        return TaskStatus.InProgress;
      case TaskStatus.InProgress:
        return TaskStatus.Completed;
      case TaskStatus.Completed:
        return TaskStatus.Pending;
      case TaskStatus.Cancelled:
        return TaskStatus.Pending;
      default:
        return TaskStatus.Pending;
    }
  };
  
  // Отображение соединительных линий для дочерних задач
  const renderConnectingLines = () => {
    if (level === 0) return null;
    
    return (
      <>
        {/* Вертикальная линия от родителя к текущей задаче */}
        <div 
          className="absolute border-l-2 border-blue-300"
          style={{ 
            left: `${(level - 1) * 24 + 12}px`,
            top: 0,
            height: isLastChild ? '20px' : '100%',
            width: '2px',
            zIndex: 1
          }}
        ></div>
        
        {/* Горизонтальная линия к текущей задаче */}
        <div 
          className="absolute border-t-2 border-blue-300"
          style={{ 
            left: `${(level - 1) * 24 + 12}px`,
            top: '20px',
            width: '12px',
            height: '2px',
            zIndex: 1
          }}
        ></div>
      </>
    );
  };
  
  const handleStatusClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    // Предотвращаем перезагрузку страницы
    e.preventDefault();
    
    const nextStatus = getNextStatus(task.status);
    // Если следующий статус - "Выполнено", проверяем возможность изменения
    if (nextStatus === TaskStatus.Completed && !canMarkCompleted) {
      toast.warning('Сначала выполните все дочерние задачи');
    } else {
      onUpdateStatus(task.id, nextStatus);
    }
  };
  
  return (
    <div className={`task-item relative ${level === 0 ? 'mt-2' : ''}`}>
      {/* Соединительные линии */}
      {renderConnectingLines()}
      
      <div 
        className={`flex items-center p-3 hover:bg-gray-50 relative bg-white mb-2 rounded-lg 
          ${level === 0 ? 'shadow-md border-l-4 border-blue-500' : 'shadow-sm border-l-4 border-blue-300'} 
          ${hasChildren && !canMarkCompleted && task.status !== TaskStatus.Completed ? 'bg-amber-50 border-l-4 border-amber-400' : ''}
          ${hasChildren && task.isExpanded ? 'bg-blue-50' : ''}`}
        style={{ marginLeft: `${level * 24}px` }}
      >
        {/* Кнопка сворачивания/разворачивания, только если есть дочерние задачи */}
        <div className="w-8 h-8 mr-2 flex items-center justify-center">
          {hasChildren ? (
            <button 
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                onToggleExpand(task.id);
              }}
              className="rounded-full bg-blue-100 hover:bg-blue-200 w-8 h-8 flex items-center justify-center transition-colors border border-blue-300 shadow-sm hover:shadow-md animation-pulse"
              title={task.isExpanded ? "Скрыть подзадачи" : "Показать подзадачи"}
              style={{
                animation: !task.isExpanded ? 'pulse 2s infinite' : 'none'
              }}
            >
              {task.isExpanded ? 
                <ChevronDown className="h-5 w-5 text-blue-600" /> : 
                <ChevronRight className="h-5 w-5 text-blue-600" />
              }
              <span className="sr-only">{task.isExpanded ? "Скрыть подзадачи" : "Показать подзадачи"}</span>
            </button>
          ) : (
            <div className="w-8 h-8 flex items-center justify-center opacity-30">
              <div className="w-2 h-2 rounded-full bg-gray-300"></div>
            </div>
          )}
        </div>
        
        <div className="flex-1 flex items-center">
          <button 
            onClick={(e) => {
              e.preventDefault();
              e.stopPropagation();
              handleStatusClick(e);
            }}
            className="mr-3 cursor-pointer hover:opacity-80"
            title={canMarkCompleted || task.status !== TaskStatus.InProgress ? 
              `Нажмите для изменения статуса задачи` : 
              `Требуется выполнить все дочерние задачи (${task.children.filter(c => c.status !== TaskStatus.Completed).length} не выполнены)`}
          >
            {getStatusIcon(task.status)}
          </button>
          
          <div className="flex-1">
            <div className="text-sm font-medium text-gray-900 flex flex-wrap items-center">
              <span>{task.title}</span>
              <span className={`ml-2 px-1.5 py-0.5 text-xs rounded-full ${getStatusColorClass(task.status)}`}>
                {task.status}
              </span>
              {task.priority && (
                <span className="ml-2 px-1.5 py-0.5 text-xs rounded-full bg-gray-100 text-gray-800">
                  {task.priority}
                </span>
              )}
              {hasChildren && !canMarkCompleted && task.status !== TaskStatus.Completed && (
                <span className="ml-2 px-1.5 py-0.5 text-xs rounded-full bg-amber-100 text-amber-800">
                  Ожидает завершения подзадач
                </span>
              )}
              {task.parent_id && (
                <span className="text-xs ml-2 text-gray-500 italic">
                  ↳ {allTasks.find(t => t.id === task.parent_id)?.title || `Задача #${task.parent_id}`}
                </span>
              )}
              {hasChildren && (
                <span className={`text-xs ml-2 px-1.5 py-0.5 rounded-full ${canMarkCompleted ? 'bg-blue-100 text-blue-700' : 'bg-amber-100 text-amber-700'}`}>
                  {childrenStats.completed}/{childrenStats.total} выполнено
                  {!canMarkCompleted && task.status !== TaskStatus.Completed && 
                    <span className="ml-1">⚠️</span>
                  }
                </span>
              )}
              {hasChildren && task.status !== TaskStatus.Completed && childrenStats.percent < 100 && (
                <span className="text-xs ml-2 text-gray-500">
                  {childrenStats.percent}% завершено
                </span>
              )}
              {onDelete && (
                <button
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    onDelete(task.id);
                  }}
                  className="ml-auto text-red-500 hover:text-red-700 transition-colors"
                  title="Удалить задачу"
                >
                  <Trash2 size={16} />
                </button>
              )}
            </div>
            {task.description && (
              <div className="text-xs text-gray-500 mt-0.5 line-clamp-1">{task.description}</div>
            )}
            <div className="flex flex-wrap items-center mt-1 text-xs text-gray-500 gap-1">
              <span className="text-xs">
                {usernames && task.assigned_to && usernames.get(task.assigned_to) || `User #${task.assigned_to || 'Не назначен'}`}
              </span>
              {task.story_points && (
                <span className="px-1.5 py-0.5 rounded-full bg-purple-100 text-purple-800 text-xs">
                  {task.story_points}p
                </span>
              )}
              <span className="text-xs">{formatDate(task.created_at)}</span>
            </div>
          </div>
        </div>
      </div>
      
      {/* Дочерние задачи отображаются только если родительская развернута */}
      {hasChildren && task.isExpanded && (
        <div 
          key={`children-${task.id}-${animationKey}`}
          className="children relative ml-4 animate-slideDown overflow-hidden"
          style={{ animationDuration: '0.2s' }}
        >
          {task.children.map((childTask, index) => (
            <TaskItem 
              key={childTask.id} 
              task={childTask} 
              level={level + 1} 
              isLastChild={index === task.children.length - 1}
              allTasks={allTasks}
              usernames={usernames}
              animationKey={animationKey}
              onToggleExpand={onToggleExpand}
              onUpdateStatus={onUpdateStatus}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export default TaskItem; 