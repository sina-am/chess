import './App.css';
import { Game } from './components/Game'
import { Profile } from './components/Profile';
import { APP_ROUTES } from './utils/constants';
import { SignIn, SignOut, SignUp } from './components/Auth';
import { Players } from './components/Player';

import {
    BrowserRouter,
    Route, Routes
} from "react-router-dom";
import { Navbar } from './components/Navbar';


function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path={APP_ROUTES.SIGN_UP} element={<SignUp />} />
                <Route path={APP_ROUTES.SIGN_IN} element={<SignIn />} />
                <Route path={APP_ROUTES.SIGN_OUT} element={<SignOut />}></Route>
                <Route path={APP_ROUTES.GAMES} element={<Game />}></Route>
                <Route path={APP_ROUTES.PLAYERS} element={<Players />}></Route>
                <Route path={APP_ROUTES.PROFILE} element={<Profile />}></Route>
                <Route path={APP_ROUTES.DASHBOARD} element={<Navbar />}></Route>
            </Routes>
        </BrowserRouter>
    );
}

export default App;
