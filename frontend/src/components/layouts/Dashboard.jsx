import { Navbar } from "../Navbar";
import { Footer } from "../Footer";

export const Dashboard = ({children}) => {
    return (
        <>
            <Navbar />
            <main> {children} </main>
            <Footer />
        </>
    )
}