import React from 'react';
import './ToggleSwitch.scss';

class ToggleButton extends React.Component {
    constructor(props) {
        super(props);
        this.state = { toggled: false };
        this.triggerToggle = this.triggerToggle.bind(this);
    }

    componentDidMount() {
        const toggled = localStorage.getItem('light-mode') === 'true'
        this.setState({toggled});
        this.props.onToggle(toggled);
    }

    triggerToggle(e) {
        e.preventDefault();
        if (this.props.disabled) { return; }
        const toggled = !this.state.toggled;
        this.setState({toggled});
        this.props.onToggle(toggled);
        localStorage.setItem('light-mode', toggled.toString())
    }

    render() {
        return (
            <button onClick={this.triggerToggle} 
            className={
                `ToggleSwitch wrg-toggle 
                ${this.state.toggled ? 'wrg-toggle--checked' : ''}
                ${this.props.disabled ? 'wrg-toggle--disabled' : ''}`
            }>
                <div className="relative">
                    <div className="wrg-toggle-container">
                        <div className="wrg-toggle-check">
                            <span><>🌞</></span>
                        </div>
                        <div className="wrg-toggle-uncheck">
                            <span><>🌜</></span>
                        </div>
                    </div>
                    <div className="wrg-toggle-circle"></div>
                    <input className="wrg-toggle-input" type="checkbox" aria-label="Toggle Button"
                    tabIndex="-1" />
                </div>
            </button>
        )
      }
}

export default ToggleButton;