export type SignUpPayload = {
    username: string,
    email: string,
    password: string
}
export type LoginPayload = {
    email: string,
    password: string
}
export type ChangePasswordPayload = {
    user_id: number,
    token: string,
    new_password: string
}

export type User = {
    id: number
    username: string
    email: string
    roleId: number
}
