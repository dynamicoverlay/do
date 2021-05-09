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


const SIGNUP = gql`
    mutation Signup($username: String!, $email: String!, $password: String!){
        createUser(email: $email, password: $password, username: $username){
            token
        }
    }
`


type EmailSignupState = {
    loading: boolean
}

export default () => {
    const history = useHistory();
    const [doSignup, {data, error:signupError, loading:loggingIn}] = useMutation(SIGNUP, { errorPolicy: 'all' , onError: () => {}, onCompleted: (data) => {console.log(data); localStorage.setItem('token', data.createUser.token); history.push('/dashboard');}});
    const {value:emailValue, bind:emailBind, reset:emailReset} = useInput('');
    const {value:passwordValue, bind:passwordBind, reset:passwordReset} = useInput('');
    const {value:usernameValue, bind:usernameBind, reset:usernameReset} = useInput('');
    const [errors, setErrors] = useState({email: false, password: false, username: false});
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
        if(usernameValue.length <= 0) {
            copyOfErrors = {...copyOfErrors, username: true}
            hasErrored = true;
        } else {
            copyOfErrors = {...copyOfErrors, username: false}
            if(!hasErrored){
                hasErrored = false;
            }
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
        {signupError && signupError.graphQLErrors.map((err, index) => <ErrorBox key={index}>{err.message}</ErrorBox>)}
        <form className="flex flex-col justify-center w-full" onSubmit={() => { 
            if(validateForm()){
                localStorage.removeItem('token');
                doSignup({variables: {email: emailValue, password: passwordValue, username: usernameValue}})
            }
        }}>
            <TextInput type="email" placeholder="Email" {...emailBind} error={errors.email}></TextInput>
            <TextInput type="text" placeholder="Username" {...usernameBind} error={errors.username}></TextInput>
            <TextInput type="password" placeholder="Password" {...passwordBind} error={errors.password}></TextInput>
            <Button color="green-700" type={"submit"}>SUBMIT</Button>
        </form>
        </>)
    }
    return content
}
