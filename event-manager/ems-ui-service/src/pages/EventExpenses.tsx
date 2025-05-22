import React, { useEffect, useState, useCallback } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';
import { expenseService } from '../services/expenseService';
import { eventService } from '../services/eventService';
import { taskService } from '../services/taskService';
import { 
  ExpenseResponse, 
  ExpenseFormData, 
  ExpenseCreateRequest, 
  EventResponse, 
  BalanceReportResponse,
  EventParticipantsResponse,
  TaskResponse,
  UserBalance,
  ExpensesResponse
} from '../types/api';
import { toast } from 'react-toastify';
import { 
  Plus, 
  Trash2, 
  Edit, 
  DollarSign, 
  Check, 
  CheckCircle, 
  XCircle 
} from 'lucide-react';

interface ExpenseDetailProps {
  expense: ExpenseResponse;
  onClose: () => void;
}

const EventExpenses: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const eventId = parseInt(id || '0');
  const [loading, setLoading] = useState(true);
  const [event, setEvent] = useState<EventResponse | null>(null);
  const [expenses, setExpenses] = useState<ExpenseResponse[]>([]);
  const [participants, setParticipants] = useState<Array<{id: number, username: string}>>([]);
  const [showAddExpenseForm, setShowAddExpenseForm] = useState(false);
  const [balanceReport, setBalanceReport] = useState<BalanceReportResponse | null>(null);
  const { user } = useSelector((state: RootState) => state.auth);
  const [tasks, setTasks] = useState<TaskResponse[]>([]);
  const [error, setError] = useState<string | null>(null);
  
  const [expenseForm, setExpenseForm] = useState<ExpenseFormData>({
    description: '',
    amount: 0,
    currency: 'RUB',
    split_method: 'equal',
    user_ids: [],
    isEditing: false,
    task_id: undefined
  });

  // Нормализация долей расходов с мемоизацией для оптимизации производительности
  const normalizeShare = useCallback((share: any) => {
    if (!share) return null;
    return {
      ShareID: share.ShareID !== undefined ? share.ShareID : share.shareID,
      ExpenseID: share.ExpenseID !== undefined ? share.ExpenseID : share.expenseID,
      UserID: share.UserID !== undefined ? share.UserID : share.userID,
      Amount: share.Amount !== undefined ? share.Amount : share.amount,
      IsPaid: share.IsPaid !== undefined ? share.IsPaid : share.isPaid,
      PaidAt: share.PaidAt !== undefined ? share.PaidAt : share.paidAt
    };
  }, []);

  const fetchEventData = useCallback(async () => {
    if (!eventId) return;
    
    try {
      setLoading(true);
      setError(null);
      
      // Получаем все данные события одним запросом
      const eventData = await eventService.getEvent(eventId);
      
      // Проверка на наличие данных и их структуры
      if (!eventData) {
        throw new Error('Не удалось получить данные события');
      }
      
      // Устанавливаем основные данные события
      setEvent(eventData.eventData || null);
      
      // Обновляем список участников события
      if (eventData.eventParticipants && Array.isArray(eventData.eventParticipants.participants)) {
        const participantsList = eventData.eventParticipants.participants
          .filter(participant => participant && participant.user) // Проверка на валидность данных
          .map(participant => ({
            id: participant.user.id,
            username: participant.user.username
          }));
        setParticipants(participantsList);
      } else {
        setParticipants([]);
      }
      
      // Устанавливаем расходы из данных события
      if (eventData.expenses && Array.isArray(eventData.expenses.items)) {
        // Ограничиваем количество загружаемых расходов для предотвращения проблем с памятью
        const limitedExpenses = eventData.expenses.items.slice(0, 100); // Ограничиваем до 100 расходов
        
        const normalizedExpenses = limitedExpenses.map(expense => {
          if (!expense || !Array.isArray(expense.shares)) {
            return {
              ...expense,
              shares: []
            };
          }
          
          return {
            ...expense,
            // Нормализуем только валидные доли
            shares: expense.shares
              .filter(share => share !== null && typeof share === 'object')
              .map(normalizeShare)
              .filter(share => share !== null)
          };
        });
        
        setExpenses(normalizedExpenses || []);
      } else {
        setExpenses([]);
      }
      
      // Устанавливаем баланс из данных события
      if (eventData.balanceReport) {
        setBalanceReport(eventData.balanceReport);
      } else {
        setBalanceReport(null);
      }

      // Получаем список задач события
      if (eventData.tasks && Array.isArray(eventData.tasks.tasks)) {
        setTasks(eventData.tasks.tasks || []);
      } else {
        setTasks([]);
      }
    } catch (error) {
      console.error('Ошибка при загрузке данных события:', error);
      setError('Не удалось загрузить данные события');
      toast.error('Не удалось загрузить данные события');
    } finally {
      setLoading(false);
    }
  }, [eventId, normalizeShare]);

  useEffect(() => {
    if (eventId) {
      fetchEventData();
    }
  }, [eventId, fetchEventData]);

  // Функция для расчета рекомендаций на основе балансов пользователей
  const generateSettlementRecommendations = useCallback((balances: UserBalance[]) => {
    if (!balances || balances.length === 0) return [];
    
    // Разделяем пользователей на должников и кредиторов
    const debtors = balances
      .filter(user => user.balance < -0.01)  // Те, чей баланс отрицательный - должны деньги
      .map(user => ({
        userId: user.user_id,
        username: user.username,
        balance: user.balance,   // Отрицательное значение (долг)
        toPay: user.unpaid_amount || 0  // Сколько этот пользователь должен заплатить
      }))
      .sort((a, b) => a.balance - b.balance); // Сортировка по возрастанию, самые большие долги первыми
      
    const creditors = balances
      .filter(user => user.balance > 0.01)   // Те, чей баланс положительный - им должны деньги
      .map(user => ({
        userId: user.user_id,
        username: user.username,
        balance: user.balance,   // Положительное значение (кредит)
        toReceive: user.unpaid_amount || 0   // Сколько этот пользователь должен получить
      }))
      .sort((a, b) => b.balance - a.balance); // Сортировка по убыванию, самые большие кредиты первыми
      
    // Рекомендации по переводам
    const recommendations: {from: string; to: string; amount: number}[] = [];
    
    // Для каждого должника определяем, кому и сколько он должен перевести
    debtors.forEach(debtor => {
      let remainingDebt = Math.abs(debtor.balance); // Сколько всего должен перевести (положительное число)
      
      // Перебираем кредиторов, пока долг не будет распределен
      for (const creditor of creditors) {
        if (remainingDebt <= 0.01 || creditor.balance <= 0.01) continue;
        
        // Сколько может получить этот кредитор (не больше его баланса и не больше оставшегося долга)
        const transferAmount = Math.min(creditor.balance, remainingDebt);
        
        if (transferAmount > 0.01) {
          // Добавляем рекомендацию
          recommendations.push({
            from: debtor.username,
            to: creditor.username,
            amount: Number(transferAmount.toFixed(2)) // Округляем до 2 знаков
          });
          
          // Обновляем оставшиеся суммы
          remainingDebt -= transferAmount;
          creditor.balance -= transferAmount;
        }
      }
    });
    
    return recommendations;
  }, []);

  const handleCreateExpense = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!user || !user.id) {
      toast.error('Необходимо войти в систему');
      return;
    }
    
    if (expenseForm.amount <= 0) {
      toast.error('Сумма должна быть больше нуля');
      return;
    }
    
    if (expenseForm.user_ids.length === 0) {
      toast.error('Выберите хотя бы одного участника');
      return;
    }

    try {
      if (expenseForm.isEditing && expenseForm.expense_id) {
        // Обновляем существующий расход
        await expenseService.updateExpense(expenseForm.expense_id, {
          description: expenseForm.description,
          amount: expenseForm.amount,
          currency: expenseForm.currency,
          split_method: expenseForm.split_method,
          user_ids: expenseForm.user_ids,
          task_id: expenseForm.task_id
        });
        toast.success('Расход успешно обновлен');
      } else {
        // Создаем новый расход
        const expenseData: ExpenseCreateRequest = {
          event_id: eventId,
          created_by: user.id,
          description: expenseForm.description,
          amount: expenseForm.amount,
          currency: expenseForm.currency,
          split_method: expenseForm.split_method,
          user_ids: expenseForm.user_ids,
          task_id: expenseForm.task_id
        };
        
        await expenseService.createExpense(expenseData);
        toast.success('Расход успешно создан');
      }
      
      // Сбрасываем форму
      setExpenseForm({
        description: '',
        amount: 0,
        currency: 'RUB',
        split_method: 'equal',
        user_ids: [],
        isEditing: false,
        task_id: undefined
      });
      
      // Закрываем форму
      setShowAddExpenseForm(false);
      
      // Обновляем данные события
      fetchEventData();
    } catch (error) {
      console.error('Ошибка при создании/обновлении расхода:', error);
      toast.error('Не удалось сохранить расход');
    }
  };

  const handleDeleteExpense = async (expenseId: number) => {
    if (window.confirm('Вы уверены, что хотите удалить этот расход?')) {
      try {
        await expenseService.deleteExpense(expenseId);
        toast.success('Расход успешно удален');
        fetchEventData();
      } catch (error) {
        console.error('Ошибка при удалении расхода:', error);
        toast.error('Не удалось удалить расход');
      }
    }
  };

  const handleEditExpense = (expense: ExpenseResponse) => {
    // Заполняем форму данными существующего расхода
    setExpenseForm({
      description: expense.description,
      amount: expense.amount,
      currency: expense.currency,
      split_method: expense.split_method,
      user_ids: expense.shares.map(share => share.UserID),
      isEditing: true,
      expense_id: expense.expense_id,
      task_id: expense.task_id
    });
    
    // Показываем форму
    setShowAddExpenseForm(true);
  };

  // Функция погашения долга по рекомендации - максимально простая и точная
  const handleSettleRecommendation = async (fromUsername: string, toUsername: string, amount: number) => {
    try {
      console.log('Погашение рекомендации:', {fromUsername, toUsername, amount});
      
      // 1. Находим пользователей
      const fromUser = participants.find(p => p.username === fromUsername);
      const toUser = participants.find(p => p.username === toUsername);
      
      if (!fromUser || !toUser) {
        toast.error(`Не удалось найти пользователей: ${fromUsername} или ${toUsername}`);
        return;
      }
      
      // 2. Находим все расходы, где toUser является создателем и fromUser должен ему
      const relevantExpenses = expenses.filter(expense => 
        // Расход создан получателем денег (toUser)
        expense.created_by === toUser.id &&
        // В расходе есть неоплаченные доли от отправителя (fromUser)
        expense.shares.some(share => 
          share.UserID === fromUser.id && 
          !share.IsPaid
        )
      );
      
      if (relevantExpenses.length === 0) {
        toast.error(`У ${fromUsername} нет неоплаченных долей перед ${toUsername}`);
        return;
      }
      
      // 3. Собираем все неоплаченные доли в один список
      const unpaidShares = relevantExpenses.flatMap(expense => 
        expense.shares
          .filter(share => share.UserID === fromUser.id && !share.IsPaid)
          .map(share => ({
            expenseId: expense.expense_id,
            expenseTitle: expense.description,
            shareId: share.ShareID,
            amount: share.Amount
          }))
      );
      
      console.log('Найдены неоплаченные доли:', unpaidShares);
      
      // 4. Сортируем доли по размеру (от меньших к большим)
      unpaidShares.sort((a, b) => a.amount - b.amount);
      
      // 5. Погашаем доли до исчерпания нужной суммы
      let remainingAmount = amount;
      const settledShares = [];
      
      for (const share of unpaidShares) {
        // Если осталось погасить меньше минимальной значимой суммы, завершаем
        if (remainingAmount < 0.01) break;
        
        // Если достаточно денег для погашения этой доли
        if (remainingAmount >= share.amount) {
          try {
            // Отмечаем долю как оплаченную
            await expenseService.markShareAsPaid(share.shareId, true);
            
            // Вычитаем сумму из оставшейся
            remainingAmount -= share.amount;
            
            // Сохраняем информацию о погашенной доле
            settledShares.push({
              ...share,
              status: 'погашена'
            });
            
          } catch (error) {
            console.error(`Ошибка при погашении доли ${share.shareId}:`, error);
            settledShares.push({
              ...share,
              status: 'ошибка'
            });
          }
        }
      }
      
      // 6. Выводим результаты
      if (settledShares.length > 0) {
        const totalPaid = settledShares.reduce((sum, share) => 
          share.status === 'погашена' ? sum + share.amount : sum, 0);
        
        toast.success(
          `Погашено долгов на сумму ${totalPaid.toFixed(2)} RUB (${settledShares.length} долей)`, 
          { autoClose: 5000 }
        );
        
        console.log('Детали погашения:', settledShares);
        
        // Обновляем данные события для обновления интерфейса
        fetchEventData();
      } else {
        toast.warning('Не удалось погасить ни одной доли');
      }
      
    } catch (error) {
      console.error('Ошибка при погашении долгов:', error);
      toast.error('Произошла ошибка при погашении долга');
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-6xl mx-auto bg-red-50 p-4 rounded-lg mt-4">
        <h1 className="text-lg font-medium text-red-700">Ошибка при загрузке данных</h1>
        <p className="text-red-600">{error}</p>
        <button 
          onClick={() => fetchEventData()}
          className="mt-2 bg-red-100 text-red-800 px-4 py-2 rounded hover:bg-red-200"
        >
          Попробовать снова
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto relative px-4">
      {/* Заголовок и кнопка добавления */}
      <div className="flex justify-between items-center mb-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">
            {event ? event.title : 'Загрузка...'} - Расходы
          </h1>
          <p className="text-sm text-gray-500">
            Управление расходами и расчет долгов
          </p>
        </div>
        <button
          onClick={() => {
            setExpenseForm({
              description: '',
              amount: 0,
              currency: 'RUB',
              split_method: 'equal',
              user_ids: [],
              isEditing: false,
              task_id: undefined
            });
            setShowAddExpenseForm(true);
          }}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center hover:bg-blue-700 transition"
        >
          <Plus className="h-5 w-5 mr-1" />
          Добавить расход
        </button>
      </div>

      {/* Основное содержимое */}
      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        {/* Шапка */}
        <div className="border-b border-gray-200 bg-gray-50 p-4">
          <div className="flex justify-between items-center">
            <h2 className="text-lg font-medium text-gray-900">
              Список расходов ({expenses.length})
            </h2>
            <Link to={`/events/${eventId}`} className="text-blue-600 hover:text-blue-800 text-sm">
              Вернуться к событию
            </Link>
          </div>
        </div>
        
        {/* Баланс участников */}
        {balanceReport && (
          <div className="p-6 border-b">
            <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
              <DollarSign className="h-5 w-5 mr-2 text-blue-600" />
              Отчет о балансе
            </h3>

            <div className="bg-blue-50 p-4 rounded-lg mb-4">
              <p className="text-sm text-blue-800">
                <span className="font-medium">Общая сумма расходов:</span> {balanceReport.total_amount.toFixed(2)} RUB
              </p>
              <p className="text-xs text-blue-600 mt-1">
                Этот отчет показывает, кто должен получить деньги (положительный баланс), а кто должен вернуть (отрицательный баланс).
              </p>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full divide-y divide-gray-200 table-fixed">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider w-1/3">Участник</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-700 uppercase tracking-wider w-1/3">Баланс</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-700 uppercase tracking-wider w-1/3">Всего к оплате</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {balanceReport.user_balances
                    .sort((a, b) => b.balance - a.balance)
                    .map(userBalance => (
                      <tr key={`balance-${userBalance.user_id}`} className="hover:bg-gray-50">
                        <td className="px-4 py-3 whitespace-nowrap text-sm font-medium">{userBalance.username}</td>
                        <td className={`px-4 py-3 whitespace-nowrap text-sm font-medium text-right ${
                          userBalance.balance > 0.01 ? 'text-green-600' : 
                          userBalance.balance < -0.01 ? 'text-red-600' : 'text-gray-600'
                        }`}>
                          {userBalance.balance.toFixed(2)} RUB
                        </td>
                        <td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-gray-900 text-right">
                          {userBalance.total_due ? userBalance.total_due.toFixed(2) : '0.00'} RUB
                        </td>
                      </tr>
                    ))}
                </tbody>
              </table>
            </div>

            {/* Рекомендации по распределению */}
            {balanceReport.user_balances.some(ub => Math.abs(ub.balance) > 0.01) && (
              <div className="mt-6 border-t border-gray-200 pt-4">
                <h4 className="font-medium text-gray-900 mb-2">Рекомендации по расчетам:</h4>
                <div className="space-y-2">
                  {generateSettlementRecommendations(balanceReport.user_balances).map((rec, index) => (
                    <div key={index} className="flex items-center justify-between bg-yellow-50 p-3 rounded-lg border border-yellow-200">
                      <p className="text-sm flex-grow">
                        <span className="font-medium">{rec.from}</span> должен перевести <span className="font-bold">{rec.amount.toFixed(2)} RUB</span> пользователю <span className="font-medium">{rec.to}</span>
                      </p>
                      {user && (
                        <button 
                          onClick={() => {
                            if (window.confirm(`Вы действительно хотите отметить оплату ${rec.amount.toFixed(2)} RUB от ${rec.from} пользователю ${rec.to}?`)) {
                              handleSettleRecommendation(rec.from, rec.to, rec.amount);
                            }
                          }}
                          className="ml-4 px-3 py-1.5 bg-gradient-to-r from-green-500 to-green-600 text-white text-sm font-medium rounded-md hover:from-green-600 hover:to-green-700 transition-all shadow-md hover:shadow-lg flex items-center whitespace-nowrap"
                          title="Отметить как оплачено"
                        >
                          <DollarSign className="h-4 w-4 mr-1" /> Погасить долг
                        </button>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
        
        {/* Список расходов */}
        <div className="p-4">
          {expenses.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider" style={{width: '20%'}}>Описание</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider" style={{width: '10%'}}>Сумма</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider" style={{width: '15%'}}>Оплатил(а)</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider" style={{width: '15%'}}>Способ деления</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider" style={{width: '12%'}}>Дата</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider" style={{width: '18%'}}>Неоплаченные доли</th>
                    <th className="px-4 py-3 text-center text-xs font-medium text-gray-700 uppercase tracking-wider" style={{width: '5%'}}>Задача</th>
                    <th className="px-4 py-3 text-center text-xs font-medium text-gray-700 uppercase tracking-wider" style={{width: '5%'}}>Действия</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {expenses.map((expense) => (
                    <tr key={expense.expense_id} className="hover:bg-gray-50">
                      <td className="px-4 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{expense.description}</td>
                      <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-600 font-medium">{expense.amount.toFixed(2)} {expense.currency}</td>
                      <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-600">
                        {participants.find(p => p.id === expense.created_by)?.username || `Пользователь #${expense.created_by}`}
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-600">
                        {expense.split_method === 'equal' ? 'Поровну' : 
                         expense.split_method === 'percentage' ? 'Процентами' : 
                         expense.split_method === 'amount' ? 'Суммами' : 
                         expense.split_method === 'custom' ? 'Пользовательский' : expense.split_method}
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-600">
                        {new Date(expense.created_at).toLocaleString('ru-RU', {
                          day: '2-digit',
                          month: '2-digit',
                          year: 'numeric',
                          hour: '2-digit',
                          minute: '2-digit'
                        })}
                      </td>
                      <td className="px-4 py-4 text-sm text-red-600">
                        <div className="flex flex-col space-y-1">
                          {expense.shares
                            // Только неоплаченные доли и не доли создателя
                            .filter(share => !share.IsPaid && share.UserID !== expense.created_by)
                            .map(share => {
                              const username = participants.find(p => p.id === share.UserID)?.username || `#${share.UserID}`;
                              return (
                                <div key={share.ShareID} className="flex items-center justify-between">
                                  <span className="text-red-700 mr-2">{username}</span>
                                  <span className="text-xs font-medium text-red-600">
                                    {share.Amount.toFixed(2)} {expense.currency}
                                  </span>
                                </div>
                              );
                            })}
                        </div>
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap text-center text-sm font-medium">
                        {expense.task_id && (
                          <Link to={`/events/${eventId}/tasks`} className="text-blue-600 hover:text-blue-900">
                            #{expense.task_id}
                          </Link>
                        )}
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap text-center text-sm font-medium">
                        <div className="flex justify-center space-x-2">
                          <button 
                            onClick={() => handleEditExpense(expense)}
                            className="text-blue-600 hover:text-blue-900"
                          >
                            <Edit className="h-5 w-5" />
                          </button>
                          <button 
                            onClick={() => handleDeleteExpense(expense.expense_id)}
                            className="text-red-600 hover:text-red-900"
                          >
                            <Trash2 className="h-5 w-5" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">
              <DollarSign className="h-12 w-12 mx-auto text-gray-300 mb-2" />
              <p>У этого события пока нет расходов.</p>
              <p className="text-sm mt-1">Нажмите "Добавить расход", чтобы создать первый расход.</p>
            </div>
          )}
        </div>
      </div>

      {/* Форма добавления/редактирования расхода */}
      {showAddExpenseForm && (
        <div className="fixed inset-0 overflow-hidden z-50">
          <div className="absolute inset-0 overflow-hidden">
            <div 
              className="absolute inset-0 bg-gray-500 bg-opacity-75 transition-opacity" 
              onClick={() => setShowAddExpenseForm(false)}
            ></div>
            
            <div className="absolute inset-y-0 right-0 max-w-full flex">
              <div className="relative w-screen max-w-md">
                <div className="h-full flex flex-col bg-white shadow-xl overflow-y-auto">
                  {/* Заголовок панели */}
                  <div className="px-4 py-6 bg-blue-600 sm:px-6">
                    <div className="flex items-center justify-between">
                      <h2 className="text-lg font-medium text-white">
                        {expenseForm.isEditing ? 'Редактировать расход' : 'Создать новый расход'}
                      </h2>
                      <button
                        className="text-white hover:text-gray-200"
                        onClick={() => setShowAddExpenseForm(false)}
                      >
                        <span className="sr-only">Close</span>
                        <Trash2 className="h-6 w-6" />
                      </button>
                    </div>
                  </div>
                  
                  {/* Форма */}
                  <div className="p-6">
                    <form onSubmit={handleCreateExpense}>
                      <div className="space-y-6">
                        {/* Описание */}
                        <div>
                          <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                            Описание расхода *
                          </label>
                          <input
                            type="text"
                            id="description"
                            value={expenseForm.description}
                            onChange={(e) => setExpenseForm({...expenseForm, description: e.target.value})}
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
                              value={expenseForm.amount}
                              onChange={(e) => setExpenseForm({...expenseForm, amount: parseFloat(e.target.value)})}
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
                              value={expenseForm.currency}
                              onChange={(e) => setExpenseForm({...expenseForm, currency: e.target.value})}
                              className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                            >
                              <option value="RUB">RUB - Российский рубль</option>
                              <option value="USD">USD - Доллар США</option>
                              <option value="EUR">EUR - Евро</option>
                            </select>
                          </div>
                        </div>
                        
                        {/* Способ деления */}
                        <div>
                          <label htmlFor="split_method" className="block text-sm font-medium text-gray-700">
                            Метод разделения *
                          </label>
                          <select
                            id="split_method"
                            value={expenseForm.split_method}
                            onChange={(e) => setExpenseForm({...expenseForm, split_method: e.target.value})}
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                          >
                            <option value="equal">Поровну между всеми</option>
                            <option value="percentage">По процентам</option>
                            <option value="amount">По сумме</option>
                          </select>
                        </div>

                        {/* Привязка к задаче */}
                        <div>
                          <label htmlFor="task_id" className="block text-sm font-medium text-gray-700">
                            Привязать к задаче (опционально)
                          </label>
                          <select
                            id="task_id"
                            value={expenseForm.task_id || ''}
                            onChange={(e) => {
                              const value = e.target.value;
                              setExpenseForm({
                                ...expenseForm, 
                                task_id: value ? parseInt(value) : undefined
                              });
                            }}
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                          >
                            <option value="">Не привязывать к задаче</option>
                            {tasks.map(task => (
                              <option key={task.id} value={task.id}>
                                {task.title} (#{task.id})
                              </option>
                            ))}
                          </select>
                        </div>
                        
                        {/* Участники */}
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-2">
                            Участники *
                          </label>
                          <div className="space-y-2 max-h-48 overflow-y-auto">
                            {participants.map((participant) => (
                              <div key={participant.id} className="flex items-center">
                                <input
                                  type="checkbox"
                                  id={`participant-${participant.id}`}
                                  checked={expenseForm.user_ids.includes(participant.id)}
                                  onChange={(e) => {
                                    if (e.target.checked) {
                                      setExpenseForm({
                                        ...expenseForm,
                                        user_ids: [...expenseForm.user_ids, participant.id]
                                      });
                                    } else {
                                      setExpenseForm({
                                        ...expenseForm,
                                        user_ids: expenseForm.user_ids.filter(id => id !== participant.id)
                                      });
                                    }
                                  }}
                                  className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                                />
                                <label htmlFor={`participant-${participant.id}`} className="ml-2 block text-sm text-gray-900">
                                  {participant.username}
                                </label>
                              </div>
                            ))}
                          </div>
                          {participants.length === 0 && (
                            <p className="text-sm text-red-500 mt-1">
                              У события нет участников. Пожалуйста, добавьте участников в событие.
                            </p>
                          )}
                        </div>
                      </div>
                      
                      {/* Кнопки */}
                      <div className="mt-6 flex justify-end space-x-3">
                        <button
                          type="button"
                          onClick={() => setShowAddExpenseForm(false)}
                          className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        >
                          Отмена
                        </button>
                        <button
                          type="submit"
                          className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        >
                          {expenseForm.isEditing ? 'Сохранить' : 'Создать расход'}
                        </button>
                      </div>
                    </form>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default EventExpenses; 