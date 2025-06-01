import type { ChangePasswordPayload, LoginPayload, SignUpPayload } from "@/types/auth";
import api from "../axios";

export function signUp(data: SignUpPayload) {
    return api.post("/authentication/user", data)
}

export function login(data: LoginPayload) {
    return api.post("/authentication/token", data)
}
export function activateUser(token: string) {
    return api.put(`/authentication/activate/${token}`);
}

export function passwordCheckValidation(token: string) {
    return api.get(`/password/request/verify?token=${token}`)
}

export function passwordChange(data: ChangePasswordPayload) {
    return api.post("/password/reset", data)
}

export function passwordResetRequest(data: { email: string }) {
    return api.post("/password/reset-request", data)
}

export function getUser() {
    return api.get("/users/me")
}
export function getUserByID(id: number) {
    return api.get(`/users/${id}`)
}
