import * as React from 'react';

import Backdrop from './modules/backdrop';
import {AppProps, AppState, BaseMessage, UpdateStateMessage} from './types';

export default class App extends React.Component<AppProps, AppState> {

    private socket: WebSocket | undefined;

    constructor(props: AppProps){
        super(props);
        this.state = {
            authenticated: false,
            connected: false,
            globalState: {}
        }
        this.connectToSocket();
    }


    connectToSocket(){
        this.socket = new WebSocket('ws://localhost:8083/ws');
        this.socket.addEventListener('open', (event) => {
            console.log("Connected to Websocket server");
            this.setState({connected: true})
            this.socket?.send(JSON.stringify({overlay: this.props.overlayID, passcode: this.props.overlayPin}));
        });
        this.socket.addEventListener('message', (event) => {
            try {
                let message = JSON.parse(event.data) as BaseMessage;
                if (message.message === "authRequest" && message.success){
                    this.setState({authenticated: true})
                    console.log("Successfully authenticated with the server")
                } else if (message.message === "updateState" && this.state.authenticated){
                    let updateStateMessage = message as UpdateStateMessage;
                    console.log("Received updateState request", message);
                    if(updateStateMessage.value){
                        this.setState({globalState: updateStateMessage.value})
                        console.log("Updated state to", this.state.globalState);
                    }
                }
            } catch (e) {
                console.error(e);
                console.log("Received message that was not JSON, this is unexpected")
                return
            }
            console.log("Received message from server", event.data);
        });
        this.socket.addEventListener('close', (event) => {
            console.log("Disconnected");
            this.socket = undefined;
        });
        this.socket.addEventListener('error', (event) => {
            console.log("Errored", event)
        });
    }

    getModuleState<T>(name: string): T {
        return this.state.globalState[name] as T;
    }

    render(){
        if(this.state.connected && this.state.authenticated){
            return (
                <Backdrop state={this.getModuleState('backdrop')} />
            )
        }
        return (
            <div></div>
        );
    }
}
