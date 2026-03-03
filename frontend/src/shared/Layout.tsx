import { Outlet } from "react-router-dom";

function Layout() {
    return (
        <div>
            <header style={{ padding: "1rem", borderBottom: "1px solid #ccc" }}>
                EVNTX
            </header>
            <main style={{ padding: "1rem" }}>
                <Outlet />
            </main>
        </div>
    );
}

export default Layout;