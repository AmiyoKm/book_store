import api from "../axios";


export function fetchBooks() {
    return api.get("/books/search")
}