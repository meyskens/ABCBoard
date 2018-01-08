import React, { Component } from 'react';
import { Card } from 'react-materialize';
import Hotkeys from 'react-hot-keys';

class SoundPanel extends Component {
    constructor(props) {
        super(props);
        console.log(props)

        this.playSound = this.playSound.bind(this)
    }

    playSound() {
        window.panelController.play(this.props.file)
    }

    render() {
        return <Card title={this.props.shortcut.toUpperCase()} actions={[<a href="#" onClick={this.playSound}>Play</a>]}>
            <p>{this.props.name}</p>
            <Hotkeys keyName={this.props.shortcut} onKeyUp={this.playSound}/>
        </Card>
    }
}

export default SoundPanel