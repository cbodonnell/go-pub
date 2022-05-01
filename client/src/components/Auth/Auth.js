import Modal from 'react-bootstrap/Modal';
import './Auth.scss';
import axios from 'axios';
import { formDataToURLString } from '../../utils/forms';
import logError from '../../utils/errors';
import HCaptcha from '@hcaptcha/react-hcaptcha';
import React from 'react';
import { withRouterAndRef } from '../../utils/router';
import { environment } from '../../environment';
import { AuthLoading } from '../Loading/AuthLoading';


class Auth extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            showLogin: false,
            showRegister: false,
            processing: false,
            registerError: '',
            loginError: '',
        }

        this.handleShowRegister = this.handleShowRegister.bind(this);
        this.handleCloseRegister = this.handleCloseRegister.bind(this);
        this.handleRegister = this.handleRegister.bind(this);
        this.handleShowLogin = this.handleShowLogin.bind(this);
        this.handleCloseLogin = this.handleCloseLogin.bind(this);
        this.handleLogin = this.handleLogin.bind(this);
        this.handleLogout = this.handleLogout.bind(this);
    }

    handleCloseLogin = (e) => this.setState({ showLogin: false });

    handleShowLogin(e) { 
        e.preventDefault();
        this.setState({
            showRegister: false,
            showLogin: true,
        });
    }

    handleCloseRegister = () => this.setState({ showRegister: false });

    handleShowRegister(e) {
        e.preventDefault();
        this.setState({
            showRegister: true,
            showLogin: false,
        });
    }

    handleLogin(e) {
        e.preventDefault();
        this.setState({
            processing: true,
            loginError: ''
        })
        const formData = new FormData(e.target);
        axios.post(environment.REACT_APP_AUTH_URL + '/login', formDataToURLString(formData), {
            headers: { 
                'Content-Type': 'application/x-www-form-urlencoded',
                'Accept': 'application/json'
            },
            withCredentials: true
        }).then(res => {
            const auth = res.data;
            this.props.onAuth(auth);
            this.setState({ showLogin: false })
        }).catch(error => {
            logError(error);
            this.setState({ loginError: 'Login failed' })
        }).finally(() => this.setState({ processing: false }));
    };

    handleRegister(e) {
        e.preventDefault();
        this.setState({
            processing: true,
            registerError: ''
        });
        const formData = new FormData(e.target);
        axios.post(environment.REACT_APP_AUTH_URL + '/register', formDataToURLString(formData), {
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            withCredentials: true
        }).then(res => {
            this.setState({ showRegister: false });
        }).catch(error => {
            logError(error);
            this.setState({ registerError: 'Registration failed' });
        }).finally(() => this.setState({ processing: false }));
    };

    handleLogout() {
        axios.get(environment.REACT_APP_AUTH_URL + '/logout', {
            withCredentials: true
        })
        .then(() => {
            // let history = useHistory();
            this.props.onAuth(null);
            this.props.history.push('/');
        })
        .catch(error => logError(error))
        .finally(() => this.props.onAuth(null));
    };

    render() {
        return (
            <>
                <div>
                    <div className="Auth-dialog">
                        {/* TODO: Break into component? */}
                        <Modal show={this.state.showLogin} onHide={this.handleCloseLogin} centered dialogClassName="Auth-dialog">
                            {this.state.processing && 
                                <AuthLoading />
                                // <div className="Auth-processing-overlay w-full h-full absolute flex justify-center items-center">
                                //     <div className="lds-grid"><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div></div>
                                // </div>
                            }
                            <Modal.Header closeButton>
                                <Modal.Title>Sign in</Modal.Title>
                            </Modal.Header>
                            <form onSubmit={this.handleLogin}>
                            <Modal.Body>
                            <div className="w-full h-full relative">
                                <div className="flex flex-col">
                                    {this.state.loginError && <p className="Auth-error">{this.state.loginError}</p>}
                                    <p className="Auth-label">Username:</p>
                                    <input className="Auth-input" type="text" name="username" id="username" required
                                    autoCorrect="off" autoCapitalize="off" />
                                    <p className="Auth-label">Password:</p>
                                    <input className="Auth-input" type="password" name="password" id="password" required />
                                </div>
                            </div>
                            </Modal.Body>
                            <Modal.Footer>
                            <div className="flex justify-between items-center w-full">
                                <div>
                                    { environment.REACT_APP_AUTH_REGISTER === "true" &&
                                        <p className="m-0">New user? <a href="/" className="Auth-register-link" onClick={this.handleShowRegister} >Sign up</a>.</p>
                                    }
                                </div>
                                <div>
                                    <input type="button" className="button button-secondary" value="Close" onClick={this.handleCloseLogin} />
                                    <input type="submit" className="button button-primary" value="Submit" />
                                </div>
                            </div>
                            </Modal.Footer>
                            </form>
                        </Modal>

                        {/* TODO: Break into component? */}
                        { environment.REACT_APP_AUTH_REGISTER === "true" &&
                            <Modal show={this.state.showRegister} onHide={this.handleCloseRegister} centered dialogClassName="Auth-dialog">
                                {this.state.processing && 
                                    <div className="Auth-processing-overlay w-full h-full absolute flex justify-center items-center">
                                        <div className="lds-grid"><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div></div>
                                    </div>
                                }
                                <Modal.Header closeButton>
                                    <Modal.Title>Sign up</Modal.Title>
                                </Modal.Header>
                                <form onSubmit={this.handleRegister}>
                                <Modal.Body>
                                    <div className="flex flex-col">
                                        {this.state.registerError && <p className="Auth-error">{this.state.registerError}</p>}
                                        <p className="Auth-label">Username:</p>
                                        <input className="Auth-input" type="text" name="username" id="username" required
                                        autoCorrect="off" autoCapitalize="off" />
                                        <p className="Auth-label">Password:</p>
                                        <input className="Auth-input" type="password" name="password" id="password" required />
                                        <p className="Auth-label">Confirm Password:</p>
                                        <input className="Auth-input" type="password" name="confirm-password" id="confirm-password" required />
                                        { environment.REACT_APP_AUTH_HCAPTCHA &&
                                            <div className="Auth-hCaptcha">
                                                <HCaptcha sitekey={environment.REACT_APP_AUTH_HCAPTCHA} />
                                            </div>
                                        }
                                    </div>
                                </Modal.Body>
                                <Modal.Footer>
                                <div className="flex justify-between items-center w-full">
                                    <p className="m-0">Existing user? <a href="/" className="Auth-register-link" onClick={this.handleShowLogin} >Sign in</a>.</p>
                                    <div>
                                        <input type="button" className="button button-secondary" value="Close" onClick={this.handleCloseRegister} />
                                        <input type="submit" className="button button-primary" value="Submit" />
                                    </div>
                                </div>
                                </Modal.Footer>
                                </form>
                            </Modal>
                        }
                    </div>
                </div>
            </>
        );
    }

}

export default withRouterAndRef(Auth);