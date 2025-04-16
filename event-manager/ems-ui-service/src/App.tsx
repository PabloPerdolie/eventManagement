import React, { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Provider, useDispatch, useSelector } from 'react-redux';
import { RootState } from './store/store';
import { ToastContainer } from 'react-toastify';
import { store } from './store/store';
import { getCurrentUser } from './store/slices/authSlice';
import Cookies from 'js-cookie';
import PrivateRoute from './components/PrivateRoute';
import Navbar from './components/Navbar';
import Login from './pages/Login';
import Register from './pages/Register';
import EventsList from './pages/EventsList';
import EventDetails from './pages/EventDetails';
import CreateEvent from './pages/CreateEvent';
import EventTasks from './pages/EventTasks';
import EventExpenses from './pages/EventExpenses';
import MyTasks from './pages/MyTasks';
import 'react-toastify/dist/ReactToastify.css';

// Component to handle auth initialization
const AuthInitializer: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const dispatch = useDispatch();
  const { isAuthenticated, user } = useSelector((state: RootState) => state.auth);

  useEffect(() => {
    // Only try to get current user if we have a token but no user data
    const hasToken = !!Cookies.get('access_token');
    if (hasToken && isAuthenticated && !user) {
      dispatch(getCurrentUser());
    }
  }, [dispatch, isAuthenticated, user]);

  return <>{children}</>;
};

function App() {
  return (
    <Provider store={store}>
      <Router>
        <AuthInitializer>
          <div className="min-h-screen bg-gray-50">
            <Navbar />
            <main className="container mx-auto px-4 py-8">
              <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/" element={<PrivateRoute element={<EventsList />} />} />
                <Route path="/events/create" element={<PrivateRoute element={<CreateEvent />} />} />
                <Route path="/events/:id" element={<PrivateRoute element={<EventDetails />} />} />
                <Route path="/events/:id/tasks" element={<PrivateRoute element={<EventTasks />} />} />
                <Route path="/events/:id/expenses" element={<PrivateRoute element={<EventExpenses />} />} />
                <Route path="/my-tasks" element={<PrivateRoute element={<MyTasks />} />} />
              </Routes>
            </main>
            <ToastContainer position="bottom-right" />
          </div>
        </AuthInitializer>
      </Router>
    </Provider>
  );
}

export default App;