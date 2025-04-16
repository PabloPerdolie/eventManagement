import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Calendar, MapPin, Clock, Users, Image } from 'lucide-react';
import { eventService } from '../services/eventService';
import { EventResponse } from '../types/api';
import { toast } from 'react-toastify';

// Кастомное изображение по умолчанию
const DEFAULT_EVENT_IMAGE = 'https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?q=80&w=800&auto=format&fit=crop';


function EventsList() {
  // Состояние для хранения информации о сбоях загрузки изображений
  const [failedImages, setFailedImages] = useState<{ [key: number]: boolean }>({});
  const [events, setEvents] = useState<EventResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalEvents, setTotalEvents] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  const pageSize = 9; // 3x3 grid

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        setLoading(true);
        const response = await eventService.getEvents(page, pageSize);
        setEvents(prev => page === 1 ? response.events : [...prev, ...response.events]);
        setTotalEvents(response.total);
        setHasMore(response.events.length === pageSize && response.total > page * pageSize);
      } catch (error) {
        console.error('Error fetching events:', error);
        toast.error('Failed to load events');
      } finally {
        setLoading(false);
      }
    };

    fetchEvents();
  }, [page]);

  const loadMore = () => {
    if (hasMore && !loading) {
      setPage(prev => prev + 1);
    }
  };

  return (
    <div className="max-w-6xl mx-auto relative">
      <div className="flex justify-between items-center mb-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">События</h1>
          <p className="text-sm text-gray-500">
            {totalEvents > 0 && 
              `Показано ${Math.min(events.length, totalEvents)} из ${totalEvents} событий`
            }
          </p>
        </div>
        <Link
          to="/events/create"
          className="bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center hover:bg-blue-700 transition"
        >
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Создать событие
        </Link>
      </div>

      {loading && events.length === 0 ? (
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
        </div>
      ) : events.length === 0 ? (
        <div className="bg-white rounded-lg shadow-md p-8 text-center">
          <p className="text-gray-500 text-lg">Событий пока нет. Создайте первое событие!</p>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-md overflow-hidden">
          <div className="border-b border-gray-200 bg-gray-50 p-4">
            <h2 className="text-lg font-medium text-gray-900">
              Предстоящие события
            </h2>
          </div>
          
          <div className="p-4">
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {events.map((event) => (
                <Link
                  key={event.id}
                  to={`/events/${event.id}`}
                  className="block bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow overflow-hidden"
                >
                  <div className="aspect-video w-full overflow-hidden">
                    {failedImages[event.id] ? (
                      <div className="w-full h-full flex items-center justify-center bg-gray-100">
                        <div className="text-center p-4">
                          <Image className="w-12 h-12 mx-auto text-gray-400 mb-2" />
                          <p className="text-gray-500 text-sm">{event.title}</p>
                        </div>
                      </div>
                    ) : (
                      <img
                        src={`https://source.unsplash.com/random/800x600/?event&sig=${event.id}`}
                        alt={event.title}
                        className="w-full h-full object-cover"
                        onError={() => {
                          setFailedImages(prev => ({ ...prev, [event.id]: true }));
                        }}
                      />
                    )}
                  </div>
                  <div className="p-4">
                    <h3 className="text-xl font-semibold text-gray-900 mb-2">{event.title}</h3>
                    <p className="text-gray-600 mb-4 line-clamp-2">{event.description || 'Описание отсутствует'}</p>
                    <div className="space-y-2 text-sm text-gray-500">
                      <div className="flex items-center">
                        <Calendar className="w-4 h-4 mr-2" />
                        <span>Начало: {new Date(event.start_date).toLocaleDateString('ru-RU')}</span>
                      </div>
                      <div className="flex items-center">
                        <Clock className="w-4 h-4 mr-2" />
                        <span>Окончание: {new Date(event.end_date).toLocaleDateString('ru-RU')}</span>
                      </div>
                      <div className="flex items-center">
                        <MapPin className="w-4 h-4 mr-2" />
                        <span>{event.location || 'Место не указано'}</span>
                      </div>
                      <div className="flex items-center">
                        <Users className="w-4 h-4 mr-2" />
                        <span>Организатор: Пользователь #{event.created_by}</span>
                      </div>
                    </div>
                  </div>
                </Link>
              ))}
            </div>

            {hasMore && (
              <div className="mt-8 text-center">
                <button
                  onClick={loadMore}
                  disabled={loading}
                  className="px-6 py-2 rounded-lg bg-gray-100 hover:bg-gray-200 transition-colors disabled:opacity-50"
                >
                  {loading ? 'Загрузка...' : 'Загрузить больше событий'}
                </button>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

export default EventsList;