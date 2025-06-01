import { createBrowserRouter } from "react-router";
import LandingPage from "./pages/landing";
import SignupPage from "./pages/sign-up";
import LoginPage from "./pages/sign-in";
import ActivateAccountPage from "./pages/activate_user";
import ResetPassword from "./pages/reset_password";
import ForgotPasswordPage from "./pages/forgot_password";
import RootLayout from "./components/RootOutlet";
import BooksPage from "./pages/books";
import BookDetailPage from "./pages/books/BookDetail";
import CartPage from "./pages/cart";
const CategoriesPage = () => <div>Books Categories Page</div>;
const WishlistPage = () => <div>Wishlist Page</div>;
const OrdersPage = () => <div>Orders List Page</div>;
const OrderDetailPage = () => <div>Order Detail Page</div>;
const ProfilePage = () => <div>User Profile Page</div>;
const CheckoutPage = () => <div>Checkout Page</div>;
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
		element: <RootLayout />,
		children: [
			{
				path: "/books",
				element: <BooksPage />,
			},
			{
				path: "/books/:bookId",
				element: <BookDetailPage/>,
			},
			{
				path : "/categories",
				element : <CategoriesPage />
			},
			{
				path: "/cart",
				element: <CartPage />,
			},
			{
				path: "/wishlist",
				element: <WishlistPage />,
			},
			{
				path: "/orders",
				element: <OrdersPage />,
			},
			{
				path: "/orders/:orderId",
				element: <OrderDetailPage />,
			},

			{
				path: "/profile",
				element: <ProfilePage />,
			},
			{
				path: "/checkout",
				element: <CheckoutPage />,
			},
		],
	},
]);
