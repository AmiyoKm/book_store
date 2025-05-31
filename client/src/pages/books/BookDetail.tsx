import { useParams, useNavigate, Link } from "react-router-dom";
import {
	useQueries,
	useQuery,
	useMutation,
	useQueryClient,
} from "@tanstack/react-query";
import {
	fetchBookById,
	fetchBookByTags,
	fetchReviewsOfBook,
	submitBookReview,
} from "@/config/api/books";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ChevronLeft, Star } from "lucide-react";
import type { ApiResponse, Book, Review } from "@/types/books";
import { StarRating } from "./components/StarRating";
import { useState } from "react";
import { Textarea } from "@/components/ui/textarea"; // Import Textarea for review input
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { getUser } from "@/config/api/auth";
import type { User } from "@/types/auth";

const BookDetailPage = () => {
	const { bookId } = useParams<{ bookId: string }>();
	const navigate = useNavigate();
	const queryClient = useQueryClient(); // Get query client for invalidation
	const { data: UserData, isError: isUserError } = useQuery<ApiResponse<User>>({
		queryKey: ["user"],
		queryFn: getUser,
	});
	const [reviewText, setReviewText] = useState("");
	const [reviewRating, setReviewRating] = useState(0); // 0-5 stars
	const {
		data: bookResponse,
		isLoading: isLoadingBook,
		isError: isErrorBook,
		error: bookError,
	} = useQuery<ApiResponse<Book>>({
		queryKey: ["book", bookId],
		queryFn: () => fetchBookById(bookId!),
		enabled: !!bookId,
	});
	const {
		data: reviewsResponse,
		isLoading: isLoadingReviews,
		isError: isErrorReviews,
		error: reviewsError,
	} = useQuery<ApiResponse<Review[]>>({
		queryKey: ["bookReviews", bookId],
		queryFn: () => fetchReviewsOfBook(bookId!),
		enabled: !!bookId,
		staleTime: 5 * 60 * 1000,
	});

	const reviewMutation = useMutation({
		mutationKey: ["review", bookId],
		mutationFn: submitBookReview,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["bookReviews", bookId] });
			setReviewText("");
			setReviewRating(0);
			toast.success("Review Submitted!", {
				description: "Your review has been successfully added.",
			});
		},
		onError: (error) => {
			toast.error("Review Submission Failed", {
				description:
					error.message || "Could not submit your review. Please try again.",
			});
		},
	});
	if (isUserError) {
		toast.error("User not found");
		return (
			<div className="container mx-auto px-4 py-8 text-center text-lg text-red-500">
				User not found . Please Log in .
			</div>
		);
	}
	const book = bookResponse?.data;
	const relatedBookQueriesConfig = (book?.tags || []).map((tag) => ({
		queryKey: ["relatedBooksByTag", tag],
		queryFn: () => fetchBookByTags(tag),
		enabled: !!book && book.tags?.length > 0,
		staleTime: 5 * 60 * 1000,
		cacheTime: 10 * 60 * 1000,
	}));

	const relatedQueries = useQueries({
		queries: relatedBookQueriesConfig,
	});

	if (!bookId) {
		return (
			<div className="container mx-auto px-4 py-8 text-center text-lg text-red-500">
				Book ID is missing from the URL.
			</div>
		);
	}

	if (isLoadingBook) {
		return (
			<div className="container mx-auto px-4 py-8 text-center text-xl text-muted-foreground">
				Loading book details...
			</div>
		);
	}

	if (isErrorBook) {
		return (
			<div className="container mx-auto px-4 py-8 text-center text-xl text-red-500">
				Error loading book details:{" "}
				{bookError?.message || "Something went wrong."}
			</div>
		);
	}

	const reviews = reviewsResponse?.data || [];

	// Calculate average rating
	const averageRating =
		reviews.length > 0
			? reviews.reduce((sum, review) => sum + review.Rating, 0) / reviews.length
			: 0;

	// --- Related Books Queries ---

	const isLoadingRelatedBooks = relatedQueries.some((query) => query.isLoading);
	const isErrorRelatedBooks = relatedQueries.some((query) => query.isError);
	const relatedBooksError = relatedQueries.find(
		(query) => query.isError
	)?.error;

	console.log("isLoadingRelatedBooks:", isLoadingRelatedBooks);
	console.log("isErrorRelatedBooks:", isErrorRelatedBooks);
	if (isErrorRelatedBooks) {
		console.error("Related Books Error Details:", relatedBooksError);
	}

	let combinedRelatedBooks: Book[] = [];
	const seenBookIds = new Set<number>();

	if (!isLoadingRelatedBooks && !isErrorRelatedBooks) {
		relatedQueries.forEach((queryResult) => {
			if (queryResult.data?.data && Array.isArray(queryResult.data.data)) {
				queryResult.data.data.forEach((b: Book) => {
					if (book && b.id !== book.id && !seenBookIds.has(b.id)) {
						combinedRelatedBooks.push(b);
						seenBookIds.add(b.id);
					}
				});
			}
		});
	}
	const relatedBooks = combinedRelatedBooks;

	console.log("Final Related Books for Display:", relatedBooks);

	if (!book) {
		return (
			<div className="container mx-auto px-4 py-8 text-center text-xl text-muted-foreground">
				Book not found.
			</div>
		);
	}

	const handleReviewSubmit = async () => {
		if (!bookId || !reviewText.trim() || reviewRating === 0) {
			toast.error("Missing Review Info", {
				description:
					"Please provide both text and a star rating for your review.",
			});
			return;
		}

		reviewMutation.mutate({
			bookId,
			content: reviewText,
			rating: reviewRating,
		});
	};

	return (
		<div className="container mx-auto px-4 py-8 bg-background text-foreground">
			<Button
				variant="ghost"
				onClick={() => navigate("/books")}
				className="mb-6 text-muted-foreground hover:text-foreground"
			>
				<ChevronLeft className="mr-2 h-4 w-4" /> Back to Books
			</Button>

			{/* Main layout: Grid with main content and sidebar */}
			<div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
				{/* Main Book Details (2 or 3 columns) */}
				<div className="lg:col-span-3 space-y-8">
					<Card className="mx-auto overflow-hidden rounded-lg shadow-xl bg-card border-border">
						<CardHeader className="p-0">
							<img
								src={
									book.cover_image_url ||
									`https://placehold.co/400x600/DDDDDD/666666?text=No+Image`
								}
								alt={book.title}
								className="w-full h-96 object-contain bg-muted rounded-t-lg"
								onError={(e) => {
									const target = e.target as HTMLImageElement;
									target.onerror = null;
									target.src = `https://placehold.co/400x600/DDDDDD/666666?text=No+Image`;
								}}
							/>
						</CardHeader>
						<CardContent className="p-6 space-y-4">
							<CardTitle className="text-4xl font-bold text-card-foreground">
								{book.title}
							</CardTitle>
							<div className="flex items-center space-x-4">
								<CardDescription className="text-xl text-muted-foreground">
									by {book.author}
								</CardDescription>
								{/* Average Rating as Stars */}
								{reviews.length > 0 ? (
									<StarRating rating={averageRating} />
								) : (
									<span className="text-sm text-muted-foreground">
										No ratings yet.
									</span>
								)}
							</div>

							<p className="text-base text-foreground leading-relaxed">
								{book.description || "No description available."}
							</p>

							<div className="grid grid-cols-2 gap-4 text-sm text-muted-foreground">
								{book.isbn && (
									<div>
										<strong>ISBN:</strong> {book.isbn}
									</div>
								)}
								{book.pages > 0 && (
									<div>
										<strong>Pages:</strong> {book.pages}
									</div>
								)}
								{book.tags && book.tags.length > 0 && (
									<div>
										<strong>Tags:</strong> {book.tags.join(", ")}
									</div>
								)}
								{book.stock > 0 && (
									<div>
										<strong>Stock:</strong> {book.stock}
									</div>
								)}
								{book.createdAt && (
									<div>
										<strong>Published:</strong>{" "}
										{new Date(book.createdAt).toLocaleDateString()}
									</div>
								)}
							</div>
						</CardContent>
						<CardFooter className="flex justify-between items-center p-6 border-t border-border">
							<span className="text-4xl font-extrabold text-primary">
								${book.price.toFixed(2)}
							</span>
							<Button className="bg-primary hover:bg-primary/90 text-primary-foreground text-lg px-6 py-3">
								Add to Cart
							</Button>
						</CardFooter>
					</Card>

					{/* Add Review Section */}
					<Card className="bg-card border-border shadow-lg mt-8">
						<CardHeader>
							<CardTitle className="text-2xl text-card-foreground">
								Add Your Review
							</CardTitle>
						</CardHeader>
						<CardContent className="space-y-4">
							<div>
								<Label
									htmlFor="review-rating"
									className="text-foreground mb-2 block"
								>
									Your Rating
								</Label>
								<div className="flex space-x-1">
									{[1, 2, 3, 4, 5].map((star) => (
										<Star
											key={star}
											className={`h-7 w-7 cursor-pointer transition-colors duration-200 ${
												reviewRating >= star
													? "fill-primary stroke-primary"
													: "fill-transparent stroke-muted-foreground"
											}`}
											onClick={() => setReviewRating(star)}
										/>
									))}
								</div>
							</div>
							<div>
								<Label
									htmlFor="review-text"
									className="text-foreground mb-2 block"
								>
									Your Review
								</Label>
								<Textarea
									id="review-text"
									placeholder="Write your review here..."
									value={reviewText}
									onChange={(e) => setReviewText(e.target.value)}
									className="min-h-[100px] bg-input text-foreground border-border focus:ring-primary focus:border-primary"
								/>
							</div>
							<Button
								onClick={handleReviewSubmit}
								disabled={
									reviewMutation.isPending ||
									reviewRating === 0 ||
									reviewText.trim() === ""
								}
								className="bg-primary hover:bg-primary/90 text-primary-foreground w-full"
							>
								{reviewMutation.isPending ? "Submitting..." : "Submit Review"}
							</Button>
							{reviewMutation.isError && (
								<p className="text-red-500 text-sm mt-2">
									Error:{" "}
									{reviewMutation.error?.message || "Failed to submit review."}
								</p>
							)}
						</CardContent>
					</Card>

					{/* Customer Reviews Section */}
					<Card className="bg-card border-border shadow-lg mt-8">
						<CardHeader>
							<CardTitle className="text-2xl text-card-foreground">
								Customer Reviews ({reviews.length})
							</CardTitle>
						</CardHeader>
						<CardContent className="space-y-6">
							{isLoadingReviews ? (
								<p className="text-muted-foreground">Loading reviews...</p>
							) : isErrorReviews ? (
								<p className="text-red-500">
									Error loading reviews: {reviewsError?.message}
								</p>
							) : reviews.length > 0 ? (
								reviews.map((review) => (
									<div
										key={review.ID}
										className="border-b border-border pb-4 last:border-b-0"
									>
										<div className="flex items-center justify-between mb-2">
											<StarRating rating={review.Rating} />
											<span className="text-sm text-muted-foreground">
												{review.CreatedAt
													? new Date(review.CreatedAt).toLocaleDateString()
													: "Date N/A"}
											</span>
										</div>
										<p className="text-foreground text-base leading-relaxed">
											{review.Content}
										</p>
										<p className="text-xs text-muted-foreground mt-1">
											— {" " + UserData?.data.username || "Anonymous"}
										</p>
									</div>
								))
							) : (
								<p className="text-muted-foreground text-center">
									No reviews yet. Be the first to review this book!
								</p>
							)}
						</CardContent>
					</Card>
				</div>

				{/* Related Books Sidebar (1 column) */}
				<div className="lg:col-span-1">
					<Card className="bg-card border-border shadow-lg">
						<CardHeader>
							<CardTitle className="text-2xl text-card-foreground">
								More Like This
							</CardTitle>
						</CardHeader>
						<CardContent className="space-y-4">
							{isLoadingRelatedBooks ? (
								<p className="text-muted-foreground">
									Loading related books...
								</p>
							) : isErrorRelatedBooks ? (
								<p className="text-red-500">
									Error loading related books:{" "}
									{relatedBooksError?.message || "Unknown error."}
								</p>
							) : relatedBooks.length > 0 ? (
								<div className="space-y-4">
									{relatedBooks.map((relatedBook) => (
										<Link
											to={`/books/${relatedBook.id}`}
											key={relatedBook.id}
											className="block hover:bg-muted rounded-lg transition-colors duration-200 p-2"
										>
											<div className="flex items-center space-x-3">
												<img
													src={
														relatedBook.cover_image_url ||
														`https://placehold.co/60x90?text=No+Image`
													}
													alt={relatedBook.title}
													className="w-16 h-24 object-cover rounded-md flex-shrink-0"
													onError={(e) => {
														const target = e.target as HTMLImageElement;
														target.onerror = null;
														target.src = `https://placehold.co/60x90?text=No+Image`;
													}}
												/>
												<div>
													<p className="text-lg font-semibold text-card-foreground line-clamp-2">
														{relatedBook.title}
													</p>
													<p className="text-sm text-muted-foreground">
														by {relatedBook.author}
													</p>
													<p className="text-sm font-bold text-primary">
														${relatedBook.price.toFixed(2)}
													</p>
												</div>
											</div>
										</Link>
									))}
								</div>
							) : (
								<p className="text-muted-foreground">No related books found.</p>
							)}
						</CardContent>
					</Card>
				</div>
			</div>
		</div>
	);
};

export default BookDetailPage;
