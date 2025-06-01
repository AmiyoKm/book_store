export type CartItem = {
    id: number;
    book_id: number;
    cart_id: number;
    author: string;
    title: string;
    price: number;
    cover_image_url: string,
    stock: number;
    quantity: number;
    created_at: string;
    updated_at: string;
};

export type Cart = {
    cart_id: number;
    items: CartItem[];
};