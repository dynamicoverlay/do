import React, { Component } from 'react';
import { CSSTransition,  SwitchTransition} from 'react-transition-group';
import EmailLogin from '../components/login/EmailLogin';
import LoginButtons from '../components/login/LoginButtons';

type LoginState = {
    loginType: string,
    emailAddress: string,
    password: string
}

class LoginPage extends Component<{}, LoginState> {

    constructor(props: {}) {
        super(props);
        this.state = {
            loginType: "none",
            emailAddress: "",
            password: ""
        }
    }

    handleEmailClick = () => {
        this.setState({
            loginType: "email"
        })
    }

    goBack = () => {
        this.setState({
            loginType: "none"
        })
    }

    render() {
        return (
            <>
                <div className="flex flex-col justify-center text-theme items-center min-h-screen min-w-screen bg-background theme-dark">
                    <div className="mx-auto w-3/4 lg:w-1/2 text-center ">
                        <h1 className="text-6xl font-semibold">My<span className="font-light">Overlay</span></h1>
                        {/* <h2 className="text-2xl mb-auto font-normal">Login</h2> */}
                        <div className="p-4 mt-4 w-1/2 lg:w-1/3 mx-auto rounded flex flex-row bg-gray-1000 font-light shadow-xl text-white">
                            <SwitchTransition mode={"out-in"}>
                                <CSSTransition key={this.state.loginType} addEndListener={(node, done) => {node.addEventListener("transitionend", done, false);}} classNames={"email-login"} >
                                    <div className="flex flex-col justify-center w-full">
                                        {this.state.loginType === "none" && <LoginButtons goToEmail={this.handleEmailClick}></LoginButtons>}
                                        {this.state.loginType === "email" && <EmailLogin goBack={this.goBack}></EmailLogin>}
                                    </div>
                                </CSSTransition>
                            </SwitchTransition>
                        </div>
                    </div>
                </div>
            </>
        )
    }

}
export default LoginPage