import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';
import { expenseService } from '../services/expenseService';
import { eventService } from '../services/eventService';
import { 
  ExpenseResponse, 
  ExpensesResponse, 
  EventResponse, 
  ExpenseCreateRequest,
  UserBalance,
  BalanceReportResponse,
  EventParticipantsResponse
} from '../types/api';
import { Plus, Trash2, Edit, DollarSign, Check, CheckCircle, XCircle } from 'lucide-react';
import { toast } from 'react-toastify';

interface ExpenseFormData {
  description: string;
  amount: number;
  currency: string;
  split_method: string;
  user_ids: number[];
  isEditing: boolean;
  expense_id?: number;
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
  
  const [expenseForm, setExpenseForm] = useState<ExpenseFormData>({
    description: '',
    amount: 0,
    currency: 'RUB',
    split_method: 'equal',
    user_ids: [],
    isEditing: false
  });

  useEffect(() => {
    if (eventId) {
      fetchEventData();
    }
  }, [eventId]);

  // Функция для нормализации структуры доли расходов
  const normalizeShare = (share: any) => {
    return {
      ShareID: share.ShareID !== undefined ? share.ShareID : share.shareID,
      ExpenseID: share.ExpenseID !== undefined ? share.ExpenseID : share.expenseID,
      UserID: share.UserID !== undefined ? share.UserID : share.userID,
      Amount: share.Amount !== undefined ? share.Amount : share.amount,
      IsPaid: share.IsPaid !== undefined ? share.IsPaid : share.isPaid,
      PaidAt: share.PaidAt !== undefined ? share.PaidAt : share.paidAt
    };
  };

