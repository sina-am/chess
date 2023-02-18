import { Navbar } from "../Navbar";
import { Footer } from "../Footer";
import { UserContext } from "../../App"
import { useContext, useEffect } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import { APP_ROUTES } from "../../utils/constants";

export const Dashboard = ({children}) => {
    const {user, setUser} = useContext(UserContext);

    if (!user) {
        return 
    }

    return (
        <>
            <Navbar />
            <main> {children} </main>
            <Footer />
        </>
    )
}