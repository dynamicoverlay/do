
export interface AppProps {
    overlayID: string;
    overlayPin: string;
}

export interface BaseMessage {
    success: boolean;
    message?: string;
    value?: string;
}

export interface UpdateStateMessage {
    success: boolean;
    message?: string;
    value?: GlobalState;
}

export interface GlobalState {
    [index: string]: {};
}

export interface AppState {
    authenticated: boolean;
    connected: boolean;
    globalState: GlobalState;
}

export interface BackdropState {
    visible: boolean;
    text?: string;
    transitionTime: number;
    color?: string;
    backgroundImage?: string;
    textColor?: string;
}

export interface BackdropProps {
    state: BackdropState | undefined;
}
