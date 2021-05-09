import React, { Component } from 'react';
import Button from '../Button';


type LoginButtonsProps = {
    goToEmail: Function
}

export default class LoginButtons extends Component<LoginButtonsProps, {}> {

    render() {
        return (
            <>
                <Button color="twitch">TWITCH</Button>
                <Button color="gray-700" onClick={this.props.goToEmail}>EMAIL</Button>
            </>
        )
    }

}