<template>
    <div>
        <h2>Expenses</h2>
        <div v-if="!isAuthenticated">
            <p>Please log in to view your expenses.</p>
        </div>
        <div v-else>
            <div v-for="expense in this.expenses" :key="expense.id">
                <p>{{ expense.description }} - ${{ expense.amount }}</p>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    name: 'ListExpenses',
    data() {
        return {
            groupID: -1,
            expenses: [],
            isAuthenticated: false,
            token: null,
        }
    },
    async created() {
        await this.checkExistingToken();
        if (this.isAuthenticated) {
            this.getExpenses();
        }
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
        getExpenses() {
            this.isLoading = true;
            this.errorMessage = '';

            axios.post(`${this.$apiUrl}expenses/get`, {
                token: this.token,
                group_id: this.groupID,
            }
            ).then(response => {
                this.expenses = response.data;
                this.isLoading = false;
            }).catch(error => {
                console.error("Error fetching expenses:", error);
                this.errorMessage = "Failed to load expenses.";
                this.isLoading = false;
            });
        },
    }
}
</script>
