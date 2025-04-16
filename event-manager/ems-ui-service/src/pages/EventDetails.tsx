import React, { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';
import { eventService } from '../services/eventService';
import { taskService } from '../services/taskService';
import { commentService } from '../services/commentService';
import { EventData, TaskCreateRequest, CommentCreateRequest, Comment, TaskStatus, EventParticipant, EventParticipantCreateRequest } from '../types/api';
import { toast } from 'react-toastify';

const EventDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const eventId = parseInt(id || '0');
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [eventData, setEventData] = useState<EventData | null>(null);
  const [newCommentContent, setNewCommentContent] = useState('');
  const currentUser = useSelector((state: RootState) => state.auth.user);
  
  // Состояние для модального окна добавления участника
  const [showAddParticipantModal, setShowAddParticipantModal] = useState(false);
  const [newParticipantData, setNewParticipantData] = useState({
    username: '',
    email: ''
  });

  // State to store user mappings
  const [userMap, setUserMap] = useState<Map<number, string>>(new Map());

  // Fetch event data
  useEffect(() => {
    const fetchEventData = async () => {
      console.log('Attempting to fetch event with ID:', eventId);
      if (!eventId) {
        console.error('No event ID provided');
        setLoading(false);
        toast.error('Invalid event ID');
        navigate('/');
        return;
      }
      try {
        setLoading(true);
        // Add a slight delay to ensure loading state is visible for debugging
        await new Promise(resolve => setTimeout(resolve, 500));
        
        const data = await eventService.getEvent(eventId);
        console.log('Received event data:', data);
        
        if (data && data.eventData) {
          console.log('Setting event data, event found:', data.eventData.title);
          setEventData(data);
          
          // Create a map of user IDs to usernames from participants
          const usernameMap = new Map<number, string>();
          if (data.eventParticipants && data.eventParticipants.participants) {
            data.eventParticipants.participants.forEach(participant => {
              usernameMap.set(participant.user.id, participant.user.username);
            });
          }
          setUserMap(usernameMap);
        } else {
          console.error('Event data is null or missing eventData property');
          toast.error('Event not found or data is invalid');
          navigate('/');
        }
      } catch (error) {
        console.error('Error fetching event details:', error);
        toast.error('Failed to load event details');
        navigate('/');
      } finally {
        setLoading(false);
      }
    };

    fetchEventData();
  }, [eventId, navigate]);

  // Обработчик удаления участника
  const handleRemoveParticipant = async (participantId: number) => {
    try {
      await eventService.removeParticipant(participantId);
      toast.success('Participant removed');
      // Обновляем данные события
      const data = await eventService.getEvent(eventId);
      setEventData(data);
    } catch (error) {
      console.error('Error removing participant:', error);
      toast.error('Failed to remove participant');
    }
  };

  // Обработчик добавления участника
  const handleAddParticipant = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newParticipantData.username.trim()) {
      toast.error('Username is required');
      return;
    }
    
    try {
      // Предполагаем, что у нас есть API для получения ID пользователя по имени пользователя
      // В реальном приложении это может быть встроено в API добавления участника
      // или отдельный эндпоинт для поиска пользователя
      
      // Для этого примера, предположим, что у нас есть доступ к ID пользователя
      // Например, через поле формы или другой механизм
      
      // Создаем запрос на добавление участника
      const participantData: EventParticipantCreateRequest = {
        user_id: parseInt(newParticipantData.username), // В реальном приложении здесь будет ID пользователя
        event_title: eventData?.eventData.title || ''
      };
      
      await eventService.addParticipant(eventId, participantData);
      toast.success('Participant added successfully');
      
      // Сбрасываем форму
      setNewParticipantData({
        username: '',
        email: ''
      });
      
      // Закрываем модальное окно
      setShowAddParticipantModal(false);
      
      // Обновляем данные события
      const data = await eventService.getEvent(eventId);
      setEventData(data);
    } catch (error) {
      console.error('Error adding participant:', error);
      toast.error('Failed to add participant');
    }
  };

  // Handle creating a new comment
  const handleCreateComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCommentContent.trim() || !currentUser) return;

    const commentData: CommentCreateRequest = {
      content: newCommentContent,
      event_id: eventId,
      sender_id: currentUser.id
    };

    try {
      const newComment = await commentService.createComment(commentData);
      toast.success('Comment added successfully');
      setNewCommentContent('');
      
      // Создаем новый объект комментария и добавляем его к текущим комментариям
      if (eventData) {
        const newCommentObject: Comment = {
          commentId: newComment.commentId || Date.now(), // Используем ID из ответа или временный ID
          content: commentData.content,
          eventId: commentData.event_id,
          senderId: commentData.sender_id,
          createdAt: new Date().toISOString(),
          isRead: false,
          isDeleted: false,
          taskId: undefined
        };
        
        // Создаем новый объект данных события с обновленными комментариями
        const updatedEventData = {
          ...eventData,
          comments: {
            ...eventData.comments,
            comments: [newCommentObject, ...eventData.comments.comments]
          }
        };
        
        // Обновляем состояние
        setEventData(updatedEventData);
        
        // Затем обновляем данные с сервера асинхронно
        const data = await eventService.getEvent(eventId);
        setEventData(data);
      }
    } catch (error) {
      console.error('Error adding comment:', error);
      toast.error('Failed to add comment');
    }
  };

  // Handle updating task status
  const handleUpdateTaskStatus = async (taskId: number, status: TaskStatus) => {
    try {
      await taskService.updateTask(taskId, { status });
      toast.success('Task status updated');
      // Refresh event data
      const data = await eventService.getEvent(eventId);
      setEventData(data);
    } catch (error) {
      console.error('Error updating task status:', error);
      toast.error('Failed to update task status');
    }
  };

  // Handle deleting the event
  const handleDeleteEvent = async () => {
    try {
      await eventService.deleteEvent(eventId);
      toast.success('Event deleted');
      navigate('/');
    } catch (error) {
      console.error('Error deleting event:', error);
      toast.error('Failed to delete event');
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    );
  }

  if (!eventData) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[50vh]">
        <p className="text-gray-600 text-lg mb-4">Event not found or failed to load</p>
        <button 
          onClick={() => navigate('/')}
          className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700 transition"
        >
          Return to Events List
        </button>
      </div>
    );
  }

  const { eventData: event, eventParticipants, tasks, comments } = eventData;
  
  // Проверяем, является ли текущий пользователь организатором или администратором
  const isOrganizer = currentUser && event.created_by === currentUser.id;
  const userParticipant = currentUser ? eventParticipants.participants.find(p => p.user.id === currentUser.id) : null;
  const isAdmin = userParticipant?.role === 'admin';
  const canManageParticipants = isOrganizer || isAdmin;

  return (
    <div className="max-w-6xl mx-auto relative">
      {/* Event Header */}
      <div className="flex justify-between items-center mb-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{event.title}</h1>
          <p className="text-sm text-gray-500">
            {new Date(event.start_date).toLocaleString()} - {new Date(event.end_date).toLocaleString()}
          </p>
        </div>
        {isOrganizer && (
          <button 
            onClick={handleDeleteEvent}
            className="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition"
          >
            Delete Event
          </button>
        )}
      </div>

      {/* Event Details */}
      <div className="bg-white rounded-lg shadow-md overflow-hidden mb-6">
        <div className="border-b border-gray-200 bg-gray-50 p-4">
          <h2 className="text-lg font-medium text-gray-900">
            Информация о событии
          </h2>
        </div>
        <div className="p-4">
          <div className="flex items-center text-gray-600 mb-2">
            <span className="font-semibold mr-2">Место проведения:</span>
            <span>{event.location || 'Не указано'}</span>
          </div>
          <div className="mt-4">
            <h3 className="text-md font-semibold text-gray-900 mb-2">Описание</h3>
            <p className="text-gray-600">{event.description || 'Описание отсутствует'}</p>
          </div>
        </div>
      </div>

      {/* Two column layout */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Left column - Participants */}
        <div className="md:col-span-1">
          <div className="bg-white rounded-lg shadow-md overflow-hidden">
            <div className="border-b border-gray-200 bg-gray-50 p-4">
              <div className="flex justify-between items-center">
                <h2 className="text-lg font-medium text-gray-900">
                  Участники ({eventParticipants.total})
                </h2>
                {canManageParticipants && (
                  <button
                    onClick={() => setShowAddParticipantModal(true)}
                    className="bg-blue-600 text-white px-3 py-1 rounded text-sm hover:bg-blue-700 transition"
                  >
                    Добавить
                  </button>
                )}
              </div>
            </div>
            <div className="p-4">
              {eventParticipants.participants.length > 0 ? (
                <ul className="divide-y divide-gray-200">
                  {eventParticipants.participants.map((participant: EventParticipant) => (
                    <li key={participant.id} className="py-3">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="font-medium text-gray-800">{participant.user.username}</p>
                          <p className="text-sm text-gray-500">{participant.role}</p>
                        </div>
                        {canManageParticipants && participant.user.id !== currentUser?.id && (
                          <button 
                            onClick={() => handleRemoveParticipant(participant.id)}
                            className="text-red-600 hover:text-red-800 text-sm"
                          >
                            Удалить
                          </button>
                        )}
                      </div>
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-gray-500">Участников пока нет</p>
              )}
            </div>
          </div>

          {/* Comments Section - moved under participants */}
          <div className="bg-white rounded-lg shadow-md overflow-hidden mt-6">
            <div className="border-b border-gray-200 bg-gray-50 p-4">
              <h2 className="text-lg font-medium text-gray-900">
                Обсуждение ({comments.comments.length})
              </h2>
            </div>
            
            {/* Add Comment Form */}
            {currentUser && (
              <div className="p-4 bg-white">
                <form onSubmit={handleCreateComment} className="relative">
                  <div className="flex items-start space-x-3">
                    <div className="flex-shrink-0">
                      <div className="h-10 w-10 rounded-full bg-blue-600 flex items-center justify-center text-white font-bold">
                        {currentUser.username?.charAt(0).toUpperCase() || "У"}
                      </div>
                    </div>
                    <div className="min-w-0 flex-1">
                      <textarea
                        placeholder="Напишите ваш комментарий..."
                        value={newCommentContent}
                        onChange={(e) => setNewCommentContent(e.target.value)}
                        className="w-full border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-700 placeholder-gray-400"
                        rows={2}
                        required
                      ></textarea>
                      <div className="mt-2 flex justify-end">
                        <button 
                          type="submit" 
                          className="px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 text-white transition"
                        >
                          Отправить
                        </button>
                      </div>
                    </div>
                  </div>
                </form>
              </div>
            )}
            
            {/* Comments List */}
            <div className="p-4">
              {comments.comments.length > 0 ? (
                <ul className="divide-y divide-gray-200">
                  {comments.comments
                    .filter((c: Comment) => !c.isDeleted)
                    .sort((a: Comment, b: Comment) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
                    .map((comment: Comment) => (
                      <li key={comment.commentId} className="py-4">
                        <div className="flex items-start space-x-3">
                          <div className="flex-shrink-0">
                            <div className="h-10 w-10 rounded-full bg-gray-200 flex items-center justify-center text-gray-600 font-bold">
                              {(userMap.get(comment.senderId) || `У#${comment.senderId}`)?.charAt(0).toUpperCase()}
                            </div>
                          </div>
                          <div className="min-w-0 flex-1">
                            <div className="flex justify-between items-baseline">
                              <h3 className="font-medium text-gray-800">
                                {userMap.get(comment.senderId) || `Пользователь #${comment.senderId}`}
                              </h3>
                              <div className="text-sm text-gray-400 whitespace-nowrap">
                                {new Date(comment.createdAt).toLocaleString('ru-RU', {
                                  year: 'numeric',
                                  month: 'short',
                                  day: 'numeric',
                                  hour: '2-digit',
                                  minute: '2-digit'
                                })}
                              </div>
                            </div>
                            <p className="text-gray-600 mt-1 break-words whitespace-pre-wrap">
                              {comment.content}
                            </p>
                          </div>
                        </div>
                      </li>
                    ))}
                </ul>
              ) : (
                <div className="flex flex-col items-center justify-center py-8 text-gray-500">
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-12 w-12 text-gray-300 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                  <p>Комментариев пока нет. Будьте первым!</p>
                </div>
              )}
            </div>
          </div>
        </div>
        
        {/* Right column - Tasks, Expenses */}
        <div className="md:col-span-2 space-y-6">
          {/* Tasks Section */}
          <div className="bg-white rounded-lg shadow-md overflow-hidden">
            <div className="border-b border-gray-200 bg-gray-50 p-4">
              <div className="flex justify-between items-center">
                <h2 className="text-lg font-medium text-gray-900">
                  Задачи ({tasks.tasks.length})
                </h2>
                <Link to={`/events/${eventId}/tasks`} className="text-blue-600 hover:text-blue-800 text-sm">
                  Посмотреть все задачи
                </Link>
              </div>
            </div>
            <div className="p-4">
              {tasks.tasks.length > 0 ? (
                <ul className="divide-y divide-gray-200">
                  {tasks.tasks
                    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
                    .slice(0, 3)
                    .map(task => (
                      <li key={task.id} className="py-3">
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <div className="flex items-center">
                              <span 
                                className={`w-3 h-3 rounded-full mr-2 ${task.status === TaskStatus.Completed ? 'bg-green-500' : task.status === TaskStatus.InProgress ? 'bg-yellow-500' : 'bg-gray-400'}`}
                              ></span>
                              <p className="font-medium text-gray-800">{task.title}</p>
                            </div>
                            {task.description && (
                              <p className="text-sm text-gray-600 mt-1">{task.description}</p>
                            )}
                            <div className="flex items-center mt-2 text-sm">
                              <span className="text-gray-500 mr-4">Статус: {task.status.replace('_', ' ')}</span>
                              {task.assigned_to && (
                                <span className="text-gray-500">Исполнитель: {userMap.get(task.assigned_to) || `Пользователь #${task.assigned_to}`}</span>
                              )}
                            </div>
                          </div>
                          <div className="flex items-center space-x-2 ml-4">
                            <select
                              className="text-sm border rounded p-1"
                              value={task.status}
                              onChange={(e) => handleUpdateTaskStatus(task.id, e.target.value as TaskStatus)}
                            >
                              <option value={TaskStatus.Pending}>Ожидает</option>
                              <option value={TaskStatus.InProgress}>В процессе</option>
                              <option value={TaskStatus.Completed}>Завершено</option>
                              <option value={TaskStatus.Cancelled}>Отменено</option>
                            </select>
                          </div>
                        </div>
                      </li>
                    ))}
                </ul>
              ) : (
                <p className="text-gray-500">Задач пока нет</p>
              )}
              
              <div className="mt-4 flex justify-end">
                <Link 
                  to={`/events/${eventId}/tasks`} 
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition flex items-center"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clipRule="evenodd" />
                  </svg>
                  Добавить задачу
                </Link>
              </div>
            </div>
          </div>

          {/* Expenses Section */}
          <div className="bg-white rounded-lg shadow-md overflow-hidden">
            <div className="border-b border-gray-200 bg-gray-50 p-4">
              <div className="flex justify-between items-center">
                <h2 className="text-lg font-medium text-gray-900">
                  Расходы
                </h2>
                <Link to={`/events/${eventId}/expenses`} className="text-blue-600 hover:text-blue-800 text-sm">
                  Посмотреть все расходы
                </Link>
              </div>
            </div>
            <div className="p-4">
              {eventData?.expenses?.items?.length > 0 ? (
                <div>
                  <div className="mb-4">
                    <div className="flex justify-between items-center p-3 bg-gray-50 rounded-lg">
                      <div>
                        <h3 className="text-gray-700 font-medium">Общая сумма расходов:</h3>
                        <p className="text-2xl font-bold text-green-600">
                          {eventData.balanceReport.total_amount.toFixed(2)} {eventData.expenses.items[0]?.currency || 'RUB'}
                        </p>
                      </div>
                      <div className="text-right">
                        <h3 className="text-gray-700 font-medium">Количество расходов:</h3>
                        <p className="text-xl font-bold">{eventData.expenses.total_count}</p>
                      </div>
                    </div>
                  </div>
                  
                  {eventData.balanceReport.user_balances.length > 0 && (
                    <div className="mt-4">
                      <h3 className="text-gray-700 font-medium mb-2">Балансы участников:</h3>
                      <ul className="divide-y divide-gray-200">
                        {eventData.balanceReport.user_balances
                          .filter(balance => Math.abs(balance.balance) > 0.01)
                          .map(balance => (
                            <li key={balance.user_id} className="py-2 flex justify-between items-center">
                              <span className="font-medium">{balance.username}</span>
                              <span className={`font-medium ${balance.balance > 0 ? 'text-green-600' : 'text-red-500'}`}>
                                {balance.balance > 0 ? '+' : ''}{balance.balance.toFixed(2)}
                              </span>
                            </li>
                          ))}
                      </ul>
                    </div>
                  )}
                  
                  <div className="flex items-center justify-center py-4 mt-2">
                    <Link 
                      to={`/events/${eventId}/expenses`} 
                      className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition flex items-center"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clipRule="evenodd" />
                      </svg>
                      Управление расходами
                    </Link>
                  </div>
                </div>
              ) : (
                <div className="text-center py-6">
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-12 w-12 mx-auto text-gray-400 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <p className="text-gray-500 mb-4">Расходов пока нет</p>
                  <Link 
                    to={`/events/${eventId}/expenses`} 
                    className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition inline-flex items-center"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clipRule="evenodd" />
                    </svg>
                    Добавить расход
                  </Link>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Модальное окно для добавления участника */}
      {showAddParticipantModal && (
        <div className="fixed inset-0 overflow-y-auto z-50 flex items-center justify-center">
          <div className="fixed inset-0 bg-black bg-opacity-50" onClick={() => setShowAddParticipantModal(false)}></div>
          <div className="bg-white rounded-lg shadow-xl z-10 max-w-md w-full mx-4">
            <div className="px-4 py-5 sm:px-6 border-b">
              <h3 className="text-lg font-medium text-gray-900">Add Participant</h3>
            </div>
            <div className="px-4 py-5 sm:p-6">
              <form onSubmit={handleAddParticipant}>
                <div className="space-y-4">
                  <div>
                    <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                      User ID or Username *
                    </label>
                    <input
                      type="text"
                      id="username"
                      value={newParticipantData.username}
                      onChange={(e) => setNewParticipantData({...newParticipantData, username: e.target.value})}
                      className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                      required
                    />
                    <p className="mt-1 text-sm text-gray-500">Enter user ID or username to invite.</p>
                  </div>
                  
                  <div>
                    <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                      Email (Optional)
                    </label>
                    <input
                      type="email"
                      id="email"
                      value={newParticipantData.email}
                      onChange={(e) => setNewParticipantData({...newParticipantData, email: e.target.value})}
                      className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                    />
                    <p className="mt-1 text-sm text-gray-500">If the user is not registered, we can send an invitation.</p>
                  </div>
                </div>
                
                <div className="mt-5 sm:mt-6 flex justify-end space-x-3">
                  <button
                    type="button"
                    onClick={() => setShowAddParticipantModal(false)}
                    className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none"
                  >
                    Add Participant
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default EventDetails;
