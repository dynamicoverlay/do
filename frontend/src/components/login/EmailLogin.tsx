import React, { useState} from 'react';
import gql from 'graphql-tag';
import TextInput from '../form/TextInput';
import Button from '../Button';
import Spinner from '../Spinner';
import {useMutation} from '@apollo/react-hooks';
import ErrorBox from '../ErrorBox';
import { error } from 'console';

import { useHistory } from 'react-router-dom';
import { useInput } from '../../hooks/useInput';


const LOGIN = gql`
    mutation Login($email: String!, $password: String!){
        login(email: $email, password: $password){
            token
        }
    }
`

type EmailLoginProps = {
    goBack: (event: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => void,
}

type EmailLoginState = {
    loading: boolean
}

export default (props: EmailLoginProps) => {
    const history = useHistory();
    const [doLogin, {data, error:loginError, loading:loggingIn}] = useMutation(LOGIN, { errorPolicy: 'all' , onError: () => {}, onCompleted: (data) => {localStorage.setItem('token', data.login.token); history.push('/dashboard');}});
    const {value:emailValue, bind:emailBind, reset:emailReset} = useInput('');
    const {value:passwordValue, bind:passwordBind, reset:passwordReset} = useInput('');
    const [errors, setErrors] = useState({email: false, password: false});
    const validateForm = (): boolean => {
        let hasErrored = false;
        let copyOfErrors = errors;
        if(emailValue.length <= 0 || emailValue.indexOf('@') === -1 || emailValue.indexOf('.') === -1){
            copyOfErrors = {...copyOfErrors, email: true}
            hasErrored = true;
        } else {
            hasErrored = false;
            copyOfErrors = {...copyOfErrors, email: false}
        }
        if(passwordValue.length <= 0) {
            copyOfErrors = {...copyOfErrors, password: true}
            hasErrored = true;
        } else {
            copyOfErrors = {...copyOfErrors, password: false}
            if(!hasErrored){
                hasErrored = false;
            }
        }
        setErrors(copyOfErrors);
        return !hasErrored;
    }
    let content;
    if(loggingIn){
        content = (<Spinner></Spinner>)
    } else {
        content = (<>
        {loginError && loginError.graphQLErrors.map((err, index) => <ErrorBox key={index}>{err.message}</ErrorBox>)}
        <form className="flex flex-col justify-center w-full" onSubmit={(event) => {
            event.preventDefault();
                if(validateForm()){
                    localStorage.removeItem('token');
                    doLogin({variables: {email: emailValue, password: passwordValue}})
                }
        }}>
            <TextInput type="email" placeholder="Email" {...emailBind} error={errors.email}></TextInput>
            <TextInput type="password" placeholder="Password" {...passwordBind} error={errors.password}></TextInput>
            <Button type={"submit"} color="green-700">SUBMIT</Button>
            <Button color="gray-700" onClick={props.goBack}>GO BACK</Button>
        </form>
        </>)
    }
    return content
}
