import React, { Component } from 'react';
import { Row, Col, Preloader } from 'react-materialize';
import { EventEmitter } from 'fbemitter';
import { getAllPanels } from '../apis/panels_api'
import SoundPanel from './soundpanel'

class Panel extends Component {
    constructor(props) {
        super(props)

        this.state = { items: [] }
        this.waitForPanelControllerIntetral = setInterval(this.checkPanelController.bind(this), 100)

        this.panelControllerDidApear = this.panelControllerDidApear.bind(this)
    }

    panelControllerDidApear() {
        window.eventEmitter = new EventEmitter(); // sorry you have to see this but we're fighting webkit here
        getAllPanels().then((response) => this.setState({ items: response.data }))
    }

    checkPanelController() {
        if (window.panelController) {
            clearInterval(this.waitForPanelControllerIntetral);        
            this.panelControllerDidApear();
        }
    }

    render() {
        if (this.state.items.length === 0) {
            return <Row><Col s={4} offset='s6'><Preloader size='big' flashing={true} /></Col></Row>
        }

        console.log(this.state.items)
        return <Row>{this.state.items.map((item) => <Col s={2}><SoundPanel {...item}/></Col>)}</Row>
    }

}

export default Panel