import React, { Component } from 'react';
import { Card } from 'react-materialize';
import Hotkeys from 'react-hot-keys';
import Moment from 'react-moment';
import moment from 'moment-timezone';

Moment.globalMoment = moment;

class SoundPanel extends Component {
    interval = null

    constructor(props) {
        super(props);

        this.state = { playing: false, time: 0, interval: null  }

        this.playStopSound = this.playStopSound.bind(this)
        window.eventEmitter.addListener("endSound", this.onEndSound.bind(this))

        this.addTime = this.addTime.bind(this)
        this.resetTime = this.resetTime.bind(this)
    }

    addTime() {
        this.setState({ time: this.state.time + 1 })
    }

    resetTime() {
        console.log(this.state)
        this.setState({ time: 0 })
        clearInterval(this.interval)
    }

    onEndSound(file) {
        if (file == this.props.file) {
            this.setState({ playing: false })
            this.resetTime()
        }
    } 

    playStopSound() {
        if (this.state.playing) {
            window.panelController.cancel(this.props.file)
            this.setState({ playing: false })
            this.resetTime()
        } else {
            window.panelController.play(this.props.file)
            this.setState({ playing: true })
            this.interval = setInterval(this.addTime, 1000)
        }   
    }

    render() {
        let timer = <br/>
        if (this.state.playing) {
            timer = <p><Moment tz="UTC" format="HH:mm:ss"unix>{this.state.time}</Moment></p>
        }
        return <Card className={this.state.playing ? "deep-orange" : ""}title={this.props.shortcut.toUpperCase()} actions={[<a href="#" onClick={this.playStopSound}>{this.state.playing ? "Stop" : "Play"}</a>]}>
            <p>{this.props.name}</p>
            {timer}
            <Hotkeys keyName={this.props.shortcut} onKeyUp={this.playStopSound}/>
        </Card>
    }
}

export default SoundPanel