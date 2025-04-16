import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { Calendar, LogOut, User, ListTodo, CalendarDays } from 'lucide-react';
import { RootState } from '../store/store';
import { logout } from '../store/slices/authSlice';

const Navbar: React.FC = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { isAuthenticated, user } = useSelector((state: RootState) => state.auth);

  const handleLogout = async () => {
    await dispatch(logout());
    navigate('/login');
  };

  if (!isAuthenticated) return null;

  return (
    <nav className="bg-white shadow-lg">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <Link to="/" className="flex items-center space-x-2">
            <Calendar className="h-6 w-6 text-indigo-600" />
            <span className="text-xl font-semibold">EventHub</span>
          </Link>

          <div className="flex items-center space-x-6">
            <Link to="/" className="flex items-center space-x-1 text-gray-600 hover:text-gray-800">
              <CalendarDays className="h-5 w-5" />
              <span>My Events</span>
            </Link>
            <Link to="/my-tasks" className="flex items-center space-x-1 text-gray-600 hover:text-gray-800">
              <ListTodo className="h-5 w-5" />
              <span>My Tasks</span>
            </Link>
            <div className="flex items-center space-x-2">
              <User className="h-5 w-5 text-gray-600" />
              <span className="text-gray-800">{user?.username}</span>
            </div>
            <button
              onClick={handleLogout}
              className="flex items-center space-x-1 text-gray-600 hover:text-gray-800"
            >
              <LogOut className="h-5 w-5" />
              <span>Logout</span>
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;