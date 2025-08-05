<template>
    <div class="login-container">
        <div class="login-card">
            <div v-if="LoginOrRegister === 'Login'">
                <h1 class="title">Sign In</h1>
                <form @submit.prevent="handleLogin" class="form">
                    <div v-if="errorMessage" class="error-message">
                        {{ errorMessage }}
                    </div>
                    <div class="input-group">
                        <label for="username">Username</label>
                        <input 
                            type="text" 
                            id="username" 
                            v-model="loginForm.username" 
                            required 
                            :disabled="isLoading"
                        />
                    </div>
                    <div class="input-group">
                        <label for="password">Password</label>
                        <input 
                            type="password" 
                            id="password" 
                            v-model="loginForm.password" 
                            required 
                            :disabled="isLoading"
                        />
                    </div>
                    <button type="submit" class="btn btn-primary" :disabled="isLoading">
                        {{ isLoading ? 'Signing In...' : 'Sign In' }}
                    </button>
                </form>
                <p class="switch-text">Don't have an account?</p>
                <button @click="log" class="btn btn-secondary">Create Account</button>
            </div>
            
            <div v-else>
                <h1 class="title">Create Account</h1>
                <form @submit.prevent="handleRegister" class="form">
                    <div v-if="errorMessage" class="error-message">
                        {{ errorMessage }}
                    </div>
                    <div class="input-group">
                        <label for="new-username">Username</label>
                        <input 
                            type="text" 
                            id="new-username" 
                            v-model="registerForm.username" 
                            required 
                            :disabled="isLoading"
                        />
                    </div>
                    <div class="input-group">
                        <label for="new-password">Password</label>
                        <input 
                            type="password" 
                            id="new-password" 
                            v-model="registerForm.password" 
                            required 
                            :disabled="isLoading"
                        />
                    </div>
                    <button type="submit" class="btn btn-primary" :disabled="isLoading">
                        {{ isLoading ? 'Creating Account...' : 'Create Account' }}
                    </button>
                </form>
                <p class="switch-text">Already have an account?</p>
                <button @click="log" class="btn btn-secondary">Sign In</button>
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
    async created() {
        const token = localStorage.getItem('authToken');
        if (token) {
            try {
                const response = await axios.post(`${this.$apiUrl}validate/token`, {
                    token: token
                });
                if (response.data.valid) {
                    this.$emit('token-validated', response.data);
                }
            } catch (error) {
                console.error("Token validation failed:", error);
                localStorage.removeItem('authToken');
                localStorage.removeItem('tokenExpiry');
            }
        }
    },
    methods: {
        log() {
            this.errorMessage = ''; 
            this.LoginOrRegister = this.LoginOrRegister === "Login" ? "Register" : "Login";
        },
        
        async handleLogin() {
            this.isLoading = true;
            this.errorMessage = '';
            
            try {
                const response = await axios.post(`${this.$apiUrl}login`, {
                    username: this.loginForm.username,
                    password: this.loginForm.password
                });

                if (response.data.token) {
                    localStorage.setItem('authToken', response.data.token);
                    localStorage.setItem('tokenExpiry', Date.now() + (24 * 60 * 60 * 1000)); // 24 hours
                    
                    this.$emit('login-success', response.data);
                    
                    // Clear form
                    this.loginForm.username = '';
                    this.loginForm.password = '';
                }
                
            } catch (error) {
                if (error.response) {
                    this.errorMessage = error.response.data.error || 'Login failed';
                } else if (error.request) {
                    this.errorMessage = 'Cannot connect to server. Please check if the backend is running.';
                } else {
                    this.errorMessage = 'An unexpected error occurred';
                }
            } finally {
                this.isLoading = false;
            }
        },
        
        async handleRegister() {
            this.isLoading = true;
            this.errorMessage = '';
            
            try {
                const response = await axios.post(`${this.$apiUrl}register`, {
                    username: this.registerForm.username,
                    password: this.registerForm.password
                });

                if (response.data.token) {
                    localStorage.setItem('authToken', response.data.token);
                    localStorage.setItem('tokenExpiry', Date.now() + (24 * 60 * 60 * 1000)); // 24 hours
                    
                    this.$emit('registration-success', response.data);
                    
                    // Clear form
                    this.registerForm.username = '';
                    this.registerForm.password = '';
                }

            } catch (error) {
                if (error.response) {
                    this.errorMessage = error.response.data.error || 'Registration failed';
                } else if (error.request) {
                    this.errorMessage = 'Cannot connect to server. Please check if the backend is running.';
                } else {
                    this.errorMessage = 'An unexpected error occurred';
                }
            } finally {
                this.isLoading = false;
            }
        },
    },
};
</script>

<style scoped>
.login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #1e1e1e 0%, #2d2d2d 100%);
    padding: 20px;
}

.login-card {
    background: #1a1a1a;
    border-radius: 12px;
    padding: 40px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
    border: 1px solid #333;
    max-width: 400px;
    width: 100%;
}

.title {
    color: #ffffff;
    text-align: center;
    margin-bottom: 30px;
    font-size: 28px;
    font-weight: 600;
}

.form {
    margin-bottom: 20px;
}

.input-group {
    margin-bottom: 20px;
}

.input-group label {
    display: block;
    color: #b3b3b3;
    margin-bottom: 8px;
    font-size: 14px;
    font-weight: 500;
}

.input-group input {
    width: 100%;
    padding: 12px 16px;
    background: #2a2a2a;
    border: 1px solid #404040;
    border-radius: 8px;
    color: #ffffff;
    font-size: 16px;
    transition: all 0.3s ease;
}

.input-group input:focus {
    outline: none;
    border-color: #4285f4;
    box-shadow: 0 0 0 2px rgba(66, 133, 244, 0.2);
}

.input-group input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn {
    width: 100%;
    padding: 12px 16px;
    border: none;
    border-radius: 8px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    margin-bottom: 10px;
}

.btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-primary {
    background: linear-gradient(135deg, #4285f4 0%, #34a853 100%);
    color: white;
}

.btn-primary:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 10px 20px rgba(66, 133, 244, 0.3);
}

.btn-secondary {
    background: transparent;
    color: #4285f4;
    border: 1px solid #4285f4;
}

.btn-secondary:hover:not(:disabled) {
    background: #4285f4;
    color: white;
}

.switch-text {
    text-align: center;
    color: #b3b3b3;
    margin: 20px 0 10px 0;
    font-size: 14px;
}

.error-message {
    background: #ff4444;
    color: white;
    padding: 12px;
    border-radius: 8px;
    margin-bottom: 20px;
    text-align: center;
    font-size: 14px;
}
</style>
