import React from 'react';
import { Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';

interface PrivateRouteProps {
  element: React.ReactElement;
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ element }) => {
  const { isAuthenticated, loading } = useSelector((state: RootState) => state.auth);

  // Show loading indicator while checking authentication
  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <p className="text-gray-600 text-lg">Проверка авторизации...</p>
      </div>
    );
  }

  return isAuthenticated ? element : <Navigate to="/login" replace />;
};

export default PrivateRoute;