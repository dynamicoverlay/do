import React, {Component, useState} from 'react';
import Navigation from '../Navigation';
import HomePage from '../../pages/dashboard/HomePage';
import LogoutPage from '../../pages/dashboard/LogoutPage';
import {useUserQuery} from '../../hooks/useUserQuery';

import {Route, useLocation} from 'react-router-dom';
import OverlaysPage from '../../pages/dashboard/OverlaysPage';
import OverlayControls from '../../pages/dashboard/overlay/OverlayControls';

export default ({match}: any) => {
    const {data} = useUserQuery();
    const [theme, setTheme] = useState('dark');
    let location = useLocation();
    return (
        <>
            <div className={`flex flex-col min-h-screen min-w-screen theme-${theme}`}>
               <Navigation userData={data.user} setTheme={setTheme} theme={theme} activeURL={location.pathname}></Navigation>
                <div className="bg-background flex-grow min-w-screen text-theme">
                    <Route path={match.url + "/"} exact={true} component={HomePage} />
                    <Route path={match.url + "/overlays"} exact={true}  component={OverlaysPage} />
                    <Route path={match.url + "/overlays/:id/control"} component={OverlayControls} />
                    <Route path={match.url + "/logout"} component={LogoutPage} />
                </div>
            </div>
        </>
    )
}


