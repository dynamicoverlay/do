import React from 'react';
import Button from '../components/Button';
import {Link} from 'react-router-dom';


export default () => {
    return (
        <div className="flex flex-col justify-center text-theme items-center min-h-screen min-w-screen">
            <div className="mx-auto w-3/4 lg:w-1/2 text-center">
                <h1 className="text-6xl font-semibold">My<span className="font-light">Overlay</span></h1>
                <div className="p-4 mt-4 w-1/2 lg:w-1/3 mx-auto rounded flex flex-row bg-gray-1000 font-light shadow-xl text-white">
                    <div className="flex flex-col justify-center w-full">
                        <Link to="/login" ><Button expand={true} color="green-400">Login</Button></Link>
                        <Link to="/signup"  ><Button expand={true} color="green-400">Signup</Button></Link>
                    </div>
                </div>
            </div>
        </div>
    )
}

