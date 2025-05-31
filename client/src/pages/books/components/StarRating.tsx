import { Star } from "lucide-react";

export const StarRating = ({
	rating,
	maxRating = 5,
}: {
	rating: number;
	maxRating?: number;
}) => {
	const fullStars = Math.floor(rating);
	const hasHalfStar = rating % 1 !== 0;
	const emptyStars = maxRating - fullStars - (hasHalfStar ? 1 : 0);

	return (
		<div className="flex items-center text-primary">
			{[...Array(fullStars)].map((_, i) => (
				<Star
					key={`full-${i}`}
					className="h-5 w-5 fill-primary stroke-primary"
				/>
			))}
			{hasHalfStar && (
				<Star
					key="half"
					className="h-5 w-5 fill-transparent stroke-primary"
					style={{ clipPath: "inset(0 50% 0 0)" }}
				/>
			)}
			{[...Array(emptyStars)].map((_, i) => (
				<Star
					key={`empty-${i}`}
					className="h-5 w-5 fill-transparent stroke-current text-muted-foreground"
				/>
			))}
			<span className="ml-2 text-sm font-semibold text-foreground">
				{rating.toFixed(1)} / {maxRating}
			</span>
		</div>
	);
};
