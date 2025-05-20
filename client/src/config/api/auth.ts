import type { LoginPayload, SignUpPayload } from "@/types/auth";
import api from "../axios";

export function signUp(data: SignUpPayload) {
    return api.post("/authentication/user", data)
}

export function login(data: LoginPayload) {
    return api.post("/authentication/token", data)
}