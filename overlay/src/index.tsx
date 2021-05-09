import * as React from 'react';
import * as ReactDOM from 'react-dom';

import App from './App';
import './styles.css';

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
const overlayID: string = urlParams.get("o") as string;
const overlayPin: string = urlParams.get("p") as string;
console.log("ID:", overlayID);
console.log("Pin:", overlayPin);

var mountNode = document.getElementById("app");
ReactDOM.render(<App overlayID={overlayID} overlayPin={overlayPin}/>, mountNode);
