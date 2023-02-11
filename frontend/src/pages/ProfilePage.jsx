import { Dashboard } from "../components/layouts/Dashboard";
import { Profile } from "../components/Profile";

export function ProfilePage() {
    return (
        <>
            <Dashboard>
                <Profile />
            </Dashboard>
        </>
    )
}