<template>
    <div>
        <h2>Expenses</h2>
        <div v-if="!isAuthenticated">
            <p>Please log in to view your expenses.</p>
        </div>
        <div v-else>
            <div v-for="month in this.months" :key="month">
                <h3>{{ month.name }}</h3>
                <div v-if="month.expenses.length === 0">
                    <p>No expenses for this month.</p>
                </div>
                <div v-else>
                    <ul>
                        <li v-for="expense in month.expenses" :key="expense.id">
                            {{ expense.description }} - ${{ expense.amount }}
                        </li>
                    </ul>
                </div>
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
            months: {
                1: {name: 'January', expenses : []},
                2: {name: 'February', expenses : []},
                3: {name: 'March', expenses : []},
                4: {name: 'April', expenses : []},
                5: {name: 'May', expenses : []},
                6: {name: 'June', expenses : []},
                7: {name: 'July', expenses : []},
                8: {name: 'August', expenses : []},
                9: {name: 'September', expenses : []},
                10: {name: 'October', expenses : []},
                11: {name: 'November', expenses : []},
                12: {name: 'December', expenses : []},
            },
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
        sortExpensesByMonth() {
            // Clear existing expenses from months
            Object.values(this.months).forEach(month => month.expenses = []);
            
            for (const expense of this.expenses) {
                const month = new Date(expense.date).getMonth() + 1; // getMonth() returns 0-11
                console.log("Sorting expense:", expense, "into month:", month);
                if (this.months[month]) {
                    this.months[month].expenses.push(expense);
                }
                console.log(this.months[month].expenses);
            }
        },
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
                this.sortExpensesByMonth(); // Call sorting after expenses are loaded
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
