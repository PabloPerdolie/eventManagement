import React, { useState } from 'react';
import { X } from 'lucide-react';
import { TaskStatus, TaskCreateRequest } from '../../types/api';

interface AddTaskFormProps {
  eventId: number;
  tasks: Array<{ id: number; title: string; status: TaskStatus }>;
  participants: Array<{ id: number; username: string }>;
  onCancel: () => void;
  onSubmit: (taskData: TaskCreateRequest) => Promise<void>;
}

const AddTaskForm: React.FC<AddTaskFormProps> = ({
  eventId,
  tasks,
  participants,
  onCancel,
  onSubmit
}) => {
  const [newTask, setNewTask] = useState<{
    title: string;
    description: string;
    assigned_to: number | null;
    parent_id: number | null;
    priority: string;
    story_points: number | null;
  }>({
    title: '',
    description: '',
    assigned_to: null,
    parent_id: null,
    priority: 'medium',
    story_points: null
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newTask.title.trim()) {
      return;
    }
    
    const taskData: TaskCreateRequest = {
      title: newTask.title,
      event_id: eventId,
      description: newTask.description || undefined,
      assigned_to: newTask.assigned_to || undefined,
      parent_id: newTask.parent_id || undefined,
      priority: newTask.priority || undefined,
      story_points: newTask.story_points || undefined
    };
    
    await onSubmit(taskData);
  };

  return (
    <div className="fixed inset-0 overflow-hidden z-50">
      <div className="absolute inset-0 overflow-hidden">
        <div 
          className="absolute inset-0 bg-gray-500 bg-opacity-75 transition-opacity" 
          onClick={onCancel}
        ></div>
        
        <div className="absolute inset-y-0 right-0 max-w-full flex">
          <div className="relative w-screen max-w-md">
            <div className="h-full flex flex-col bg-white shadow-xl overflow-y-auto">
              {/* Заголовок панели */}
              <div className="px-4 py-6 bg-blue-600 sm:px-6">
                <div className="flex items-center justify-between">
                  <h2 className="text-lg font-medium text-white">Создать новую задачу</h2>
                  <button
                    className="text-white hover:text-gray-200"
                    onClick={onCancel}
                  >
                    <X className="h-6 w-6" />
                  </button>
                </div>
              </div>
              
              {/* Форма создания задачи */}
              <div className="p-6">
                <form onSubmit={handleSubmit}>
                  <div className="space-y-6">
                    {/* Название задачи */}
                    <div>
                      <label htmlFor="title" className="block text-sm font-medium text-gray-700">
                        Название задачи *
                      </label>
                      <input
                        type="text"
                        id="title"
                        value={newTask.title}
                        onChange={(e) => setNewTask({...newTask, title: e.target.value})}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                        required
                      />
                    </div>
                    
                    {/* Описание */}
                    <div>
                      <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                        Описание
                      </label>
                      <textarea
                        id="description"
                        value={newTask.description}
                        onChange={(e) => setNewTask({...newTask, description: e.target.value})}
                        rows={3}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                      ></textarea>
                    </div>
                    
                    {/* Назначена */}
                    <div>
                      <label htmlFor="assigned_to" className="block text-sm font-medium text-gray-700">
                        Назначить исполнителя
                      </label>
                      <select
                        id="assigned_to"
                        value={newTask.assigned_to || ''}
                        onChange={(e) => setNewTask({...newTask, assigned_to: e.target.value ? parseInt(e.target.value) : null})}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                      >
                        <option value="">Не назначено</option>
                        {participants.map((p) => (
                          <option key={p.id} value={p.id}>{p.username}</option>
                        ))}
                      </select>
                    </div>
                    
                    {/* Родительская задача */}
                    <div>
                      <label htmlFor="parent_id" className="block text-sm font-medium text-gray-700">
                        Родительская задача
                      </label>
                      <select
                        id="parent_id"
                        value={newTask.parent_id || ''}
                        onChange={(e) => {
                          const selectedParentId = e.target.value ? parseInt(e.target.value) : null;
                          setNewTask({...newTask, parent_id: selectedParentId});
                        }}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                      >
                        <option value="">Без родительской задачи</option>
                        {tasks.map((t) => (
                          <option key={t.id} value={t.id}>{t.title} ({t.status})</option>
                        ))}
                      </select>
                      <p className="mt-1 text-xs text-gray-500">
                        Родительская задача не может быть отмечена как выполненная, 
                        пока не выполнены все её дочерние задачи
                      </p>
                    </div>
                    
                    {/* Приоритет */}
                    <div>
                      <label htmlFor="priority" className="block text-sm font-medium text-gray-700">
                        Приоритет
                      </label>
                      <select
                        id="priority"
                        value={newTask.priority}
                        onChange={(e) => setNewTask({...newTask, priority: e.target.value})}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                      >
                        <option value="low">Низкий</option>
                        <option value="medium">Средний</option>
                        <option value="high">Высокий</option>
                        <option value="critical">Критический</option>
                      </select>
                    </div>
                    
                    {/* Story Points */}
                    <div>
                      <label htmlFor="story_points" className="block text-sm font-medium text-gray-700">
                        Story Points
                      </label>
                      <input
                        type="number"
                        id="story_points"
                        value={newTask.story_points || ''}
                        onChange={(e) => setNewTask({...newTask, story_points: e.target.value ? parseInt(e.target.value) : null})}
                        min="0"
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                      />
                    </div>
                  </div>
                  
                  {/* Кнопки */}
                  <div className="mt-6 flex justify-end space-x-3">
                    <button
                      type="button"
                      onClick={onCancel}
                      className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      Отмена
                    </button>
                    <button
                      type="submit"
                      className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      Создать задачу
                    </button>
                  </div>
                </form>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AddTaskForm; 