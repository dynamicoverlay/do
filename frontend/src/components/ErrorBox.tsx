import React from 'react';

export default (props: {children: any}) => {
    return (
        <div className="p-2 border-red-700 bg-red-500 text-white my-2 capitalize">
            {props.children}
        </div>
    )
}