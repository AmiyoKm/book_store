import { ModeToggle } from "@/components/mode-toggle";
import {
	BookOpen,
	ShoppingCart,
	Menu,
	ChevronRight,
	Mail,
	MapPin,
	Phone,
} from "lucide-react";

function App() {
	return (
		<div className="container mx-auto flex min-h-screen flex-col">
			<header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
				<div className="container flex h-16 items-center justify-between">
					<div className="flex items-center gap-2">
						<BookOpen className="h-6 w-6" />
						<span className="text-xl font-bold">BookBond</span>
					</div>
					<nav className="hidden md:flex gap-6">
						<a
							href="#"
							className="text-sm font-medium transition-colors hover:text-primary"
						>
							Home
						</a>
						<a
							href="#books"
							className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
						>
							Books
						</a>
						<a
							href="#categories"
							className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
						>
							Categories
						</a>
						<a
							href="#about"
							className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
						>
							About
						</a>
						<a
							href="#contact"
							className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
						>
							Contact
						</a>
					</nav>
					<div className="flex items-center gap-4">
						<button className="md:hidden inline-flex items-center justify-center rounded-md text-sm font-medium h-9 w-9 text-foreground">
							<Menu className="h-5 w-5" />
							<span className="sr-only">Toggle menu</span>
						</button>
						<button className="inline-flex items-center justify-center rounded-md text-sm font-medium h-9 w-9 text-foreground">
							<ModeToggle />
						</button>
						<button className="inline-flex items-center justify-center rounded-md text-sm font-medium h-9 w-9 text-foreground">
							<ShoppingCart className="h-5 w-5" />
							<span className="sr-only">Cart</span>
						</button>
						<a href="/sign-up">
							<button className="hidden md:inline-flex items-center justify-center rounded-md text-sm font-medium h-9 px-4 bg-primary text-primary-foreground hover:bg-primary/90">
								Sign up
							</button>
						</a>
					</div>
				</div>
			</header>
			<main className="flex-1">
				<section className="w-full py-12 md:py-24 lg:py-32 xl:py-48 bg-gradient-to-b from-background to-muted">
					<div className="container px-4 md:px-6">
						<div className="grid gap-6 lg:grid-cols-2 lg:gap-12 xl:grid-cols-2">
							<div className="flex flex-col justify-center space-y-4">
								<div className="space-y-2">
									<span className="inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 border-transparent bg-primary text-primary-foreground mb-2">
										New Arrivals Weekly
									</span>
									<h1 className="text-3xl font-bold tracking-tighter sm:text-5xl xl:text-6xl/none">
										Discover Your Next Favorite Book
									</h1>
									<p className="max-w-[600px] text-muted-foreground md:text-xl">
										BookBond is your cozy corner for literary adventures. Browse
										our curated collection and find stories that resonate with
										your soul.
									</p>
								</div>
								<div className="flex flex-col gap-2 min-[400px]:flex-row">
									<button className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none bg-primary text-primary-foreground hover:bg-primary/90 h-11 px-8">
										Browse Collection
										<ChevronRight className="ml-2 h-4 w-4" />
									</button>
									<button className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none border border-input bg-background hover:bg-accent hover:text-accent-foreground h-11 px-8">
										Book Club
									</button>
								</div>
							</div>
							<div className="flex items-center justify-center">
								<img
									src="https://plus.unsplash.com/premium_photo-1681488394409-5614ef55488c?q=80&w=1964&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
									width={550}
									height={450}
									alt="BookBond Store Interior"
									className="rounded-lg object-cover"
								/>
							</div>
						</div>
					</div>
				</section>

				<section
					id="books"
					className="w-full py-12 md:py-24 lg:py-32 bg-background"
				>
					<div className="container px-4 md:px-6">
						<div className="flex flex-col items-center justify-center space-y-4 text-center">
							<div className="space-y-2">
								<h2 className="text-3xl font-bold tracking-tighter sm:text-5xl">
									Featured Books
								</h2>
								<p className="max-w-[900px] text-muted-foreground md:text-xl">
									Explore our handpicked selection of this month's most
									captivating reads.
								</p>
							</div>
						</div>
						<div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4 mt-8">
							{[1, 2, 3, 4].map((book) => (
								<div
									key={book}
									className="rounded-lg border bg-card text-card-foreground shadow-sm overflow-hidden"
								>
									<div className="aspect-[2/3] relative">
										<img
											src={`https://m.media-amazon.com/images/I/61GNpAHFttL._AC_UF1000,1000_QL80_.jpg`}
											alt={`Book ${book}`}
											className="h-full w-full object-cover transition-all hover:scale-105"
										/>
									</div>
									<div className="p-4">
										<h3 className="font-semibold">The Great Adventure</h3>
										<p className="text-sm text-muted-foreground">Jane Doe</p>
										<div className="flex items-center justify-between mt-2">
											<span className="font-bold">$19.99</span>
											<button className="inline-flex items-center justify-center rounded-md text-xs font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none border border-input bg-background hover:bg-accent hover:text-accent-foreground h-8 px-3">
												<ShoppingCart className="mr-2 h-4 w-4" />
												Add to Cart
											</button>
										</div>
									</div>
								</div>
							))}
						</div>
						<div className="flex justify-center mt-8">
							<button className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none border border-input bg-background hover:bg-accent hover:text-accent-foreground h-11 px-8">
								View All Books
							</button>
						</div>
					</div>
				</section>

				<section
					id="categories"
					className="w-full py-12 md:py-24 lg:py-32 bg-muted"
				>
					<div className="container px-4 md:px-6">
						<div className="flex flex-col items-center justify-center space-y-4 text-center">
							<div className="space-y-2">
								<h2 className="text-3xl font-bold tracking-tighter sm:text-5xl">
									Browse by Genre
								</h2>
								<p className="max-w-[900px] text-muted-foreground md:text-xl">
									Find your perfect read from our diverse collection of genres.
								</p>
							</div>
						</div>
						<div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6 mt-8">
							{[
								"Fiction",
								"Mystery",
								"Romance",
								"Sci-Fi",
								"Biography",
								"History",
							].map((genre) => (
								<div
									key={genre}
									className="rounded-lg border bg-card text-card-foreground shadow-sm overflow-hidden"
								>
									<div className="p-6 flex flex-col items-center justify-center text-center">
										<h3 className="font-semibold">{genre}</h3>
										<p className="text-sm text-muted-foreground mt-1">
											Explore
										</p>
									</div>
								</div>
							))}
						</div>
					</div>
				</section>

				<section
					id="about"
					className="w-full py-12 md:py-24 lg:py-32 bg-background"
				>
					<div className="container px-4 md:px-6">
						<div className="grid gap-6 lg:grid-cols-2 lg:gap-12 items-center">
							<div className="flex justify-center">
								<img
									src="https://plus.unsplash.com/premium_photo-1706061121381-d7a7a5a06998?q=80&w=2071&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
									width={600}
									height={400}
									alt="BookBond Store"
									className="rounded-lg object-cover"
								/>
							</div>
							<div className="flex flex-col justify-center space-y-4">
								<div className="space-y-2">
									<h2 className="text-3xl font-bold tracking-tighter sm:text-4xl">
										Our Story
									</h2>
									<p className="max-w-[600px] text-muted-foreground md:text-xl">
										Founded in 2010, BookBond began as a small corner shop with
										a big dream: to connect readers with stories that move them.
									</p>
								</div>
								<p className="text-muted-foreground">
									Today, we've grown into a community hub for book lovers,
									hosting regular events, book clubs, and author signings. Our
									knowledgeable staff is passionate about literature and
									dedicated to helping you find your next literary adventure.
								</p>
								<p className="text-muted-foreground">
									We believe in the power of stories to transform lives, spark
									imagination, and build bridges between people. That's why we
									carefully curate our collection to include diverse voices and
									perspectives.
								</p>
								<div>
									<button className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2">
										Learn More About Us
									</button>
								</div>
							</div>
						</div>
					</div>
				</section>

				<section className="w-full py-12 md:py-24 lg:py-32 bg-muted">
					<div className="container px-4 md:px-6">
						<div className="flex flex-col items-center justify-center space-y-4 text-center">
							<div className="space-y-2">
								<h2 className="text-3xl font-bold tracking-tighter sm:text-5xl">
									What Our Readers Say
								</h2>
								<p className="max-w-[900px] text-muted-foreground md:text-xl">
									Don't just take our word for it. Here's what our community has
									to say.
								</p>
							</div>
						</div>
						<div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 mt-8">
							{[
								{
									name: "Alex Johnson",
									quote:
										"BookBond has transformed my reading experience. Their staff recommendations are always spot on!",
								},
								{
									name: "Sarah Williams",
									quote:
										"I love the cozy atmosphere and the carefully curated selection. It's my happy place in the city.",
								},
								{
									name: "Michael Chen",
									quote:
										"The book club has introduced me to genres I never thought I'd enjoy. Now I can't get enough!",
								},
							].map((testimonial, index) => (
								<div
									key={index}
									className="rounded-lg border bg-card text-card-foreground shadow-sm p-6"
								>
									<div className="pt-4">
										<p className="mb-4 italic">"{testimonial.quote}"</p>
										<p className="font-semibold">— {testimonial.name}</p>
									</div>
								</div>
							))}
						</div>
					</div>
				</section>

				<section className="w-full py-12 md:py-24 lg:py-32 bg-primary text-primary-foreground">
					<div className="container px-4 md:px-6">
						<div className="flex flex-col items-center justify-center space-y-4 text-center">
							<div className="space-y-2">
								<h2 className="text-3xl font-bold tracking-tighter sm:text-4xl">
									Join Our Newsletter
								</h2>
								<p className="max-w-[600px] md:text-xl">
									Stay updated with new arrivals, events, and exclusive offers.
								</p>
							</div>
							<div className="w-full max-w-md space-y-2">
								<form className="flex space-x-2">
									<input
										className="flex h-10 w-full rounded-md border border-input bg-primary-foreground px-3 py-2 text-sm text-primary ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-primary/60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 max-w-lg flex-1"
										placeholder="Enter your email"
										type="email"
									/>
									<button className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none bg-secondary text-secondary-foreground hover:bg-secondary/80 h-10 px-4 py-2">
										Subscribe
									</button>
								</form>
								<p className="text-xs">
									By subscribing, you agree to our terms and privacy policy.
								</p>
							</div>
						</div>
					</div>
				</section>
			</main>
			<footer
				id="contact"
				className="w-full border-t bg-background py-6 md:py-12"
			>
				<div className="container px-4 md:px-6">
					<div className="grid gap-8 sm:grid-cols-2 md:grid-cols-4">
						<div className="space-y-4">
							<div className="flex items-center gap-2">
								<BookOpen className="h-6 w-6" />
								<span className="text-lg font-bold">BookBond</span>
							</div>
							<p className="text-sm text-muted-foreground">
								Your cozy corner for literary adventures since 2010.
							</p>
						</div>
						<div className="space-y-4">
							<h3 className="text-base font-medium">Quick Links</h3>
							<ul className="space-y-2 text-sm">
								<li>
									<a
										href="#"
										className="text-muted-foreground hover:text-foreground"
									>
										Home
									</a>
								</li>
								<li>
									<a
										href="#books"
										className="text-muted-foreground hover:text-foreground"
									>
										Books
									</a>
								</li>
								<li>
									<a
										href="#categories"
										className="text-muted-foreground hover:text-foreground"
									>
										Categories
									</a>
								</li>
								<li>
									<a
										href="#about"
										className="text-muted-foreground hover:text-foreground"
									>
										About
									</a>
								</li>
							</ul>
						</div>
						<div className="space-y-4">
							<h3 className="text-base font-medium">Customer Service</h3>
							<ul className="space-y-2 text-sm">
								<li>
									<a
										href="#"
										className="text-muted-foreground hover:text-foreground"
									>
										FAQ
									</a>
								</li>
								<li>
									<a
										href="#"
										className="text-muted-foreground hover:text-foreground"
									>
										Shipping & Returns
									</a>
								</li>
								<li>
									<a
										href="#"
										className="text-muted-foreground hover:text-foreground"
									>
										Privacy Policy
									</a>
								</li>
								<li>
									<a
										href="#"
										className="text-muted-foreground hover:text-foreground"
									>
										Terms of Service
									</a>
								</li>
							</ul>
						</div>
						<div className="space-y-4">
							<h3 className="text-base font-medium">Contact Us</h3>
							<ul className="space-y-2 text-sm">
								<li className="flex items-start gap-2">
									<MapPin className="h-4 w-4 mt-0.5" />
									<span className="text-muted-foreground">
										123 Book Lane, Reading, RG1 2CD
									</span>
								</li>
								<li className="flex items-center gap-2">
									<Phone className="h-4 w-4" />
									<span className="text-muted-foreground">(123) 456-7890</span>
								</li>
								<li className="flex items-center gap-2">
									<Mail className="h-4 w-4" />
									<span className="text-muted-foreground">
										hello@bookbond.com
									</span>
								</li>
							</ul>
						</div>
					</div>
					<div className="mt-8 border-t pt-8 text-center text-sm text-muted-foreground">
						<p>© {new Date().getFullYear()} BookBond. All rights reserved.</p>
					</div>
				</div>
			</footer>
		</div>
	);
}

export default App;
