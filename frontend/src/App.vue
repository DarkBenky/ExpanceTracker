<template>
  <div class="app">
    <div v-if="isAuthenticated">
      <ListExpenses/>
      <button @click="clearAuth" class="btn btn-secondary">Logout</button>
    </div>
    <div v-else>
      <LoginRegister 
        @login-success="handleLoginSuccess" 
        @registration-success="handleRegistrationSuccess"
        @token-validated="handleTokenValidated"
      />
    </div>
  </div>
</template>

<script>
import LoginRegister from "./components/LoginRegister.vue";
import ListExpenses from "./components/ListExpenses.vue";
import axios from "axios";

export default {
  name: "App",
  components: {
    LoginRegister,
    ListExpenses,
  },
  data() {
    return {
      isAuthenticated: false,
      token: null,
      user: null,
    };
  },
  async created() {
    await this.checkExistingToken();
  },
  methods: {
    async checkExistingToken() {
      const token = localStorage.getItem("authToken");
      const expiry = localStorage.getItem("tokenExpiry");
      
      if (token && expiry && Date.now() < parseInt(expiry)) {
        console.log("Token found in localStorage:", token);
        try {
          const response = await axios.post(`${this.$apiUrl}validate/token`, {
            token: token,
          });
          
          if (response.data.valid) {
            console.log("Token is valid, user is authenticated");
            this.isAuthenticated = true;
            this.token = token;
            this.user = response.data.user || { id: response.data.user_id };
          } else {
            this.clearAuth();
          }
        } catch (error) {
          console.error("Token validation failed:", error);
          this.clearAuth();
        }
      } else {
        this.clearAuth();
      }
    },
    
    handleLoginSuccess(data) {
      console.log("Login successful in App.vue:", data);
      this.isAuthenticated = true;
      this.token = data.token;
      this.user = data.user || { id: data.user_id };
    },
    
    handleRegistrationSuccess(data) {
      console.log("Registration successful in App.vue:", data);
      this.isAuthenticated = true;
      this.token = data.token;
      this.user = data.user || { id: data.user_id };
    },
    
    handleTokenValidated(data) {
      console.log("Token validated in App.vue:", data);
      this.isAuthenticated = true;
      this.user = data.user || { id: data.user_id };
    },
    
    clearAuth() {
      this.isAuthenticated = false;
      this.token = null;
      this.user = null;
      localStorage.removeItem("authToken");
      localStorage.removeItem("tokenExpiry");
    }
  },
};
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  background-color: #121212;
  color: #ffffff;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.app {
  min-height: 100vh;
  background-color: #121212;
  color: #ffffff;
}

#app {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: #121212;
  color: #ffffff;
  min-height: 100vh;
}
</style>
