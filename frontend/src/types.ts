export interface Overlay {
    identifier: string;
    name: string;
    pin: string;
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
