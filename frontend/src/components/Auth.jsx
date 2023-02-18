import React, { useContext } from 'react';
import axios from 'axios';
import { Alert, Box, Button, TextField } from '@mui/material'
import { useState, useEffect } from 'react';
import { API_ROUTES, APP_ROUTES } from '../utils/constants';
import { Link, useNavigate } from 'react-router-dom';
import { useUser, storeTokenInLocalStorage, removeTokenFromLocaStorage } from '../lib/auth';
import { UserContext } from '../App';

export function SignOut() {
    const navigate = useNavigate();
    removeTokenFromLocaStorage();
    navigate(APP_ROUTES.SIGN_IN);
}

export function SignUp() {
    const navigate = useNavigate()
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState("");

    const signUp = async () => {
        if(!email) {
            setError("invalid email")
        }
        try {
            await axios({
                method: 'POST',
                url: API_ROUTES.SIGN_UP,
                data: {
                    email,
                    password,
                }
            });
            navigate(APP_ROUTES.SIGN_UP);
        }
        catch (err) {
            if(err.response.status === 401) {
                setError(err.response.data.message);
            } else if(err.response.status === 400) {
                setError(err.response.data.message)
            } else if(err.response.status === 201) {
                setError(err.response.data.message)
            }
        }
    };

    return (
        <Box 
            sx={{
                marginLeft: "10%",
                marginTop: "10%",
                maxWidth: "500px",
                '& .MuiTextField-root': { 
                    display: "flex",
                    marginTop: "20px",
                },

            }} 
        >
            {error && <Alert severity="error">{error}</Alert>}
            <TextField 
                variant="filled"
                label="Username"
                value={email}
                onChange={(e) => { setEmail(e.target.value); }}
            />
            <TextField
                variant="filled"
                label="Password"
                value={password}
                type="password"
                onChange={(e) => { setPassword(e.target.value); }}
            />
            <div style={{display: "inline-flex", marginTop: "20px"}}>
                <Button 
                    variant="contained"
                    onClick={signUp}
                >Register</Button>
                <p style={{marginLeft: "40px"}}>
                    Already a user?
                    <Link to={APP_ROUTES.SIGN_IN} style={{marginLeft: "10px"}}>
                        Sign In
                    </Link>
                </p>
            </div>
        </Box>
    );
}
export function SignIn() {
    const {user, setUser} = useContext(UserContext);
    const navigate = useNavigate();
    const [error, setError] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    useEffect((user) => {
        if(user) {
            navigate(APP_ROUTES.PROFILE)
        }
    }, [user]);

    const signIn = async () => {
        try {
            const response = await axios({
                method: 'post',
                url: API_ROUTES.SIGN_IN,
                data: {
                    email,
                    password
                }
            });
            storeTokenInLocalStorage(response.data.token);
            setUser()
            navigate(APP_ROUTES.DASHBOARD)
        }
        catch (err) {
            if(err?.response?.status) {
                setError(err.response.data.message);
            }
        }
    };


    return (
        <Box 
            sx={{
                marginLeft: "10%",
                marginTop: "10%",
                maxWidth: "500px",
                '& .MuiTextField-root': { 
                    display: "flex",
                    marginTop: "20px",
                },

            }} 
        >
        {error && <Alert severity="error">{error}</Alert>}
            <TextField 
                variant="filled"
                label="Username"
                value={email}
                onChange={(e) => { setEmail(e.target.value); }}
            />
            <TextField
                variant="filled"
                label="Password"
                value={password}
                type="password"
                onChange={(e) => { setPassword(e.target.value); }}
            />
            <div style={{display: "inline-flex", marginTop: "20px"}}>
                <Button 
                    variant="contained"
                    onClick={signIn}
                >Login</Button>
                <p style={{marginLeft: "40px"}}>
                    Don't have an account? 
                    <Link to={APP_ROUTES.SIGN_UP} style={{marginLeft: "10px"}}>
                        Sign Up
                    </Link>
                </p>
            </div>
        </Box>
    );
}
