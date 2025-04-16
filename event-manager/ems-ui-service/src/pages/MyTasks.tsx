import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';
import { taskService } from '../services/taskService';
import { eventService } from '../services/eventService';
import { TaskResponse, TaskStatus, EventResponse } from '../types/api';
import { CheckCircle, Clock, AlertCircle, XCircle } from 'lucide-react';
import { toast } from 'react-toastify';

const MyTasks: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [tasks, setTasks] = useState<TaskResponse[]>([]);
  const [events, setEvents] = useState<Map<number, EventResponse>>(new Map());
  const [users, setUsers] = useState<Map<number, string>>(new Map());
  const [page, setPage] = useState(1);
  const [totalTasks, setTotalTasks] = useState(0);
  const pageSize = 20;
  const { user } = useSelector((state: RootState) => state.auth);

  useEffect(() => {
    fetchUserTasks();
  }, [page]);

  const fetchUserTasks = async () => {
    try {
      setLoading(true);
      const data = await taskService.getTasks(page, pageSize);
      console.log('User tasks data:', data);
      setTasks(data.tasks || []);
      setTotalTasks(data.total || 0);

      // Fetch event details for all task event_ids
      const eventIds = [...new Set(data.tasks.map(task => task.event_id))];
      const eventMap = new Map<number, EventResponse>();
      const userIds = [...new Set(data.tasks.map(task => task.assigned_to))];
      const userMap = new Map<number, string>();

      // Fetch event details for all unique event IDs
      await Promise.all(eventIds.map(async (eventId) => {
        try {
          const eventData = await eventService.getEvent(eventId);
          eventMap.set(eventId, eventData.eventData);
          
          // Also collect username mappings from event participants
          eventData.eventParticipants.participants.forEach(participant => {
            userMap.set(participant.user.id, participant.user.username);
          });
        } catch (error) {
          console.error(`Error fetching event ${eventId}:`, error);
        }
      }));

      setEvents(eventMap);
      setUsers(userMap);
    } catch (error) {
      console.error('Error fetching user tasks:', error);
      toast.error('Failed to load your tasks');
    } finally {
      setLoading(false);
    }
  };

  const getStatusIcon = (status: TaskStatus) => {
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

  const getStatusLabel = (status: TaskStatus) => {
    switch (status) {
      case TaskStatus.Pending:
        return 'Pending';
      case TaskStatus.InProgress:
        return 'In Progress';
      case TaskStatus.Completed:
        return 'Completed';
      case TaskStatus.Cancelled:
        return 'Cancelled';
      default:
        return 'Unknown';
    }
  };

  const getStatusColorClass = (status: TaskStatus) => {
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

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const handlePreviousPage = () => {
    if (page > 1) {
      setPage(page - 1);
    }
  };

  const handleNextPage = () => {
    if (page * pageSize < totalTasks) {
      setPage(page + 1);
    }
  };

  if (loading && tasks.length === 0) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    );
  }

  return (
    <div className="max-w-5xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">My Assigned Tasks</h1>
      </div>

      {tasks.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-6 text-center">
          <p className="text-gray-500 mb-4">You don't have any assigned tasks yet.</p>
          <Link
            to="/"
            className="inline-block bg-indigo-100 text-indigo-700 px-4 py-2 rounded-md hover:bg-indigo-200 transition"
          >
            View My Events
          </Link>
        </div>
      ) : (
        <>
          <div className="bg-white rounded-lg shadow-md overflow-hidden">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Task
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Event
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Assigned To
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Created
                  </th>
                  <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {tasks.map((task) => (
                  <tr key={task.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-medium text-gray-900">{task.title}</div>
                      {task.description && (
                        <div className="text-sm text-gray-500">{task.description.substring(0, 50)}{task.description.length > 50 ? '...' : ''}</div>
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        {events.get(task.event_id) ? (
                          <Link 
                            to={`/events/${task.event_id}`}
                            className="text-indigo-600 hover:text-indigo-900"
                          >
                            {events.get(task.event_id)?.title}
                          </Link>
                        ) : (
                          `Event #${task.event_id}`
                        )}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        {users.get(task.assigned_to) || `User #${task.assigned_to}`}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        {getStatusIcon(task.status)}
                        <span className={`ml-2 px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColorClass(task.status)}`}>
                          {getStatusLabel(task.status)}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {formatDate(task.created_at)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <Link 
                        to={`/events/${task.event_id}`}
                        className="text-indigo-600 hover:text-indigo-900 mr-4"
                      >
                        View Event
                      </Link>
                      <Link 
                        to={`/events/${task.event_id}/tasks`}
                        className="text-indigo-600 hover:text-indigo-900"
                      >
                        All Tasks
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Pagination controls */}
          {totalTasks > pageSize && (
            <div className="flex justify-between items-center mt-8">
              <button
                onClick={handlePreviousPage}
                disabled={page === 1}
                className={`px-4 py-2 rounded-md ${
                  page === 1
                    ? 'bg-gray-200 text-gray-400 cursor-not-allowed'
                    : 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200'
                }`}
              >
                Previous
              </button>
              <span className="text-gray-600">
                Page {page} of {Math.ceil(totalTasks / pageSize)}
              </span>
              <button
                onClick={handleNextPage}
                disabled={page * pageSize >= totalTasks}
                className={`px-4 py-2 rounded-md ${
                  page * pageSize >= totalTasks
                    ? 'bg-gray-200 text-gray-400 cursor-not-allowed'
                    : 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200'
                }`}
              >
                Next
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default MyTasks;
