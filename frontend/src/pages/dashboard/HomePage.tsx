import React, {Component, useState} from 'react';
import DashboardBox from '../../components/DashboardBox';
import {useMutation} from '@apollo/react-hooks';
import gql from 'graphql-tag';




export default () => {
        return (
            <div className={"grid grid-flow-row items-start md:grid-cols-4 gap-4 mt-2"}>
                <DashboardBox title={"Hello World"} columns={2}>
                    <div className={"p-2 "}>
                       
                        <h2>HELLOOOOO WORLDLDDD</h2>
                    </div>
                </DashboardBox> 
                <DashboardBox title={"Hello World"}>
                    <div className={"p-2"}>
                        <h2>HELLOOOOO WORLDLDDD</h2>
                    </div>
                </DashboardBox>
                <DashboardBox title={"Hello World"}>
                    <div className={"p-2"}>
                        <h2>HELLOOOOO WORLDLDDD</h2>
                    </div>
                </DashboardBox>
                <DashboardBox title={"Hello World"}>
                    <div className={"p-2 "}>
                        <h2>HELLOOOOO WORLDLDDD</h2>
                    </div>
                </DashboardBox> 
                <DashboardBox title={"Hello World"}>
                    <div className={"p-2"}>
                        <h2>HELLOOOOO WORLDLDDD</h2>
                    </div>
                </DashboardBox>
                <DashboardBox title={"Hello World"}>
                    <div className={"p-2"}>
                        <h2>HELLOOOOO WORLDLDDD</h2>
                    </div>
                </DashboardBox> 
            </div>
       )
}