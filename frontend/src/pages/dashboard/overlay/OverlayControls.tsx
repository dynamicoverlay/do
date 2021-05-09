import React, {Component, useState, useEffect, useRef} from 'react';
import {useMutation, useQuery, useLazyQuery} from '@apollo/react-hooks';
import gql from 'graphql-tag';
import Button from '../../../components/Button';
import { useParams } from 'react-router';
import Spinner from '../../../components/Spinner';
import DashboardBox from '../../../components/DashboardBox';
import { useOverlayState } from '../../../hooks/useOverlayState';
import BackdropModule from './modules/BackdropModule';


const GET_OVERLAY = gql`
    query GetOverlay($id: String!){
        overlay(id: $id) {
            identifier
            name
            pin
            modules {
                module {
                    name
                    identifier
                }
                enabled
                settings
            }
        }
        modules {
            name
            identifier
        }
    }
`


const UPDATE_STATE = gql`
    mutation UpdateState($overlay: String!, $key: String!, $value: String!){
        updateState(overlay: $overlay, key: $key, value: $value){
            updated
        }
    }
`

export default () => {
        let { id } = useParams();
        let pin = "";
        const {data, loading} = useQuery(GET_OVERLAY, { variables: {id: id}, errorPolicy: 'all' , onError: () => {}, onCompleted: (data) => {}});
        if(!loading && data && data.overlay){
            pin = data.overlay.pin;
        }
        function getOverlayURL(){
            return `http://localhost:1234/?o=${data.overlay.identifier}&p=${data.overlay.pin}&preview=true`
        }
        const ws = useRef<WebSocket | null>(null);
        const overlayState = useOverlayState(ws, id, pin);
        const [doUpdate, {error:loginError, loading:loggingIn}] = useMutation(UPDATE_STATE, { errorPolicy: 'all' , onError: () => {}, onCompleted: (data) => {}});
        function getModuleState<T>(name: string): T {
            return overlayState[name] as T;
        }
        let overlayPage;
        if(!loading && data && data.overlay){
            overlayPage = (
                <>
                <div className="flex flex-col flex-1">
                    <div className="flex flex-row items-end">
                        <h1 className={"text-5xl"}>{data.overlay.name}</h1>
                        <div className="ml-auto">
                            <Button color="green-600" className="mx-2">Modules</Button>
                            <Button color="orange-600" className="mx-2">Edit</Button>
                        </div>
                    </div>
                    <div className="flex flex-row">
                        <div className="flex flex-col w-1/3">
                            <DashboardBox title="Backdrop" marginX={"mr-4"}>
                                <BackdropModule state={getModuleState('backdrop')} updateState={(key: string, value: string) => doUpdate({variables: {overlay: data.overlay.identifier, key, value}})}></BackdropModule>
                            </DashboardBox>
                        </div>
                        <div className="flex-1 border-lightGray border" >
                            <div style={{paddingBottom: '51%'}} className="relative">
                                <iframe src={getOverlayURL()} className={"w-full h-full absolute"}></iframe>
                            </div>
                        </div>
                    </div>
                </div>
                </>
            )
        }
        return (
            <div className={"flex flex-col py-2 mx-auto"} style={{maxWidth: "80%"}}>
                <div className="flex flex-row items-center">
                    {loading && <Spinner />}
                    {!loading && overlayPage}
                </div>
            </div>
       )
}