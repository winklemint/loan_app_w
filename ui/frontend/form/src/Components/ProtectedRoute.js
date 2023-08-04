import React from 'react';
import { Route } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';

const isAuthenticated = () => {
  return true;
};

const ProtectedRoute = ({ component: Component, ...rest }) => {
  const isAuthenticatedUser = isAuthenticated();
  const navigate = useNavigate();

  if (!isAuthenticatedUser) {
    navigate('/login');
    return null;
  }

  return <Route {...rest} element={<Component />} />;
};

export default ProtectedRoute;
