import React from 'react';
import { Route, Redirect, RouteProps } from 'react-router-dom';
import {useUserQuery} from './hooks/useUserQuery';
import LoadingPage from './pages/LoadingPage';

export const ProtectedRoute = ({component: Component, ...rest}: RouteProps) => {
    const authToken = localStorage.getItem("token");
    const userData = useUserQuery();
    if (userData.loading) return (<Route {...rest} component={LoadingPage} />)
    if (!Component) return null;
    return (
        <Route {...rest} render={props => (authToken && userData.data ? <Component {...props} /> : <Redirect to="/login" />)}/> 
    );
}


