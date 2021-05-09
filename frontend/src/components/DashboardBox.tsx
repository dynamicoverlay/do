import React, {useState} from 'react';


type DashboardBoxProps = {
    title: string;
    columns?: number;
    marginX?: string;
}
export default (props: DashboardBoxProps & {children: React.ReactNode}) => {
    const [showing, setShowing] = useState(true);
    return (
        <div className={`${props.marginX ? props.marginX : 'mx-4'} rounded flex flex-col font-light shadow-xl col-span-${props.columns ? props.columns : 1}`}>
            <div className={"w-full p-2 bg-boxTitle select-none cursor-pointer"} onClick={(() => setShowing(!showing))}>
                <h1 className={"text-lg uppercase"}>{props.title}</h1>
            </div>
            {showing && <div className={"bg-gray-1000"}>
                {props.children}
            </div>}
        </div>
    );
};
