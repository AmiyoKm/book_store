import { useState, useEffect } from "react";
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Trash2, ShoppingCart } from "lucide-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getCart, removeFromCart, updateQuantity } from "@/config/api/cart";
import type { ApiResponse } from "@/types/books";
import type { Cart } from "@/types/cart";
import { Link } from "react-router";
import { toast } from "sonner";

const CartPage = () => {
	const queryClient = useQueryClient();

	const {
		data: cartData,
		isError,
		isLoading,
	} = useQuery<ApiResponse<Cart>>({
		queryKey: ["cart"],
		queryFn: getCart,
	});

	const [cart, setCart] = useState<Cart | undefined>(undefined);

	const { mutate: mutateCart, isPending } = useMutation({
		mutationFn: updateQuantity,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["cart"] });
		},
	});
	const { mutate: removeItem } = useMutation({
		mutationFn: removeFromCart,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["cart"] });
		},
		onError: () => {
			toast.error("Something went wrong , Please try again later");
		},
	});

	useEffect(() => {
		if (cartData?.data) setCart(cartData.data);
	}, [cartData]);

	if (isLoading) return <p>Loading...</p>;

	if (isError || !cart || !cart.items || cart.items.length === 0) {
		return (
			<div className="flex flex-col items-center justify-center h-[60vh] text-center">
				<ShoppingCart className="w-16 h-16 text-muted-foreground mb-4" />
				<h2 className="text-2xl font-semibold mb-2">Your cart is empty</h2>
				<p className="text-muted-foreground mb-6">
					Looks like you haven't added anything yet.
				</p>
				<Link to={"/books"}>
					<Button variant="outline">Go Shopping</Button>
				</Link>
			</div>
		);
	}

	const handleRemove = (id: number) => {
		removeItem({ cartItemId: id });
	};

	const total = cart.items.reduce(
		(sum, item) => sum + item.price * item.quantity,
		0
	);

	return (
		<div className="container mx-auto px-4 py-8">
			<h1 className="text-3xl font-bold mb-8">Your Cart</h1>
			<div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
				<div className="lg:col-span-2 space-y-6">
					{cart.items.map((item) => (
						<Card
							key={item.id}
							className="flex flex-col md:flex-row gap-4 p-4 shadow-md hover:shadow-lg transition"
						>
							<div className="w-full md:w-32 h-48 shrink-0 overflow-hidden rounded-md border bg-muted">
								<img
									src={item.cover_image_url}
									alt={item.title}
									className="w-full h-full object-cover"
								/>
							</div>

							<CardContent className="flex-1 p-2.5 w-full">
								<CardTitle className="text-lg">{item.title}</CardTitle>
								<p className="text-muted-foreground text-sm mb-1">
									by {item.author}
								</p>
								<p className="font-bold text-primary mb-2">
									${item.price.toFixed(2)}
								</p>

								<div className="flex items-center flex-wrap gap-3 mb-2">
									<label className="text-sm flex items-center gap-1">
										Qty:
										<Input
											type="number"
											min={1}
											max={item.stock}
											value={item.quantity}
											disabled={isPending}
											onChange={(e) => {
												const newQuantity = Number(e.target.value);
												setCart((prevCart) => {
													if (!prevCart) return prevCart;
													return {
														...prevCart,
														items: prevCart.items.map((cartItem) =>
															cartItem.id === item.id
																? { ...cartItem, quantity: newQuantity }
																: cartItem
														),
													};
												});
											}}
											onBlur={(e) => {
												mutateCart({
													quantity: Number(e.target.value),
													itemId: item.id,
												});
											}}
											className="w-20"
										/>
									</label>
									<span className="text-xs text-muted-foreground">
										In Stock: {item.stock}
									</span>
									<Button
										variant="ghost"
										size="icon"
										onClick={() => handleRemove(item.id)}
										className="text-red-500 hover:bg-red-100"
									>
										<Trash2 className="h-5 w-5" />
										<span className="sr-only">Remove</span>
									</Button>
								</div>
							</CardContent>
						</Card>
					))}
				</div>
				<div>
					<Card className="sticky top-24 shadow-md">
						<CardHeader>
							<CardTitle>Order Summary</CardTitle>
						</CardHeader>
						<CardContent>
							<div className="flex justify-between mb-2">
								<span>Subtotal</span>
								<span>${total.toFixed(2)}</span>
							</div>
							<div className="flex justify-between mb-2">
								<span>Shipping</span>
								<span className="text-muted-foreground">Free</span>
							</div>
							<div className="flex justify-between font-bold text-lg mt-4">
								<span>Total</span>
								<span>${total.toFixed(2)}</span>
							</div>
						</CardContent>
						<CardFooter>
							<Button className="w-full bg-primary hover:bg-primary/90 text-primary-foreground text-lg">
								Proceed to Checkout
							</Button>
						</CardFooter>
					</Card>
				</div>
			</div>
		</div>
	);
};

export default CartPage;
