import { createBrowserRouter } from "react-router";
import LandingPage from "./pages/landing";
import SignupPage from "./pages/sign-up";
import LoginPage from "./pages/sign-in";
import ActivateAccountPage from "./pages/activate_user";
import ResetPassword from "./pages/reset_password";
import ForgotPasswordPage from "./pages/forgot_password";

export const router = createBrowserRouter([
	{
		path: "/",
		element: <LandingPage />,
	},
	{
		path: "/sign-up",
		element: <SignupPage />,
	},
	{
		path: "/sign-in",
		element: <LoginPage />,
	},
	{
		path: "/confirm/:token",
		element: <ActivateAccountPage />,
	},
	{
		path: "/reset-password",
		element: <ResetPassword />,
	},
	{
		path: "/forgot-password",
		element: <ForgotPasswordPage />,
	},
	{
		path: "/home",
		element: <div>HOME</div>,
	},
]);
