import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';
import { taskService } from '../services/taskService';
import { eventService } from '../services/eventService';
import { expenseService } from '../services/expenseService';
import { TaskResponse, TaskStatus, EventResponse, TaskCreateRequest, TaskUpdateRequest, ExpenseCreateRequest, ExpenseResponse } from '../types/api';
import { Plus } from 'lucide-react';
import { toast } from 'react-toastify';

// Импортируем компоненты
import TaskList from '../components/tasks/TaskList';
import AddTaskForm from '../components/tasks/AddTaskForm';
import AddExpenseForm from '../components/expenses/AddExpenseForm';
import TaskStyles from '../components/tasks/TaskStyles';
import { buildHierarchicalTasks, canChangeTaskStatus } from '../components/tasks/taskUtils';
import { HierarchicalTask, TaskExpense } from '../components/tasks/types';

const EventTasks: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const eventId = parseInt(id || '0');
  const [loading, setLoading] = useState(true);
  const [event, setEvent] = useState<EventResponse | null>(null);
  const [tasks, setTasks] = useState<TaskResponse[]>([]);
  const [hierarchicalTasks, setHierarchicalTasks] = useState<HierarchicalTask[]>([]);
  const [usernames, setUsernames] = useState<Map<number, string>>(new Map());
  const [showAddTaskDrawer, setShowAddTaskDrawer] = useState(false);
  const [showAddExpenseDrawer, setShowAddExpenseDrawer] = useState(false);
  const [selectedTaskId, setSelectedTaskId] = useState<number | undefined>(undefined);
  const { user } = useSelector((state: RootState) => state.auth);
  const [participants, setParticipants] = useState<Array<{id: number, username: string}>>([]);
  const [expandedTaskIds, setExpandedTaskIds] = useState<Set<number>>(new Set());
  const [taskForm, setTaskForm] = useState<{
    title: string;
    description: string;
    assigned_to?: number;
    parent_id?: number;
    priority: string;
    story_points: number;
    task_id?: number;
  }>({
    title: '',
    description: '',
    assigned_to: user?.id || 0,
    parent_id: undefined,
    priority: 'medium',
    story_points: 1
  });
  const [showAddTaskForm, setShowAddTaskForm] = useState(false);
  const [expenses, setExpenses] = useState<ExpenseResponse[]>([]);

  useEffect(() => {
    if (eventId) {
      fetchEventAndTasks();
    }
  }, [eventId]);

  // Функция для форматирования списка участников
  const updateParticipantsList = (eventData: any) => {
    if (eventData.eventParticipants && eventData.eventParticipants.participants) {
      const participantsList = eventData.eventParticipants.participants.map((participant: any) => ({
        id: participant.user.id,
        username: participant.user.username
      }));
      setParticipants(participantsList);
    }
  };

  // Объединенная функция fetchEventAndTasks, включающая загрузку участников
  const fetchEventAndTasks = async () => {
    try {
      setLoading(true);
      
      // Fetch event details first
      const eventData = await eventService.getEvent(eventId);
      setEvent(eventData.eventData);
      
      // Загружаем расходы для события
      let eventExpenses: ExpenseResponse[] = [];
      if (eventData.expenses && eventData.expenses.items) {
        eventExpenses = eventData.expenses.items;
        setExpenses(eventExpenses);
      }
      
      // Extract tasks from event data
      const taskData = eventData.tasks;
      if (taskData && taskData.tasks) {
        // Нормализуем структуру задач и обрабатываем различные форматы parent_id
        const normalizedTasks = taskData.tasks.map(task => {
          // Создаем копию задачи для работы с ней
          const taskCopy = {...task};
          
          // Обрабатываем все возможные варианты написания parent_id
          let parentId = null;
          
          // Проверяем все возможные варианты написания parent_id
          if ('parent_id' in taskCopy && taskCopy.parent_id !== undefined) {
            parentId = taskCopy.parent_id;
          } else if ('ParentId' in taskCopy && (taskCopy as any).ParentId !== undefined) {
            parentId = (taskCopy as any).ParentId;
          } else if ('parentId' in taskCopy && (taskCopy as any).parentId !== undefined) {
            parentId = (taskCopy as any).parentId;
          } else if ('parent-id' in taskCopy && (taskCopy as any)['parent-id'] !== undefined) {
            parentId = (taskCopy as any)['parent-id'];
          }
          
          // Преобразуем строковое значение parent_id в число, если необходимо
          if (parentId !== null && typeof parentId === 'string') {
            const numericId = parseInt(parentId, 10);
            if (!isNaN(numericId)) {
              parentId = numericId;
            }
          }
          
          // Если parent_id равен 0, считаем его как null
          if (parentId === 0) {
            parentId = null;
          }
          
          // Устанавливаем нормализованное значение parent_id
          taskCopy.parent_id = parentId;
          
          return taskCopy;
        });
        
        // Sort tasks by creation date in descending order
        const sortedTasks = [...normalizedTasks].sort((a, b) => 
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        );
        setTasks(sortedTasks);
        
        // Добавляем расходы к задачам и расчитываем общие суммы
        const tasksWithExpenses = addExpensesToTasks(sortedTasks, eventExpenses);
        
        // Build hierarchical task structure
        const hierarchicalTasks = buildHierarchicalTasks(tasksWithExpenses);
        
        // Расчитываем общие суммы расходов для иерархии задач
        const tasksWithTotalExpenses = calculateTotalExpenses(hierarchicalTasks);
        setHierarchicalTasks(tasksWithTotalExpenses);
      }

      // Build username map from event participants
      const usernameMap = new Map<number, string>();
      if (eventData.eventParticipants && eventData.eventParticipants.participants) {
        eventData.eventParticipants.participants.forEach(participant => {
          usernameMap.set(participant.user.id, participant.user.username);
        });
        // Обновить список участников для формы создания задачи
        updateParticipantsList(eventData);
      }
      setUsernames(usernameMap);
      
    } catch (error) {
      console.error('Error fetching event and tasks:', error);
      toast.error('Failed to load event tasks');
    } finally {
      setLoading(false);
    }
  };

  // Функция для добавления расходов к соответствующим задачам
  const addExpensesToTasks = (tasks: TaskResponse[], expenses: ExpenseResponse[]): TaskResponse[] => {
    // Дебаг вывод
    console.log(`Найдено задач: ${tasks.length}, расходов: ${expenses.length}`);
    console.log(`Расходы с привязкой к задачам:`, expenses.filter(e => e.task_id).length);

    // Создаем карту для быстрого доступа к задачам
    const taskMap = new Map<number, TaskResponse>();
    tasks.forEach(task => {
      taskMap.set(task.id, { ...task } as any);
    });
    
    // Распределяем расходы по задачам
    expenses.forEach(expense => {
      if (expense.task_id && taskMap.has(expense.task_id)) {
        const task = taskMap.get(expense.task_id)!;
        if (!('expenses' in task)) {
          (task as any).expenses = [];
        }
        (task as any).expenses.push({
          expense_id: expense.expense_id,
          description: expense.description,
          amount: expense.amount,
          currency: expense.currency
        });
        console.log(`Задача #${task.id}: добавлен расход ${expense.expense_id} - ${expense.amount}`);
      }
    });
    
    // Возвращаем обновленные задачи
    return Array.from(taskMap.values());
  };

  // Функция для расчета общей суммы расходов для иерархии задач
  const calculateTotalExpenses = (tasks: HierarchicalTask[]): HierarchicalTask[] => {
    // Создаем функцию для рекурсивного расчета суммы расходов для одной задачи
    const calculateTaskExpenses = (task: HierarchicalTask): HierarchicalTask => {
      const taskCopy = { ...task };
      
      // Считаем сумму собственных расходов задачи
      const ownExpensesTotal = (task.expenses || []).reduce((sum, exp) => 
        sum + exp.amount, 0);
      
      // Если есть дочерние задачи, обрабатываем их
      let childrenTotal = 0;
      if (task.children && task.children.length > 0) {
        // Рекурсивно вычисляем для всех дочерних задач
        taskCopy.children = task.children.map(calculateTaskExpenses);
        
        // Суммируем общие суммы дочерних задач
        childrenTotal = taskCopy.children.reduce((sum, child) => 
          sum + (child.totalExpenses || 0), 0);
      }
      
      // Устанавливаем общую сумму расходов только если она больше 0
      const totalAmount = ownExpensesTotal + childrenTotal;
      if (totalAmount > 0.01) {
        taskCopy.totalExpenses = totalAmount;
        
        console.log(`Задача #${task.id} (${task.title}): собственные расходы ${ownExpensesTotal}, дочерние ${childrenTotal}, всего ${totalAmount}`);
      }
      
      return taskCopy;
    };
    
    // Применяем эту функцию ко всем корневым задачам
    return tasks.map(calculateTaskExpenses);
  };

  const handleTaskStatusChange = async (taskId: number, newStatus: TaskStatus) => {
    try {
      const task = tasks.find(t => t.id === taskId);
      if (!task) return;

      // Проверяем статус родительской задачи
      if (task.parent_id) {
        const parentTask = tasks.find(t => t.id === task.parent_id);
        if (parentTask && parentTask.status === TaskStatus.Completed) {
          toast.error('Нельзя изменить статус подзадачи, если родительская задача завершена');
          return;
        }
      }

      // Проверяем статус подзадач
      const subtasks = tasks.filter(t => t.parent_id === taskId);
      if (newStatus === TaskStatus.Completed && subtasks.some(t => t.status !== TaskStatus.Completed)) {
        toast.error('Нельзя завершить задачу, пока не завершены все подзадачи');
        return;
      }

      await taskService.updateTask(taskId, { status: newStatus });
      
      // Обновляем данные, но сохраняем состояние развернутых задач
      await fetchEventAndTasks();
      
      // Если задача завершена, автоматически раскрываем родительскую задачу
      if (newStatus === TaskStatus.Completed && task.parent_id) {
        setExpandedTaskIds(prev => {
          const newSet = new Set(prev);
          newSet.add(task.parent_id!);
          return newSet;
        });
      }
      
      toast.success('Статус задачи обновлен');
    } catch (error) {
      console.error('Ошибка при обновлении статуса задачи:', error);
      toast.error('Не удалось обновить статус задачи');
    }
  };

  const handleTaskDelete = async (taskId: number) => {
    if (!window.confirm('Вы уверены, что хотите удалить эту задачу?')) return;

    try {
      // Сохраняем информацию о задаче перед удалением
      const task = tasks.find(t => t.id === taskId);
      
      await taskService.deleteTask(taskId);
      await fetchEventAndTasks();
      
      // Если у удалённой задачи был родитель, автоматически раскрываем родительскую задачу
      if (task && task.parent_id) {
        setExpandedTaskIds(prev => {
          const newSet = new Set(prev);
          newSet.add(task.parent_id!);
          return newSet;
        });
      }
      
      toast.success('Задача успешно удалена');
    } catch (error) {
      console.error('Ошибка при удалении задачи:', error);
      toast.error('Не удалось удалить задачу');
    }
  };

  const handleTaskCreate = async (taskData: TaskCreateRequest) => {
    try {
      const response = await taskService.createTask(taskData);
      
      setTaskForm({
        title: '',
        description: '',
        assigned_to: user?.id || 0,
        parent_id: undefined,
        priority: 'medium',
        story_points: 1
      });
      // Закрываем форму создания задачи
      setShowAddTaskDrawer(false);
      
      await fetchEventAndTasks();
      
      // Если у новой задачи есть родитель, автоматически раскрываем родительскую задачу
      if (taskData.parent_id) {
        setExpandedTaskIds(prev => {
          const newSet = new Set(prev);
          newSet.add(taskData.parent_id!);
          return newSet;
        });
      }
      
      toast.success('Задача успешно создана');
    } catch (error) {
      console.error('Ошибка при создании задачи:', error);
      toast.error('Не удалось создать задачу');
    }
  };

  const handleTaskUpdate = async (taskId: number, taskData: TaskUpdateRequest) => {
    try {
      await taskService.updateTask(taskId, taskData);
      
      setTaskForm({
        title: '',
        description: '',
        assigned_to: user?.id || 0,
        parent_id: undefined,
        priority: 'medium',
        story_points: 1
      });
      // Закрываем форму редактирования задачи
      setShowAddTaskDrawer(false);
      
      await fetchEventAndTasks();
      
      // Если у задачи есть родитель, автоматически раскрываем родительскую задачу
      if (taskData.parent_id) {
        setExpandedTaskIds(prev => {
          const newSet = new Set(prev);
          newSet.add(taskData.parent_id!);
          return newSet;
        });
      }
      
      toast.success('Задача успешно обновлена');
    } catch (error) {
      console.error('Ошибка при обновлении задачи:', error);
      toast.error('Не удалось обновить задачу');
    }
  };

  const handleTaskFormSubmit = async (taskData: TaskCreateRequest) => {
    if (!user || !user.id) {
      toast.error('Необходимо войти в систему');
      return;
    }

    if (taskForm.task_id) {
      await handleTaskUpdate(taskForm.task_id, taskData);
    } else {
      await handleTaskCreate(taskData);
    }
  };

  // Обработчик добавления расхода к задаче
  const handleAddExpense = (taskId: number) => {
    setSelectedTaskId(taskId);
    setShowAddExpenseDrawer(true);
  };

  // Обработчик создания расхода
  const handleCreateExpense = async (expenseData: ExpenseCreateRequest) => {
    try {
      await expenseService.createExpense(expenseData);
      setShowAddExpenseDrawer(false);
      toast.success('Расход успешно создан');
      
      // Обновляем данные после создания расхода для отображения изменений
      await fetchEventAndTasks();
    } catch (error) {
      console.error('Ошибка при создании расхода:', error);
      toast.error('Не удалось создать расход');
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto relative">
      {/* Добавляем CSS-анимацию для плавного раскрытия */}
      <TaskStyles />
      
      {/* Кнопка создания задачи */}
      <div className="flex justify-between items-center mb-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">
            {event ? event.title : 'Loading...'} - Задачи
          </h1>
          <p className="text-sm text-gray-500">
            Управление задачами и отслеживание прогресса
          </p>
        </div>
        <button
          onClick={() => setShowAddTaskDrawer(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center hover:bg-blue-700 transition"
        >
          <Plus className="h-5 w-5 mr-1" />
          Добавить задачу
        </button>
      </div>

      {/* Основное содержимое */}
      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        {loading ? (
          <div className="p-8 text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
            <p className="mt-4 text-gray-600">Загрузка задач...</p>
          </div>
        ) : (
          <div>
            {/* Шапка списка задач */}
            <div className="border-b border-gray-200 bg-gray-50 p-4">
              <div className="flex justify-between items-center">
                <h2 className="text-lg font-medium text-gray-900">
                  Иерархия задач ({tasks.length})
                </h2>
                <Link to={`/events/${eventId}`} className="text-blue-600 hover:text-blue-800 text-sm">
                  Вернуться к событию
                </Link>
              </div>
            </div>
            
            {/* Список задач */}
            <div className="p-4">
              <TaskList 
                hierarchicalTasks={hierarchicalTasks}
                tasks={tasks}
                usernames={usernames}
                onUpdateStatus={handleTaskStatusChange}
                onDelete={handleTaskDelete}
                expandedTaskIds={expandedTaskIds}
                setExpandedTaskIds={setExpandedTaskIds}
                onAddExpense={handleAddExpense}
              />
            </div>
          </div>
        )}
      </div>

      {/* Форма добавления задачи */}
      {showAddTaskDrawer && (
        <AddTaskForm 
          eventId={eventId}
          tasks={tasks}
          participants={participants}
          onCancel={() => setShowAddTaskDrawer(false)}
          onSubmit={handleTaskFormSubmit}
        />
      )}

      {/* Форма добавления расхода */}
      {showAddExpenseDrawer && (
        <AddExpenseForm 
          eventId={eventId}
          userId={user?.id || 0}
          participants={participants}
          onCancel={() => setShowAddExpenseDrawer(false)}
          onSubmit={handleCreateExpense}
          taskId={selectedTaskId}
        />
      )}
    </div>
  );
};

export default EventTasks;
