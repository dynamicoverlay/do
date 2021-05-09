import React, {Component, useState} from 'react';
import {Link, useHistory} from 'react-router-dom';


type UserData = {
    email: string;
    role: string;
    username: string;
}

type NavigationProps = {
    userData: UserData,
    setTheme: Function,
    theme: string,
    activeURL: string
}

export default (props: NavigationProps) => {
    const [droppedDown, setDroppedDown] = useState(false);
    const history = useHistory();

    function isActivePage(link: string){
        return props.activeURL === link;
    }

    return (
        <>
           <div className="flex flex-row items-center bg-gray-1000 min-w-full shadow text-theme min-h-16">
                <ul className="flex flex-row min-w-full">
                    <div className="flex flex-row">
                        <li className="mr-4 font-bold p-5">DynamicOverlay</li>
                        <li className={`${isActivePage("/dashboard") ? 'bg-background' : 'hover:bg-background'} p-5 cursor-pointer`} onClick={() => history.push('/dashboard')}>Dashboard</li>
                        <li className={`${isActivePage("/dashboard/overlays") ? 'bg-background' : 'hover:bg-background'} p-5 cursor-pointer`}  onClick={() => history.push('/dashboard/overlays')}>Overlays</li>
                        <li className="hover:bg-background p-5 cursor-pointer">Chat</li>
                        <li className="hover:bg-background p-5 cursor-pointer">Settings</li>
                    </div>
                    <div className="flex flex-row ml-auto relative">
                        <li className="hover:bg-background p-5 cursor-pointer" onClick={() => setDroppedDown(!droppedDown)}>{props.userData.username}</li>
                        {droppedDown && (<ul className="absolute bg-gray-1000 w-full" style={{top: '4rem'}}>
                            <li onClick={() => props.setTheme(props.theme === 'light' ? 'dark' : 'light')} className="hover:bg-background p-5 block cursor-pointer whitespace-no-wrap">
                                {props.theme === 'light' ? 'Dark mode' : 'Light mode'}
                            </li>
                            <li>
                                <Link to="/dashboard/logout" ><div className="hover:bg-background p-5 block cursor-pointer whitespace-no-wrap">Logout</div></Link>
                            </li>
                        </ul>)}
                    </div>
                </ul>
           </div>
        </>
    )

};


