import React, {Component, useState} from 'react';
import DashboardBox from '../../components/DashboardBox';
import {useMutation, useQuery, useLazyQuery} from '@apollo/react-hooks';
import gql from 'graphql-tag';
import { Overlay } from '../../types';
import Button from '../../components/Button';
import TextInput from '../../components/form/TextInput';
import { useHistory } from 'react-router';
import { useInput } from '../../hooks/useInput';

const GET_OVERLAYS = gql`
    query GetOverlays{
        overlays {
            identifier
            name
            pin
        }
    }
`

const CREATE_OVERLAY = gql`
    mutation CreateOverlay($name: String!){
        createOverlay(name: $name) {
            identifier
            name
            pin
        }
    }
`

export default () => {
        const history = useHistory();
        const [creating, setCreating] = useState(false);
        const {data, refetch: refetchOverlays} = useQuery(GET_OVERLAYS, { errorPolicy: 'all' , onError: () => {}, onCompleted: (data) => {}});
        const {value:newNameValue, bind:newNameBind, reset:newNameReset} = useInput('');
        const [createOverlay, {data:createData, error:createError, loading:creatingOverlay}] = useMutation(CREATE_OVERLAY, { errorPolicy: 'all' , onError: () => {}, onCompleted: (data) => {
            setCreating(false)
            newNameReset()
            refetchOverlays()
        }});
        return (
            <div className={"flex flex-col py-2 mx-auto"} style={{maxWidth: "80%"}}>
                <div className="flex flex-row items-center">
                    <h1 className={"text-5xl"}>Overlays</h1>
                    <div className="ml-auto">
                        {creating && <TextInput {...newNameBind} type="text" placeholder="New Overlay"></TextInput>}
                        <Button color="green-700" className={"ml-auto"} onClick={() => {
                            if(creating){
                                if(newNameValue.length > 0){
                                    createOverlay({variables: {name: newNameValue}})
                                }
                            } else {
                                setCreating(true)
                            }
                        }}>
                            {creating ? 'Create' : 'Create New'}
                        </Button>
                    </div>
                </div>
                {data && data.overlays && data.overlays.map((overlay: Overlay) => {
                    return (
                        <div className={"flex flex-row bg-gray-1000 p-4 w-full my-2 items-center"}>
                            {overlay.name}
                            <div className="ml-auto">
                                <Button color="green-700" className={"mx-2"} onClick={() => {navigator.clipboard.writeText(`https://ovly.io/o?o=${overlay.identifier}&p=${overlay.pin}`)}}>
                                    Copy URL
                                </Button>
                                <Button color="blue-700" className={"mr-2"} onClick={() => history.push(`/dashboard/overlays/${overlay.identifier}/control`)}>Controls</Button>
                                <Button color="orange-700">Edit</Button>
                            </div>
                        </div>
                    )
                })}
            </div>
       )
}