import React, { CSSProperties, useEffect, useState } from 'react';
import { BackdropState } from '../../../../types';
import Button from '../../../../components/Button';
import Spinner from '../../../../components/Spinner';
import TextInput from '../../../../components/form/TextInput';
import useDebounce from '../../../../hooks/useDebounce';
import {ChromePicker} from 'react-color';

interface BackdropModuleProps {
    state: BackdropState | undefined;
    updateState: (key: string, value: string) => void;
}
const popover: CSSProperties = {
    position: 'absolute' ,
    zIndex: 2,
}
const cover: CSSProperties = {
    position: 'fixed',
    top: '0px',
    right: '0px',
    bottom: '0px',
    left: '0px',
}

export default (props: BackdropModuleProps) => {
    let state: BackdropState = props.state || {visible: false, transitionTime: 1};
    if(state){
        const updateState = () => {
            props.updateState("backdrop", JSON.stringify(state))
        };
        const [transitionTime, setTransitionTime] = useState(state.transitionTime);
        const [text, setText] = useState(state.text);
        const [textColour, setTextColour] = useState(state.textColor);
        const [colour, setColour] = useState(state.color);
        const [image, setImage] = useState(state.backgroundImage);
        const [showTextColourPicker, setShowTextColourPicker] = useState(false);
        const [showColourPicker, setShowColourPicker] = useState(false);
        const debouncedTime = useDebounce(transitionTime, 500);
        const debouncedText = useDebounce(text, 500);
        const debouncedTextColour = useDebounce(textColour, 500);
        const debouncedColour = useDebounce(colour, 500);
        const debouncedImage = useDebounce(image, 500);
        useEffect(
            () => {
                if(debouncedTime) {
                    state.transitionTime = debouncedTime;
                }
                state.text = debouncedText;
                state.textColor = debouncedTextColour;
                state.color = debouncedColour;
                state.backgroundImage = debouncedImage;
                updateState();
            },
            [debouncedTime, debouncedText, debouncedTextColour, debouncedColour, debouncedImage]
        )
        return  (
            <div className="p-2">
                <Button color="blue-600" onClick={(event: any) => {
                    event.preventDefault();
                    state.visible = !state.visible;
                    updateState();
                }}>
                    {state?.visible ? 'Hide' : 'Show'}
                </Button>
                <TextInput type="number" label="Transition Time (s) Max: 10" max={10} placeholder="Transition Time" value={transitionTime.toString()} onChange={(event) => setTransitionTime(Number(event.target.value))}></TextInput>
                <div className={"grid grid-flow-col gap-4"}>
                    <TextInput type="text" label="Colour" placeholder="Colour" value={colour} onChange={(event) => setColour(event.target.value)} onClick={() => setShowColourPicker(true)}></TextInput>
                    { showColourPicker ? <div style={ popover }>
                    <div style={ cover } onClick={() => setShowColourPicker(false) }/>
                    <ChromePicker onChange={color => setColour(color.hex)} color={colour}/>
                    </div> : null }
                    <TextInput type="text" label="Image (URL)" placeholder="https://i.imgur.com/...." value={image} onChange={(event) => setImage(event.target.value)}></TextInput>
                </div>
                <div className={"grid grid-flow-col gap-4"}>
                    <TextInput type="text" label="Text" placeholder="Text" value={text} onChange={(event) => setText(event.target.value)}></TextInput>
                    <TextInput type="text" label="Text Colour" placeholder="Text colour" value={textColour} onChange={(event) => setTextColour(event.target.value)} onClick={() => setShowTextColourPicker(true)}></TextInput>
                    { showTextColourPicker ? <div style={ popover }>
                    <div style={ cover } onClick={() => setShowTextColourPicker(false) }/>
                    <ChromePicker onChange={color => setTextColour(color.hex)} color={textColour}/>
                    </div> : null }
                </div>
            </div>
        )
    }
    return (
        <Spinner/>
    )
}