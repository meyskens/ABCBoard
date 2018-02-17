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

        this.state = { playing: false, time: 0, interval: null, paused: false }

        this.playStopSound = this.playStopSound.bind(this)
        window.eventEmitter.addListener("endSound", this.onEndSound.bind(this))

        this.addTime = this.addTime.bind(this)
        this.resetTime = this.resetTime.bind(this)
        this.startTime = this.startTime.bind(this)
        this.pauseTime = this.pauseTime.bind(this)
        this.pauseToggle = this.pauseToggle.bind(this)
    }

    addTime() {
        this.setState({ time: this.state.time + 1 })
    }

    startTime() {
        this.interval = setInterval(this.addTime, 1000)
    }

    pauseTime() {
        clearInterval(this.interval)
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
            this.startTime()
        }   
    }

    pauseToggle() {
        if (!this.state.playing) {
            return
        }

        if (this.state.paused) {
            window.panelController.resume(this.props.file)
            this.setState({ paused: false })
            this.startTime()
        } else {
            window.panelController.pause(this.props.file)
            this.setState({ paused: true })
            this.pauseTime()
        } 
    }

    render() {
        let timer = <br/>
        if (this.state.playing) {
            timer = <p><Moment tz="UTC" format="HH:mm:ss"unix>{this.state.time}</Moment></p>
        }

        let buttons = [<a href="#" onClick={this.playStopSound}>{this.state.playing ? "Stop" : "Play"}</a>]
        if (this.state.playing) {
            buttons.push(<a href="#" onClick={this.pauseToggle}>{this.state.paused ? "Resume" : "Pause"}</a>)
        }

        return <Card className={this.state.playing ? "deep-orange" : ""}title={this.props.shortcut.toUpperCase()} actions={buttons}>
            <p>{this.props.name}</p>
            {timer}
            <Hotkeys keyName={this.props.shortcut} onKeyUp={this.playStopSound}/>
            <Hotkeys keyName={`shift+${this.props.shortcut}`} onKeyUp={this.pauseToggle}/>
        </Card>
    }
}

export default SoundPanel