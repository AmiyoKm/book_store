import { passwordCheckValidation } from "@/config/api/auth";
import { useQuery } from "@tanstack/react-query";
import { useLocation } from "react-router";
import {
	Card,
	CardContent,
} from "@/components/ui/card";
import ResetPasswordForm from "./ResetPasswordForm";

const ResetPassword = () => {
	const { search } = useLocation();
	const params = new URLSearchParams(search);
	const token = params.get("token");

	if (!token) {
		return (
			<div className="flex min-h-screen items-center justify-center">
				<Card className="w-full max-w-md shadow-lg">
					<CardContent className="py-8 text-center">
						<span className="text-red-500 text-lg">
							No token provided. Please request a new password reset.
						</span>
					</CardContent>
				</Card>
			</div>
		);
	}

	return <ResetPasswordTokenValidator token={token} />;
};

const ResetPasswordTokenValidator = ({ token }: { token: string }) => {
    const query = useQuery({
		queryKey: ["password", "token", token],
		queryFn: () => passwordCheckValidation(token!),
		enabled: !!token,
	});

	if (query.isPending) {
		return (
			<div className="flex min-h-screen items-center justify-center">
				<Card className="w-full max-w-md shadow-lg">
					<CardContent className="py-8 text-center">
						<span className="text-lg">Validating token...</span>
					</CardContent>
				</Card>
			</div>
		);
	}

	if (query.isError) {
		return (
			<div className="flex min-h-screen items-center justify-center">
				<Card className="w-full max-w-md shadow-lg">
					<CardContent className="py-8 text-center">
						<span className="text-red-500 text-lg">
							Invalid or expired token. Please request a new password reset.
						</span>
					</CardContent>
				</Card>
			</div>
		);
	}

	if (query.isSuccess) {
        return <ResetPasswordForm token={token} userId={query.data?.data.user_id} />
    }

    return null;
}

export default ResetPassword;