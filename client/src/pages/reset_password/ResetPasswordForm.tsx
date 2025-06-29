import { useState } from "react";
import { passwordChange } from "@/config/api/auth";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router";
import type { ChangePasswordPayload } from "@/types/auth";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
	CardDescription,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";

interface ResetPasswordFormProps {
    token: string;
    userId: number;
}

const ResetPasswordForm = ({ token, userId }: ResetPasswordFormProps) => {
    const navigate = useNavigate();
	const mutation = useMutation({
		mutationFn: passwordChange,
	});

	const [password, setPassword] = useState("");
	const [confirmPassword, setConfirmPassword] = useState("");
	const [error, setError] = useState("");

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		setError("");
		if (password !== confirmPassword) {
			setError("Passwords do not match.");
			return;
		}
		const payload: ChangePasswordPayload = {
			token: token!,
			new_password: password,
			user_id: userId,
		};
		mutation.mutate(payload);
	};

    return (
        <div className="flex min-h-screen items-center justify-center bg-muted">
            <Card className="w-full max-w-md shadow-lg">
                <CardHeader>
                    <CardTitle className="text-2xl text-center">
                        Reset Password
                    </CardTitle>
                    <CardDescription className="text-center">
                        Enter your new password below.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    {error && <div className="text-red-500 text-sm mb-2">{error}</div>}
                    {mutation.isSuccess ? (
                        <div className="flex flex-col items-center py-8">
                            <div className="text-green-600 text-3xl mb-2">✔</div>
                            <p className="text-lg font-semibold mb-2">
                                Password reset successful!
                            </p>
                            <p className="text-gray-500 mb-4">
                                You can now sign in with your new password.
                            </p>
                            <Button
                                type="button"
                                className="w-full"
                                onClick={() => navigate("/sign-in")}
                            >
                                Go to Sign In
                            </Button>
                        </div>
                    ) : (
                        <form onSubmit={handleSubmit} className="space-y-6">
                            <div className="space-y-2">
                                <Label htmlFor="password">New Password</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="confirmPassword">Confirm Password</Label>
                                <Input
                                    id="confirmPassword"
                                    type="password"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    required
                                />
                            </div>
                            <Button
                                type="submit"
                                className="w-full"
                                disabled={mutation.isPending}
                            >
                                {mutation.isPending ? "Resetting..." : "Reset Password"}
                            </Button>
                        </form>
                    )}
                </CardContent>
            </Card>
        </div>
    );
};

export default ResetPasswordForm;