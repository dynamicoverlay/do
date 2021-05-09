import React from 'react';
import ReactDOM from 'react-dom';
import './assets/main.css';
import './assets/styles.css';
import { ApolloProvider } from '@apollo/react-hooks';
import ApolloClient from 'apollo-client';

import { ApolloLink } from 'apollo-link';
import { HttpLink } from 'apollo-link-http';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { setContext } from 'apollo-link-context';


import { BrowserRouter, Switch, Route } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import PageLayout from './components/layouts/PageLayout';
import LogoutPage from './pages/dashboard/LogoutPage';
import SignupPage from './pages/SignupPage';
import LandingPage from './pages/LandingPage';

import {ProtectedRoute} from './ProtectedRoute';


const httpLink = new HttpLink({uri: "http://localhost:8080/query"})

const cache = new InMemoryCache({});

const useAppApolloClient = () => {
    const authLink = setContext((request,  previousContext) => {
        const authToken = localStorage.getItem("token");
        return {
            headers: {authorization: `${authToken ? authToken : ''}`} 
        }
    })
    return new ApolloClient({
        link: authLink.concat(httpLink),
        cache
    });
};


function Router() {
    const client = useAppApolloClient();
    return (
        <ApolloProvider client={client}>
            <BrowserRouter>
                <Switch>
                    <Route exact path="/" component={LandingPage}></Route>
                    <Route path="/login" component={LoginPage}></Route>
                    <Route path="/signup" component={SignupPage}></Route>
                    <ProtectedRoute path="/dashboard" component={PageLayout}></ProtectedRoute>
                </Switch>
            </BrowserRouter>
        </ApolloProvider>
    )
}
ReactDOM.render(<Router />, document.querySelector('#app'));