  const fetchEventData = async () => {
    try {
      setLoading(true);
      
      // Получаем все данные события одним запросом
      const eventData = await eventService.getEvent(eventId);
      setEvent(eventData.eventData);
      
      // Обновляем список участников события
      if (eventData.eventParticipants && eventData.eventParticipants.participants) {
        const participantsList = eventData.eventParticipants.participants.map(participant => ({
          id: participant.user.id,
          username: participant.user.username
        }));
        setParticipants(participantsList);
      }
      
      // Устанавливаем расходы из данных события
      if (eventData.expenses) {
        const normalizedExpenses = eventData.expenses.items.map(expense => ({
          ...expense,
          shares: expense.shares.map(normalizeShare)
        }));
        setExpenses(normalizedExpenses || []);
      }
      
      // Устанавливаем баланс из данных события
      if (eventData.balanceReport) {
        setBalanceReport(eventData.balanceReport);
      }
    } catch (error) {
      console.error('Ошибка при загрузке данных события:', error);
      toast.error('Не удалось загрузить данные события');
    } finally {
      setLoading(false);
    }
  };

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
          user_ids: expenseForm.user_ids
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
          user_ids: expenseForm.user_ids
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
        isEditing: false
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
      expense_id: expense.expense_id
    });
    
    // Показываем форму
    setShowAddExpenseForm(true);
  };

  const handleMarkShareAsPaid = async (shareId: number) => {
    try {
      await expenseService.markShareAsPaid(shareId, true);
      toast.success('Доля отмечена как оплаченная');
      fetchEventData();
    } catch (error) {
      console.error('Ошибка при отметке доли как оплаченной:', error);
      toast.error('Не удалось отметить долю как оплаченную');
    }
  };

  const generateSettlementRecommendations = (balances: UserBalance[]) => {
    // Создаем копию массива балансов и фильтруем, оставляя только ненулевые балансы
    // с учетом возможной погрешности вычисления
    const sortedBalances = [...balances]
      .filter(b => Math.abs(b.balance) > 0.01)
      .map(b => ({...b}));  // Создаем полную копию объектов, чтобы не изменять оригинальные данные
    
    // Массив для хранения рекомендаций
    type Recommendation = { from: string; to: string; amount: number };
    const recommendations: Recommendation[] = [];
    
    // Пока есть пользователи с ненулевым балансом
    while (sortedBalances.length > 1) {
      // Сортируем по балансу (от наименьшего к наибольшему)
      sortedBalances.sort((a, b) => a.balance - b.balance);
      
      // Берем пользователя с наименьшим (отрицательным) и наибольшим (положительным) балансом
      const debtor = sortedBalances[0]; // Должник (отрицательный баланс)
      const creditor = sortedBalances[sortedBalances.length - 1]; // Кредитор (положительный баланс)
      
      // Сумма перевода - минимальная из абсолютных значений балансов
      const transferAmount = Math.min(Math.abs(debtor.balance), creditor.balance);
      
      // Добавляем рекомендацию, только если сумма перевода значима
      if (transferAmount > 0.01) {
        recommendations.push({
          from: debtor.username,
          to: creditor.username,
          amount: transferAmount
        });
      }
      
      // Обновляем балансы
      debtor.balance += transferAmount;
      creditor.balance -= transferAmount;
      
      // Удаляем пользователей с балансом, близким к нулю
      const updatedBalances = sortedBalances.filter(b => Math.abs(b.balance) > 0.01);
      sortedBalances.length = 0;
      sortedBalances.push(...updatedBalances);
    }
    
    return recommendations;
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
              isEditing: false
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
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Участник</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-700 uppercase tracking-wider">Баланс</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-700 uppercase tracking-wider">Оплачено</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-700 uppercase tracking-wider">Не оплачено</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-700 uppercase tracking-wider">Всего к оплате</th>
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
                        <td className="px-4 py-3 whitespace-nowrap text-sm text-green-600 text-right">
                          {userBalance.paid_amount ? userBalance.paid_amount.toFixed(2) : '0.00'} RUB
                        </td>
                        <td className="px-4 py-3 whitespace-nowrap text-sm text-red-600 text-right">
                          {userBalance.unpaid_amount ? userBalance.unpaid_amount.toFixed(2) : '0.00'} RUB
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
                    <p key={index} className="text-sm bg-yellow-50 p-2 rounded-lg border border-yellow-200">
                      <span className="font-medium">{rec.from}</span> должен перевести <span className="font-bold">{rec.amount.toFixed(2)} RUB</span> пользователю <span className="font-medium">{rec.to}</span>
                    </p>
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
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Описание</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Сумма</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Оплатил(а)</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Способ деления</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Дата</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Статус платежей</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-700 uppercase tracking-wider">Действия</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {expenses.map((expense) => (
                    <tr key={expense.expense_id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{expense.description}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 font-medium">{expense.amount.toFixed(2)} {expense.currency}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                        {participants.find(p => p.id === expense.created_by)?.username || `Пользователь #${expense.created_by}`}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                        {expense.split_method === 'equal' ? 'Поровну' : 
                         expense.split_method === 'percentage' ? 'Процентами' : 
                         expense.split_method === 'amount' ? 'Суммами' : 
                         expense.split_method === 'custom' ? 'Пользовательский' : expense.split_method}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                        {new Date(expense.created_at).toLocaleString('ru-RU', {
                          day: '2-digit',
                          month: '2-digit',
                          year: 'numeric',
                          hour: '2-digit',
                          minute: '2-digit'
                        })}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500">
                        <div className="flex flex-col space-y-1">
                          {expense.shares.map(share => {
                            const username = participants.find(p => p.id === share.UserID)?.username || `#${share.UserID}`;
                            return (
                              <div key={share.ShareID} className="flex items-center justify-between">
                                <div className="flex items-center">
                                  {share.IsPaid 
                                    ? <CheckCircle className="h-5 w-5 text-green-600 mr-2" /> 
                                    : <XCircle className="h-5 w-5 text-red-600 mr-2" />}
                                  <span className={share.IsPaid ? "text-green-700" : "text-red-700"}>
                                    {username}
                                  </span>
                                </div>
                                <div className="flex items-center space-x-2">
                                  <span className={`text-xs font-medium ${share.IsPaid ? "text-green-600" : "text-red-600"}`}>
                                    {share.Amount.toFixed(2)} {expense.currency}
                                  </span>
                                  {!share.IsPaid && user && (
                                    <button 
                                      onClick={() => handleMarkShareAsPaid(share.ShareID)}
                                      className="text-blue-600 hover:text-blue-900 p-1 rounded-full hover:bg-blue-100"
                                      title="Отметить как оплачено"
                                    >
                                      <Check className="h-4 w-4" />
                                    </button>
                                  )}
                                </div>
                              </div>
                            );
                          })}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <button 
                          onClick={() => handleEditExpense(expense)}
                          className="text-blue-600 hover:text-blue-900 mr-4"
                        >
                          <Edit className="h-5 w-5" />
                        </button>
                        <button 
                          onClick={() => handleDeleteExpense(expense.expense_id)}
                          className="text-red-600 hover:text-red-900"
                        >
                          <Trash2 className="h-5 w-5" />
                        </button>
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
                            Способ деления *
                          </label>
                          <select
                            id="split_method"
                            value={expenseForm.split_method}
                            onChange={(e) => setExpenseForm({...expenseForm, split_method: e.target.value})}
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                          >
                            <option value="equal">Поровну между всеми</option>
                            <option value="percentage">По процентам</option>
                            <option value="amount">Разными суммами</option>
                            <option value="custom">Пользовательский</option>
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