import {
	ShoppingCart,
	User,
	Heart,
	Home,
	ListOrdered,
	Menu,
    BookOpen,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { ModeToggle } from "@/components/mode-toggle";
import { useNavigate } from "react-router";

const Navbar = () => {
	const navigate = useNavigate()
	return (
		<nav className="border-b border-border px-4 py-2 flex items-center justify-between bg-background text-foreground shadow-sm">
			<div
				className="flex items-center"
				onClick={() => {
					navigate("/books");
				}}
			>
				<BookOpen className="mr-2 h-6 w-6 text-primary" />{" "}
				<span className="text-xl font-bold text-primary">BookBond</span>{" "}
			</div>
			<div className="hidden md:flex items-center space-x-6">
				<Button
					variant="ghost"
					className="flex items-center text-foreground hover:text-primary"
					onClick={() => {
						navigate("/books");
					}}
				>
					<Home className="mr-2 h-4 w-4" /> Home
				</Button>
				<Button
					variant="ghost"
					className="flex items-center text-foreground hover:text-primary"
					onClick={() => {
						navigate("/wishlist");
					}}
				>
					<Heart className="mr-2 h-4 w-4" /> Wishlist
				</Button>
				<Button
					variant="ghost"
					className="text-foreground hover:text-primary"
					onClick={() => {
						navigate("/categories");
					}}
				>
					Categories
				</Button>
				<Button
					variant="ghost"
					className="flex items-center text-foreground hover:text-primary"
					onClick={() => {
						navigate("/orders");
					}}
				>
					<ListOrdered className="mr-2 h-4 w-4" /> Orders
				</Button>
			</div>
			<div className="flex items-center space-x-4">
				<Button
					variant="ghost"
					size="icon"
					className="text-foreground hover:text-primary"
					onClick={() => {
						navigate("/cart");
					}}
				>
					<ShoppingCart className="h-5 w-5" />
					<span className="sr-only">Shopping Cart</span>{" "}
				</Button>
				<Button
					variant="ghost"
					size="icon"
					className="text-foreground hover:text-primary"
					onClick={() => {
						navigate("/profile");
					}}
				>
					<User className="h-5 w-5" />
					<span className="sr-only">Profile</span>
				</Button>
				<ModeToggle />
				<div className="md:hidden">
					<Sheet>
						<SheetTrigger asChild>
							<Button
								variant="ghost"
								size="icon"
								className="text-foreground hover:text-primary"
							>
								<Menu className="h-6 w-6" />
								<span className="sr-only">Open mobile menu</span>{" "}
							</Button>
						</SheetTrigger>
						<SheetContent
							side="right"
							className="bg-sidebar text-sidebar-foreground border-l border-sidebar-border"
						>
							<div className="flex flex-col space-y-4 pt-8">
								<Button
									variant="ghost"
									className="justify-start flex items-center text-sidebar-foreground hover:text-sidebar-primary"
								>
									<Home className="mr-2 h-4 w-4" /> Home
								</Button>
								<Button
									variant="ghost"
									className="justify-start flex items-center text-sidebar-foreground hover:text-sidebar-primary"
								>
									<Heart className="mr-2 h-4 w-4" /> Wishlist
								</Button>
								<Button
									variant="ghost"
									className="justify-start text-sidebar-foreground hover:text-sidebar-primary"
								>
									Categories
								</Button>
								<Button
									variant="ghost"
									className="justify-start flex items-center text-sidebar-foreground hover:text-sidebar-primary"
								>
									<ListOrdered className="mr-2 h-4 w-4" /> Orders
								</Button>
							</div>
						</SheetContent>
					</Sheet>
				</div>
			</div>
		</nav>
	);
};

export default Navbar;
