<template>
    <div>
        <div>
            <div v-if="LoginOrRegister === 'Login'">
                <h1>Sign In</h1>
                <form @submit.prevent="handleLogin">
                    <div>
                        <label for="username">Username:</label>
                        <input type="text" id="username" v-model="loginForm.username" required />
                    </div>
                    <div>
                        <label for="password">Password:</label>
                        <input type="password" id="password" v-model="loginForm.password" required />
                    </div>
                    <button type="submit">Login</button>
                </form>
                <p>Don't have an account?</p>
                <button @click="log">Create Account</button>
            </div>
            <div v-else>
                <h1>Register</h1>
                <form @submit.prevent="handleRegister">
                    <div>
                        <label for="new-username">Username:</label>
                        <input type="text" id="new-username" v-model="registerForm.username" required />
                    </div>
                    <div>
                        <label for="new-password">Password:</label>
                        <input type="password" id="new-password" v-model="registerForm.password" required />
                    </div>
                    <button type="submit">Register</button>
                </form>
                <p>Already have an account?</p>
                <button @click="log">Login</button>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    name: "LoginRegister",
    data() {
        return {
            LoginOrRegister: "Login",
            loginForm: {
                username: '',
                password: ''
            },
            registerForm: {
                username: '',
                password: ''
            },
            errorMessage: '',
            isLoading: false
        };
    },
    created() {
        this.checkExistingToken();
    },
    methods: {
        async checkExistingToken() {
            // get token from localStorage if it exists
            const token = localStorage.getItem('authToken');
            if (token) {
                console.log("Token found in localStorage:", token);
                try {
                    const response = await axios.post(`${this.$apiUrl}validate/token`, {
                        token: token
                    });
                    console.log("Token validation response:", response.data);
                    if (response.data.valid) {
                        console.log("Token is valid, user is authenticated");
                        this.$emit('token-validated', response.data);
                    }
                } catch (error) {
                    console.error("Token validation failed:", error);
                    localStorage.removeItem('authToken');
                    localStorage.removeItem('tokenExpiry');
                }
            }
        },
        log() {
            this.errorMessage = '';
            if (this.LoginOrRegister === "Login") {
                this.LoginOrRegister = "Register";
            } else {
                this.LoginOrRegister = "Login";
            }
        },
        async handleLogin() {
            this.isLoading = true;
            this.errorMessage = '';

            try {
                console.log("Sending login request to:", `${this.$apiUrl}login`);
                console.log("Login form data:", this.loginForm);

                const response = await axios.post(`${this.$apiUrl}login`, {
                    username: this.loginForm.username,
                    password: this.loginForm.password
                });

                console.log("Login successful:", response.data);

                // Store token in localStorage
                if (response.data.token) {
                    localStorage.setItem('authToken', response.data.token);
                    localStorage.setItem('tokenExpiry', Date.now() + (365 * 60 * 60 * 1000));
                    console.log("Token saved to localStorage");

                    // Optionally redirect or emit event to parent component
                    this.$emit('login-success', response.data);
                }

            } catch (error) {
                console.error("Login error:", error);

                if (error.response) {
                    this.errorMessage = error.response.data.error || 'Login failed';
                    console.error("Server error:", error.response.data);
                } else if (error.request) {
                    this.errorMessage = 'Cannot connect to server. Please check if the backend is running.';
                    console.error("Network error:", error.request);
                } else {
                    this.errorMessage = 'An unexpected error occurred';
                    console.error("Error:", error.message);
                }
            } finally {
                this.isLoading = false;
            }
        },
        async handleRegister() {
            this.isLoading = true;
            this.errorMessage = '';

            try {
                console.log("Sending register request to:", `${this.$apiUrl}register`);
                console.log("Register form data:", this.registerForm);

                const response = await axios.post(`${this.$apiUrl}register`, {
                    username: this.registerForm.username,
                    password: this.registerForm.password
                });

                console.log("Registration successful:", response.data);
                console.log("Token:", response.data.token);

                // Store token in localStorage
                if (response.data.token) {
                    localStorage.setItem('authToken', response.data.token);
                    localStorage.setItem('tokenExpiry', Date.now() + (365 * 60 * 60 * 1000));
                    console.log("Token saved to localStorage");

                    // Optionally redirect or emit event to parent component
                    this.$emit('registration-success', response.data);
                }

            } catch (error) {
                console.error("Registration error:", error);

                if (error.response) {
                    this.errorMessage = error.response.data.error || 'Registration failed';
                    console.error("Server error:", error.response.data);
                } else if (error.request) {
                    this.errorMessage = 'Cannot connect to server. Please check if the backend is running.';
                    console.error("Network error:", error.request);
                } else {
                    this.errorMessage = 'An unexpected error occurred';
                    console.error("Error:", error.message);
                }
            } finally {
                this.isLoading = false;
            }
        },
    },
};
</script>
