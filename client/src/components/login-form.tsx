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
import { login } from "@/config/api/auth";
import type { LoginPayload } from "@/types/auth";
import { useNavigate } from "react-router";

export function LoginForm({
	className,
	...props
}: React.ComponentPropsWithoutRef<"div">) {
	const navigate = useNavigate();
	const mutation = useMutation({
		mutationFn: login,
		onSuccess: (data) => {
			localStorage.setItem("token", data.data);
			navigate("/home");
		},
	});
	const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		const form = e.currentTarget;

		const email = (form.elements.namedItem("email") as HTMLInputElement).value;
		const password = (form.elements.namedItem("password") as HTMLInputElement)
			.value;

		const payload: LoginPayload = {
			email,
			password,
		};
		mutation.mutate(payload);
	};
	return (
		<div className={cn("flex flex-col gap-6", className)} {...props}>
			<Card className="border-purple-600">
				<CardHeader className="space-y-1">
					<div className="flex items-center justify-center mb-2">
						<BookOpen className="h-10 w-10 text-purple-600" />
					</div>
					<CardTitle className="text-2xl text-center">Welcome back</CardTitle>
					<CardDescription className="text-center">
						Login to your BookBond account
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form onSubmit={handleSubmit}>
						<div className="flex flex-col gap-4">
							<div className="grid gap-2">
								<Label htmlFor="email">Email</Label>
								<Input
									id="email"
									type="email"
									name="email"
									placeholder="you@example.com"
									required
								/>
							</div>
							<div className="grid gap-2">
								<div className="flex items-center">
									<Label htmlFor="password">Password</Label>
									<a
										href="/forgot-password"
										className="ml-auto inline-block text-sm  underline-offset-4 hover:underline"
									>
										Forgot password?
									</a>
								</div>
								<Input id="password" type="password" name="password" required />
							</div>
							<Button
								type="submit"
								className="w-full"
								disabled={mutation.isPending}
							>
								Login
							</Button>
						</div>
						<div className="mt-4 text-center text-sm">
							Don&apos;t have an account?{" "}
							<a href="/sign-up" className="underline underline-offset-4 ">
								Sign up
							</a>
						</div>
					</form>
				</CardContent>
			</Card>
		</div>
	);
}
