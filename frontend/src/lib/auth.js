import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { API_ROUTES, APP_ROUTES } from '../utils/constants';
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
    try {
        const token = getTokenFromLocalStorage();
        if (!token) {
            return null;
        }
        const response = await axios({
            method: 'GET',
            url: API_ROUTES.GET_USER,
            headers: {
                Authorization: `${token}`
            }
        });
        if (response.data.id) {
            return response.data 
        }
        return null;
    }
    catch (err) {
        console.log('getAuthenticatedUser, Something Went Wrong', err);
        return null;
    }
}
