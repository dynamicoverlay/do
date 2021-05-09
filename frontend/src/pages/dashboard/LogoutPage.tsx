import React from 'react';
import {Redirect} from 'react-router-dom';
import {useApolloClient} from '@apollo/react-hooks';

export default () => {
    const client = useApolloClient();
    localStorage.removeItem('token');
    client.resetStore();
    return (
        <Redirect to="/login" />
    )
};
