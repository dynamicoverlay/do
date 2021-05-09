import { useState, useEffect, useRef, MutableRefObject } from "react";
import {BaseMessage, UpdateStateMessage, GlobalState} from '../types';

export function useOverlayState(ws: MutableRefObject<WebSocket | null>, id: string, pin: string){
    const [overlayState, setOverlayState] = useState({} as GlobalState);
    useEffect(() => {
        if(pin.length > 0){
            ws.current = new WebSocket('ws://localhost:8083/ws');
            ws.current?.addEventListener('open', (event) => {
                console.log("Connected to Websocket server");
                ws.current?.send(JSON.stringify({overlay: id, passcode: pin}));
            });
            ws.current?.addEventListener('message', (event) => {
                try {
                    let message = JSON.parse(event.data) as BaseMessage;
                    if (message.message === "authRequest" && message.success){
                        console.log("Successfully authenticated with the server")
                    } else if (message.message === "updateState"){
                        let updateStateMessage = message as UpdateStateMessage;
                        console.log("Received updateState request", message);
                        if(updateStateMessage.value){
                            setOverlayState(updateStateMessage.value)
                        }
                    }
                } catch (e) {
                    console.error(e);
                    console.log("Received message that was not JSON, this is unexpected")
                    return
                }
                console.log("Received message from server", event.data);
            });
            ws.current?.addEventListener('close', (event) => {
                console.log("Disconnected");
            });
            ws.current?.addEventListener('error', (event) => {
                console.log("Errored", event)
            });
            return () => {
                ws.current?.close();
            }
        }
    }, [pin])
    return overlayState;
}