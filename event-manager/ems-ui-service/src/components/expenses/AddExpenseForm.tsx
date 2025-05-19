import React, { useState } from 'react';
import { X } from 'lucide-react';
import { ExpenseCreateRequest } from '../../types/api';

interface ExpenseFormData {
  description: string;
  amount: number;
  currency: string;
  split_method: string;
  user_ids: number[];
  isEditing: boolean;
  expense_id?: number;
  task_id?: number;
}

interface AddExpenseFormProps {
  eventId: number;
  userId: number;
  participants: Array<{ id: number, username: string }>;
  onCancel: () => void;
  onSubmit: (expenseData: ExpenseCreateRequest) => Promise<void>;
  initialData?: ExpenseFormData;
  taskId?: number;
}

const AddExpenseForm: React.FC<AddExpenseFormProps> = ({
  eventId,
  userId,
  participants,
  onCancel,
  onSubmit,
  initialData,
  taskId
}) => {
  const [formData, setFormData] = useState<ExpenseFormData>(
    initialData || {
      description: '',
      amount: 0,
      currency: 'RUB',
      split_method: 'equal',
      user_ids: [],
      isEditing: false,
      task_id: taskId
    }
  );

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (formData.amount <= 0) {
      alert('Сумма должна быть больше нуля');
      return;
    }
    
    if (formData.user_ids.length === 0) {
      alert('Выберите хотя бы одного участника');
      return;
    }

    if (!userId) {
      alert('Необходимо войти в систему');
      return;
    }

    try {
      const expenseData: ExpenseCreateRequest = {
        event_id: eventId,
        created_by: userId,
        description: formData.description,
        amount: formData.amount,
        currency: formData.currency,
        split_method: formData.split_method,
        user_ids: formData.user_ids,
        task_id: formData.task_id
      };
      
      await onSubmit(expenseData);
    } catch (error) {
      console.error('Ошибка при сохранении расхода:', error);
    }
  };

  const toggleParticipant = (userId: number) => {
    setFormData(prev => {
      if (prev.user_ids.includes(userId)) {
        return {
          ...prev,
          user_ids: prev.user_ids.filter(id => id !== userId)
        };
      } else {
        return {
          ...prev,
          user_ids: [...prev.user_ids, userId]
        };
      }
    });
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
                  <h2 className="text-lg font-medium text-white">
                    {formData.isEditing ? 'Редактировать расход' : 'Создать новый расход'}
                    {taskId && ' для задачи'}
                  </h2>
                  <button
                    className="text-white hover:text-gray-200"
                    onClick={onCancel}
                  >
                    <span className="sr-only">Закрыть</span>
                    <X className="h-6 w-6" />
                  </button>
                </div>
              </div>
              
              {/* Форма */}
              <div className="p-6">
                <form onSubmit={handleSubmit}>
                  <div className="space-y-6">
                    {/* Описание */}
                    <div>
                      <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                        Описание расхода *
                      </label>
                      <input
                        type="text"
                        id="description"
                        value={formData.description}
                        onChange={(e) => setFormData({...formData, description: e.target.value})}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                        required
                      />
                    </div>
                    
                    {/* Сумма и валюта */}
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label htmlFor="amount" className="block text-sm font-medium text-gray-700">
                          Сумма *
                        </label>
                        <input
                          type="number"
                          id="amount"
                          value={formData.amount}
                          onChange={(e) => setFormData({...formData, amount: parseFloat(e.target.value)})}
                          className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                          min="0"
                          step="0.01"
                          required
                        />
                      </div>
                      <div>
                        <label htmlFor="currency" className="block text-sm font-medium text-gray-700">
                          Валюта *
                        </label>
                        <select
                          id="currency"
                          value={formData.currency}
                          onChange={(e) => setFormData({...formData, currency: e.target.value})}
                          className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                        >
                          <option value="RUB">RUB - Российский рубль</option>
                          <option value="USD">USD - Доллар США</option>
                          <option value="EUR">EUR - Евро</option>
                        </select>
                      </div>
                    </div>

                    {/* Метод разделения */}
                    <div>
                      <label htmlFor="split_method" className="block text-sm font-medium text-gray-700">
                        Метод разделения *
                      </label>
                      <select
                        id="split_method"
                        value={formData.split_method}
                        onChange={(e) => setFormData({...formData, split_method: e.target.value})}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                      >
                        <option value="equal">Поровну между всеми</option>
                        <option value="percentage">По процентам</option>
                        <option value="amount">По сумме</option>
                      </select>
                    </div>
                    
                    {/* Участники */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Выберите участников для разделения расхода *
                      </label>
                      <div className="border border-gray-300 rounded-md p-2 max-h-40 overflow-y-auto">
                        {participants.map(participant => (
                          <div 
                            key={participant.id}
                            className="flex items-center p-2 hover:bg-gray-50"
                          >
                            <input
                              type="checkbox"
                              id={`participant-${participant.id}`}
                              checked={formData.user_ids.includes(participant.id)}
                              onChange={() => toggleParticipant(participant.id)}
                              className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                            />
                            <label 
                              htmlFor={`participant-${participant.id}`}
                              className="ml-3 block text-sm text-gray-700 cursor-pointer"
                            >
                              {participant.username}
                            </label>
                          </div>
                        ))}
                        {participants.length === 0 && (
                          <p className="text-sm text-gray-500 p-2">Нет участников для выбора</p>
                        )}
                      </div>
                    </div>
                    
                    {/* Кнопки действий */}
                    <div className="flex justify-end space-x-3">
                      <button
                        type="button"
                        onClick={onCancel}
                        className="py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                      >
                        Отмена
                      </button>
                      <button
                        type="submit"
                        className="py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                      >
                        {formData.isEditing ? 'Сохранить изменения' : 'Создать расход'}
                      </button>
                    </div>
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

export default AddExpenseForm; 