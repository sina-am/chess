import './App.css';
import { APP_ROUTES } from './utils/constants';
import { SignIn, SignOut, SignUp } from './components/Auth';

import {
    BrowserRouter,
    Route, Routes
} from "react-router-dom";
import { ProfilePage } from './pages/ProfilePage';
import { PlayersPage } from './pages/PlayersPage';
import { GamePage } from './pages/GamePage';
import { Game } from './components/Game';
import { createContext, useContext, useEffect, useState } from 'react';
import { getAuthenticatedUser } from './lib/auth';



export const UserContext = createContext({
    user: null,
    getUser: async () => {},
});

function App() {
    const [user, setUser] = useState()


    useEffect(() => {
        async function getAuth() {
            if(!user) {
                setUser(await getAuthenticatedUser())
            }
        }
        getAuth()
    }, []);

    return (
        <UserContext.Provider value={{user: user, getUser: async () => {setUser(await getAuthenticatedUser())} }}>
            <BrowserRouter>
                <Routes>
                    <Route path={APP_ROUTES.SIGN_UP} element={<SignUp />} />
                    <Route path={APP_ROUTES.SIGN_IN} element={<SignIn />} />
                    <Route path={APP_ROUTES.SIGN_OUT} element={<SignOut />}></Route>
                    <Route path={APP_ROUTES.GAMES} element={<GamePage />}></Route>
                    <Route path={APP_ROUTES.ONLINE_GAME} element={<Game />}></Route>
                    <Route path={APP_ROUTES.PLAYERS} element={<PlayersPage />}></Route>
                    <Route path={APP_ROUTES.PROFILE} element={<ProfilePage />}></Route>
                    <Route path={APP_ROUTES.DASHBOARD} element={<PlayersPage />}></Route>
                </Routes>
            </BrowserRouter>
        </UserContext.Provider>
    );
}

export default App;
