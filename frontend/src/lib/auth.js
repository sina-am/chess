import { useState, useEffect } from 'react';
import { APP_ROUTES } from '../utils/constants';
import { useNavigate } from 'react-router-dom';
import { API_ROUTES } from '../utils/constants';
import axios from 'axios';

export function storeTokenInLocalStorage(token) {
    localStorage.setItem('token', token);
}

export function getTokenFromLocalStorage() {
    return localStorage.getItem('token');
}

export function removeTokenFromLocaStorage() {
    return localStorage.removeItem('token');
}

export async function getAuthenticatedUser() {
    const defaultReturnObject = { authenticated: false, user: null };
    try {
        const token = getTokenFromLocalStorage();
        if (!token) {
            return defaultReturnObject;
        }
        const response = await axios({
            method: 'GET',
            url: API_ROUTES.GET_USER,
            headers: {
                Authorization: `${token}`
            }
        });
        if (response.data.id) {
            return { authenticated: true, user: response.data }
        }
        return defaultReturnObject;
    }
    catch (err) {
        console.log('getAuthenticatedUser, Something Went Wrong', err);
        return defaultReturnObject;
    }
}
export function useUser() {
    const [user, setUser] = useState(null);
    const [authenticated, setAutenticated] = useState(false);
    const navigate = useNavigate();

    useEffect((navigate) => {
        async function getUserDetails() {
            const { authenticated, user } = await getAuthenticatedUser();
            if (!authenticated) {
                navigate(APP_ROUTES.SIGN_IN);
                return;
            }
            setUser(user);
            setAutenticated(authenticated);
        }
        getUserDetails();
    }, [navigate]);

    return [user, authenticated];
}