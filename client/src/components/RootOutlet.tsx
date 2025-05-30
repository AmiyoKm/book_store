
import { Outlet } from "react-router-dom";
import Navbar from "./Navbar";

const RootLayout = () => {
	return (
		<>
			<Navbar />
			<main className="min-h-screen">
				<Outlet />
			</main>
		</>
	);
};

export default RootLayout;
