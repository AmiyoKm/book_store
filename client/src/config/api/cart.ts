import api from "../axios";

export function getCart() {
   return api.get("/carts")
}
export function addToCart({ book_id, quantity }: { book_id: number, quantity: number }) {
   return api.post(`/carts`, { book_id, quantity })
}

export function updateQuantity({quantity , itemId} : {quantity : number , itemId : number}){
   return api.patch(`/carts/items/${itemId}`,{quantity})
}

export function removeFromCart({cartItemId} : {cartItemId : number}){
   return api.delete(`/carts/items/${cartItemId}`)
}