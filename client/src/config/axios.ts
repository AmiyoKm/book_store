import axios from 'axios'


const api = axios.create({
    baseURL: import.meta.env.BACKEND_BASE_ENDPOINT || "http://localhost:8080/api/v1",
})

api.interceptors.request.use(
    config => {
        const token = localStorage.getItem("token")
        if (token) {
            config.headers["Authorization"] = `Bearer ${token}`
        }
        return config
    },
    error => Promise.reject(error)
)
api.interceptors.response.use(
    response => response,
    error => {
        if (error.response && error.response.status === 401) {
            //Optionally remove token and redirect to login
            // localStorage.removeItem("token");
            // window.location.href = "/sign-in";
        }
        return Promise.reject(error);
    }
);

api.defaults.headers.common['Content-Type'] = 'application/json';
api.defaults.timeout = 10000;

export default api