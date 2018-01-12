import React, { Component } from 'react';
import { Card } from 'react-materialize';
import Hotkeys from 'react-hot-keys';

class SoundPanel extends Component {
    constructor(props) {
        super(props);

        this.state = { playing: false }

        this.playStopSound = this.playStopSound.bind(this)
        window.eventEmitter.addListener("endSound", this.onEndSound.bind(this))
    }

    onEndSound(file) {
        if (file == this.props.file) {
            this.setState({ playing: false })
        }
    } 

    playStopSound() {
        if (this.state.playing) {
            window.panelController.cancel(this.props.file)
            this.setState({ playing: false })
        } else {
            window.panelController.play(this.props.file)
            this.setState({ playing: true })
        }   
    }

    render() {
        return <Card title={this.props.shortcut.toUpperCase()} actions={[<a href="#" onClick={this.playStopSound}>{this.state.playing ? "Stop" : "Play"}</a>]}>
            <p>{this.props.name}</p>
            <Hotkeys keyName={this.props.shortcut} onKeyUp={this.playStopSound}/>
        </Card>
    }
}

export default SoundPanel