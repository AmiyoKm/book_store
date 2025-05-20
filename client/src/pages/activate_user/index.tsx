import { useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { activateUser } from "@/config/api/auth";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
	CardDescription,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";

const ActivateAccountPage = () => {
	const { token } = useParams();
	const navigate = useNavigate();

	const query = useQuery({
		queryKey: ["token", token],
		queryFn: () => activateUser(token!),
		enabled: !!token,
	});

	useEffect(() => {
		if (query.isSuccess) {
			const timer = setTimeout(() => {
				navigate("/sign-in");
			}, 2000);
			return () => clearTimeout(timer);
		}
	}, [query.isSuccess, navigate]);

	return (
		<div className="flex min-h-screen items-center justify-center bg-muted">
			<Card className="w-full max-w-md shadow-lg">
				<CardHeader>
					<CardTitle className="text-2xl text-center">
						Activate Account
					</CardTitle>
					<CardDescription className="text-center">
						Activate your BookBound account to get started.
					</CardDescription>
				</CardHeader>
				<CardContent>
					{query.isLoading && (
						<div className="flex flex-col items-center py-8">
							<div className="animate-spin rounded-full h-10 w-10 border-b-2 border-primary-600 mb-4"></div>
							<p className="text-lg font-medium text-gray-700">
								Activating your account...
							</p>
						</div>
					)}
					{query.isError && (
						<div className="flex flex-col items-center py-8">
							<div className="text-red-500 text-3xl mb-2">✖</div>
							<p className="text-lg font-semibold mb-2">Activation Failed</p>
							<p className="text-gray-500 mb-4">
								{(query.error as Error).message || "Invalid or expired token."}
							</p>
							<Button variant="outline" onClick={() => navigate("/sign-in")}>
								Go to Sign In
							</Button>
						</div>
					)}
					{query.isSuccess && (
						<div className="flex flex-col items-center py-8">
							<div className="text-green-500 text-3xl mb-2">✔</div>
							<p className="text-lg font-semibold mb-2">Account Activated!</p>
							<p className="text-gray-500">Redirecting to sign in...</p>
						</div>
					)}
				</CardContent>
			</Card>
		</div>
	);
};

export default ActivateAccountPage;
