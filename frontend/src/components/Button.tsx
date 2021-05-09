import React, {Component} from 'react';


type ButtonProps = {
    color?: string,
    onClick?: Function,
    expand?: boolean,
    type?: "button" | "submit" | "reset" | undefined,
    className?: string
}

export default class Button extends Component<ButtonProps, {}> {

    render(){
        return (
            <>
               <button type={this.props.type ? this.props.type : 'button' } className={`p-2 bg-${this.props.color ? this.props.color : 'grey-700'} my-2 shadow hover:opacity-75 transition-opacity duration-300 ease-in-out ${this.props.expand ? 'w-full' : ''} ${this.props.className}`} onClick={this.props.onClick?.bind(this)}>{this.props.children}</button>
            </>
        )
    }

}
