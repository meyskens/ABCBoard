import React from "react"
import ReactDOM from "react-dom"

import Header from "./components/header"
import Panel from "./components/panel"

const App = function () {
    return (
        <div>
            <Header/>
            <Panel/>
        </div>
    )
}

ReactDOM.render(
    <App />,
    document.querySelector("#container")
)