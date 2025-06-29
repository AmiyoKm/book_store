import React, { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { passwordResetRequest } from "@/config/api/auth";
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

const ForgotPasswordPage = () => {
	const [email, setEmail] = useState("");
	const [error, setError] = useState("");

	const mutation = useMutation({
		mutationFn: passwordResetRequest,
		onError: (err: { response?: { data?: { message?: string } } }) => {
			setError(err?.response?.data?.message || "Something went wrong.");
		},
	});

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		setError("");
		if (!email) {
			setError("Email is required.");
			return;
		}
		mutation.mutate({ email });
	};

	return (
		<div className="flex min-h-screen items-center justify-center bg-muted">
			<Card className="w-full max-w-md shadow-lg">
				<CardHeader>
					<CardTitle className="text-2xl text-center">
						Forgot Password
					</CardTitle>
					<CardDescription className="text-center">
						Enter your email address and we’ll send you a password reset link.
					</CardDescription>
				</CardHeader>
				<CardContent>
					{mutation.isSuccess ? (
						<div className="text-green-600 text-center py-8">
							If an account with that email exists, a password reset link has
							been sent.
						</div>
					) : (
						<form onSubmit={handleSubmit} className="space-y-6">
							<div className="space-y-2">
								<Label htmlFor="email">Email Address</Label>
								<Input
									id="email"
									type="email"
									placeholder="you@example.com"
									value={email}
									onChange={(e) => setEmail(e.target.value)}
									required
								/>
							</div>
							{error && <div className="text-red-500 text-sm">{error}</div>}
							<Button
								type="submit"
								className="w-full"
								disabled={mutation.isPending}
							>
								{mutation.isPending ? "Sending..." : "Send Reset Link"}
							</Button>
						</form>
					)}
				</CardContent>
			</Card>
		</div>
	);
};

export default ForgotPasswordPage;
