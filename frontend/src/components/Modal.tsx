import React from 'react';

interface ModalProps {
    active: boolean;
    children?: any;
}

export default (props: ModalProps) => {
    if(!props.active){
        return null;
    }
    return (
        <div className={"fixed w-full h-full top-0 left-0 flex items-center justify-center"}>
            <div className={"absolute w-full h-full bg-gray-900 opacity-50"}></div>
            <div className={"bg-white w-11/12 md:max-w-md "}>

            </div>
        </div>
    )
}