import React, { Component } from 'react';
import { CSSTransition,  SwitchTransition} from 'react-transition-group';
import EmailSignup from '../components/signup/EmailSignup';

type SignupState = {
    emailAddress: string,
    password: string
}

export default class SignupPage extends Component<{}, SignupState> {

    constructor(props: {}) {
        super(props);
        this.state = {
            emailAddress: "",
            password: ""
        }
    }

    render() {
        return (
            <>
                <div className="flex flex-col justify-center text-theme items-center min-h-screen min-w-screen">
                    <div className="mx-auto w-3/4 lg:w-1/2 text-center ">
                        <h1 className="text-6xl font-semibold">My<span className="font-light">Overlay</span></h1>
                        {/* <h2 className="text-2xl mb-auto font-normal">Signup</h2> */}
                        <div className="p-4 mt-4 w-1/2 lg:w-1/3 mx-auto rounded flex flex-row bg-gray-1000 font-light shadow-xl text-white">
                            <div className="flex flex-col justify-center w-full">
                                <EmailSignup />
                            </div>
                       </div>
                    </div>
                </div>
            </>
        )
    }

}



