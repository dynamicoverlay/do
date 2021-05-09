import React, {Component} from 'react';


type TextInputProps = {
    type?: string,
    placeholder?: string,
    value?: string,
    onChange?: ((event: React.ChangeEvent<HTMLInputElement>) => void),
    error?: boolean,
    label?: string,
    max?: number,
    onClick?: ((event: any) => void)
}

export default class TextInput extends Component<TextInputProps, {}> {

    static defaultProps = {
        type: "text",
        placeholder: "",
        error: false
    }

    render(){
        let ourProps = {
            type: this.props.type,
            placeholder: this.props.placeholder,
            value: this.props.value,
        } 
        const input = <input {...ourProps} onChange={this.props.onChange} onClick={this.props.onClick} max={this.props.max} className={`p-2 my-2 text-theme shadow bg-background rounded border ${this.props.error ? 'border-red-700' : 'border-lightGray'}`}></input>;
        if(this.props.label) {
            return (
                <div className={"flex flex-col"}>
                    <p>{this.props.label}</p>
                    {input}
                </div>
            )
        } else {
            return input;
        }
    }

}