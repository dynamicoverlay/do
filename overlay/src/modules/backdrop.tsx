import * as React from 'react';

import {motion} from 'framer-motion';
import {BackdropProps} from '../types';

export default (props: BackdropProps) => {
    console.log("Transition time is ", props.state?.transitionTime);
    const transition = {
        type: 'spring', 
        duration: props.state?.transitionTime || 5
    }
    console.log("Transition is", transition);
    const backgroundStyle = {
        backgroundColor: props.state?.color || '#2980b9',
        backgroundImage: `url(${props.state?.backgroundImage})` || '',
        backgroundSize: 'cover',
        backgroundPosition: 'center'
    }
    const textStyle = {
        color: props.state?.textColor || '#ffffff'
    }
    return (
        <motion.div initial={{top:"-100%"}} animate={{top: props.state?.visible ? 0 : '-100%'}} transition={transition} className="bg-blue-500 absolute top-0 left-0 z-10 w-full h-full flex flex-row items-center" style={backgroundStyle}>
            <h1 style={textStyle} className="mx-auto text-6xl">{props.state?.text}</h1>
        </motion.div>
    )
};
