import type React from "react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { BookOpen } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { signUp } from "@/config/api/auth";
import { useState } from "react";

export function SignupForm({
	className,
	...props
}: React.ComponentPropsWithoutRef<"div">) {
	const [emailSent, setEmailSent] = useState<Boolean>(false);
	const mutation = useMutation({
		mutationFn: signUp,
		onSuccess: () => {
			setEmailSent(true);
		},
	});
	const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		const form = e.currentTarget;

		const username = (form.elements.namedItem("username") as HTMLInputElement)
			.value;
		const email = (form.elements.namedItem("email") as HTMLInputElement).value;
		const password = (form.elements.namedItem("password") as HTMLInputElement)
			.value;
		mutation.mutate({ username, email, password });
	};
	return (
		<div className={cn("flex flex-col gap-6", className)} {...props}>
			<Card className="border-primary-100">
				<CardHeader className="space-y-1">
					<div className="flex items-center justify-center mb-2">
						<BookOpen className="h-10 w-10 text-primary-600" />
					</div>
					<CardTitle className="text-2xl text-center">
						Create an account
					</CardTitle>
					<CardDescription className="text-center">
						Join BookBond and discover your next favorite book
					</CardDescription>
				</CardHeader>
				<CardContent>
					{emailSent ? (
						<div className="text-center py-8">
							<h2 className="text-xl font-semibold mb-2">Check your email!</h2>
							<p className="text-muted-foreground">
								We've sent a confirmation email to your inbox. Please follow the
								instructions to complete your registration.
							</p>
						</div>
					) : (
						<form onSubmit={handleSubmit}>
							<div className="flex flex-col gap-4">
								<div className="grid gap-2">
									<Label htmlFor="username">Username</Label>
									<Input
										id="username"
										name="username"
										placeholder="bookworm123"
										required
									/>
								</div>
								<div className="grid gap-2">
									<Label htmlFor="email">Email</Label>
									<Input
										id="email"
										name="email"
										type="email"
										placeholder="you@example.com"
										required
									/>
								</div>
								<div className="grid gap-2">
									<Label htmlFor="password">Password</Label>
									<Input
										id="password"
										name="password"
										type="password"
										required
									/>
								</div>
								<Button type="submit" disabled={mutation.isPending}>
									{mutation.isPending ? "Signing Up..." : "Sign Up"}
								</Button>
							</div>
							<div className="mt-4 text-center text-sm">
								Already have an account?{" "}
								<a
									href="/sign-in"
									className="text-primary-600 underline underline-offset-4 hover:text-primary-700"
								>
									Login
								</a>
							</div>
						</form>
					)}
				</CardContent>
			</Card>
		</div>
	);
}
