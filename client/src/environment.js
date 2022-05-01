function readEnv(variables) {
    const environment = {};
    variables.forEach(variable => {
        environment[variable] = process.env.NODE_ENV === 'production' ? window.env[variable] : process.env[variable];
    });
    return environment
}

export const environment = readEnv([
    'REACT_APP_AUTH_URL',
    'REACT_APP_AUTH_HCAPTCHA',
    'REACT_APP_AUTH_REGISTER',
    'REACT_APP_ACTIVITY_URL',
    'REACT_APP_PROXY_URL',
])