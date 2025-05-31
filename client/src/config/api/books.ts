import api from "../axios";


export async function fetchBooks() {
    return await api.get("/books/search")
}

export async function fetchBookById(id: string) {
    return await api.get(`/books/${id}`)
}

export async function fetchBookByTags(tag: string) {
    return await api.get(`/books/search?tag=${tag}`)
}
export async function fetchReviewsOfBook(id: string) {
    return await api.get(`/books/${id}/reviews`)
}

export async function submitBookReview({ bookId, content, rating }: { bookId: string, content: string, rating: number }) {
    return await api.post(`/books/${bookId}/reviews`, { content, rating })
}