export type Book = {
    id: number;
    title: string;
    author: string;
    isbn: string;
    price: number;
    tags: string[];
    description: string;
    cover_image_url: string;
    pages: number;
    stock: number;
    createdAt: string;
    updatedAt: string;
    version: number;
}
export type ApiResponse<T> = {
    data: T;
};

export type Review = {
    ID: number;
    UserID: number;
    BookID: number;
    Content: string;
    Rating: number; 
    CreatedAt: string;
    UpdatedAt: string;
}