import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import api from "@/config/axios";

const ActivateAccountPage = () => {
	const { token } = useParams();
	const navigate = useNavigate();
	const [message, setMessage] = useState("");
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		const activateAccount = async () => {
			try {
				const response = await api.put(`/authentication/activate/${token}`);
				if (response.status === 200) {
					setMessage("Your account has been successfully activated!");
					setTimeout(() => navigate("/sign-in"), 3000);
				}
			} catch (err) {
				setMessage("Activation link is invalid or expired.");
			} finally {
				setLoading(false);
			}
		};
		activateAccount();
	}, [token, navigate]);

	return (
		<div className="min-h-screen flex items-center justify-center bg-gray-100">
			{loading ? (
				<p className="text-lg">Activating your account...</p>
			) : (
				<div className="bg-white p-6 rounded-lg shadow-md">
					<p className="text-lg">{message}</p>
				</div>
			)}
		</div>
	);
};

export default ActivateAccountPage;
