import React, { Component } from "react"
import { Row, Col, Preloader } from "react-materialize"
import SoundPanel from "./soundpanel"

class Panel extends Component {
    constructor(props) {
        super(props)

        this.state = { items: [] }
        this.waitForPanelControllerIntetral = setInterval(this.checkPanelController.bind(this), 100)

        this.panelControllerDidApear = this.panelControllerDidApear.bind(this)
    }

    panelControllerDidApear() {
        window.panelController.getAllPanels().then((items) => this.setState({ items }))
    }

    checkPanelController() {
        if (window.panelController) {
            clearInterval(this.waitForPanelControllerIntetral)        
            this.panelControllerDidApear()
        }
    }

    render() {
        if (this.state.items.length === 0) {
            return <Row><Col s={4} offset='s6'><Preloader size='big' flashing={true} /></Col></Row>
        }

        let rows = []
        let c = 0
        let rc = 0
        for (let item of this.state.items) {
            if (!rows[rc]) {
                rows[rc] = []
            }
            rows[rc].push(item)
            c++
            if (c > 5) {
                c = 0
                rc++
            }
        }
        return <div>{rows.map((row) => <Row>{row.map((item) => <Col s={2}><SoundPanel {...item}/></Col>)}</Row>)}</div>
    }

}

export default Panel
