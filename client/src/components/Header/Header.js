import React from 'react';
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';
import NavDropdown from 'react-bootstrap/NavDropdown';
import ToggleSwitch from '../ToggleSwitch/ToggleSwitch';
import './Header.scss';
import Auth from '../Auth/Auth';
import axios from 'axios';
import logError from '../../utils/errors';
import { AiOutlineMenu } from "react-icons/ai";
import { Link } from 'react-router-dom';
import { environment } from '../../environment';
import Create from '../Create/Create';

export default class Header extends React.Component {

    constructor(props) {
        super(props);
        this.state = { auth: null, expanded: false };
        this.authRef = React.createRef();
        this.handleAuth = this.handleAuth.bind(this);
        this.handleShowLogin = this.handleShowLogin.bind(this);
        this.handleSelect = this.handleSelect.bind(this);
        this.handleClickBrand = this.handleClickBrand.bind(this);
        this.handleToggle = this.handleToggle.bind(this);
        this.handleLogout = this.handleLogout.bind(this);
        this.handleLightModeToggle = this.handleLightModeToggle.bind(this);
    }

    componentDidMount() {
        this.checkAuth();
    }

    checkAuth() {
        // TODO: Async loading here too
        axios.get(environment.REACT_APP_AUTH_URL + '/', {
            headers: { 'Accept': 'application/json' },
            withCredentials: true
        }).then(res => {
            const auth = res.data;
            this.handleAuth(auth);
        }).catch(error => {
            logError(error);
            this.handleAuth(null);
        });
    }

    handleAuth(auth) {
        this.setState({auth});
        this.handleSelect();
        this.props.onAuth(auth);
    }

    handleShowLogin(e) {
        this.authRef.current.handleShowLogin(e);
    }

    handleSelect() {
        // e.preventDefault();
        this.setState({ expanded: false });
    }

    handleClickBrand() {
        // e.preventDefault();
        this.setState({ expanded: false });
        const content = document.querySelector('#content');
        content.scroll({ top: 0, behavior: 'smooth' });
    }

    handleToggle(e) {
        e.preventDefault();
        this.setState({ expanded: !this.state.expanded });
    }

    handleLogout(e) {
        e.preventDefault();
        this.authRef.current.handleLogout();
    }

    handleLightModeToggle(toggled) {
        if (toggled) {
            document.body.classList.add('inverted');
        } else {
            document.body.classList.remove('inverted');
        }
    }
 
    render() {
      return (
        <div className="Header w-full flex-initial">
            <Navbar collapseOnSelect expand="md" expanded={this.state.expanded} variant="dark"
            className="fixed w-full Header-Navbar" onSelect={this.handleSelect}>
                <Link to="/">
                    <Navbar.Brand className="Header-brand" onClick={this.handleClickBrand}>
                        Studio 10B
                    </Navbar.Brand>
                </Link>
                <button type="button" className="Header-navbar-toggle" onClick={this.handleToggle}>
                    <AiOutlineMenu />
                </button>
                {/* <Navbar.Toggle aria-controls="responsive-navbar-nav" onClick={this.handleToggle} /> */}
                <Navbar.Collapse id="responsive-navbar-nav">
                    <Nav>
                        <ToggleSwitch onToggle={this.handleLightModeToggle} />
                    </Nav>
                    <Nav className="ml-auto">
                        {this.state.auth &&
                        <div className="Header-button-bar">
                            <Create />
                        </div>
                        }
                        {/* <ToggleSwitch onToggle={this.handleLightModeToggle} /> */}
                        <Auth ref={this.authRef} onAuth={this.handleAuth} />
                        {!this.state.auth &&
                        <>
                            <Nav>
                                <button className="Header-Auth-button" type="button" onClick={this.handleShowLogin}>Sign in</button>
                            </Nav>
                        </>
                        }
                    </Nav>
                    {this.state.auth &&
                    <Nav>
                        <NavDropdown title={this.state.auth.username} id="basic-nav-dropdown" className="Header-NavDropdown">
                            <Link to="/" className="dropdown-item">
                                Home
                            </Link>
                            <Link to="/" className="dropdown-item" onClick={this.handleLogout}>
                                Sign out
                            </Link>
                        </NavDropdown>
                    </Nav>
                    }
                </Navbar.Collapse>
            </Navbar>
            <div className="Header-spacer"></div>
        </div>
      )
    }
  }