import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Search } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { fetchBooks } from "@/config/api/books";
import type { ApiResponse, Book } from "@/types/books";

const BooksPage = () => {
	const {
		data: booksData,
		isLoading,
		isError,
		error,
	} = useQuery<any, Error, { data: ApiResponse<Book[]> }>({
		queryKey: ["book"],
		queryFn: fetchBooks,
	});

	const [searchTerm, setSearchTerm] = useState("");
	const [filteredBooks, setFilteredBooks] = useState<Book[]>([]);

	useEffect(() => {
		if (booksData?.data) {
			const lowercasedSearchTerm = searchTerm.toLowerCase();
			const results = booksData.data.data.filter(
				(book: Book) =>
					book.title.toLowerCase().includes(lowercasedSearchTerm) ||
					book.author.toLowerCase().includes(lowercasedSearchTerm) ||
					(book.tags &&
						book.tags.some((tag) =>
							tag.toLowerCase().includes(lowercasedSearchTerm)
						)) ||
					book.isbn.toLowerCase().includes(lowercasedSearchTerm)
			);
			setFilteredBooks(results);
		}
	}, [searchTerm, booksData]);

	if (isLoading) {
		return (
			<div className="container mx-auto px-4 py-8 text-center text-xl text-muted-foreground">
				Loading books...
			</div>
		);
	}

	if (isError) {
		return (
			<div className="container mx-auto px-4 py-8 text-center text-xl text-red-500">
				Error loading books: {error?.message || "Something went wrong."}
			</div>
		);
	}

	return (
		<div className="container mx-auto px-4 py-8 bg-background text-foreground">
			<h1 className="text-4xl font-extrabold text-center mb-8 text-foreground">
				Our Collection
			</h1>

			<div className="mb-10 flex justify-center">
				<div className="relative w-full max-w-md">
					<Input
						type="text"
						placeholder="Search books by title, author, genre, or ISBN..."
						className="w-full pl-10 pr-4 py-2 border border-border rounded-full shadow-sm focus:outline-none focus:ring-2 focus:ring-primary bg-input text-foreground"
						value={searchTerm}
						onChange={(e) => setSearchTerm(e.target.value)}
					/>
					<Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
				</div>
			</div>

			{/* Book Grid Display */}
			{filteredBooks.length > 0 ? (
				<div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
					{filteredBooks.map((book: Book) => {
						console.log(book.cover_image_url);
						return (
							<Card
								key={book.id}
								className="overflow-hidden rounded-lg shadow-lg hover:shadow-xl transition-shadow duration-300 bg-card border-border"
							>
								<CardHeader className="p-0">
									<img
										src={book.cover_image_url}
										alt={book.title}
										className="w-full h-60 object-cover rounded-t-lg"
									/>
								</CardHeader>
								<CardContent className="p-4">
									<CardTitle className="text-lg font-semibold mb-1 text-card-foreground line-clamp-2">
										{book.title}
									</CardTitle>
									<CardDescription className="text-sm text-muted-foreground mb-2">
										by {book.author}
									</CardDescription>
									<p className="text-sm text-foreground line-clamp-3 mb-2">
										{book.description}
									</p>
									{/* Displaying Pages */}
									{book.pages && (
										<p className="text-xs text-muted-foreground">
											{book.pages} pages
										</p>
									)}
								</CardContent>
								<CardFooter className="flex justify-between items-center p-4 pt-0">
									<span className="text-xl font-bold text-primary">
										${book.price.toFixed(2)}
									</span>
									<div className="flex items-center space-x-2">
										<Button className="bg-primary hover:bg-primary/90 text-primary-foreground transition-colors duration-200">
											Add to Cart
										</Button>
									</div>
								</CardFooter>
							</Card>
						);
					})}
				</div>
			) : (
				<div className="text-center text-muted-foreground text-lg py-10">
					No books found matching your search.
				</div>
			)}
		</div>
	);
};

export default BooksPage;
