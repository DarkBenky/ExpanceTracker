<template>
    <div>
        <h2>Expenses</h2>
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
            isLoading: false,
            errorMessage: ''
        }
    },
    props: {
        token: {
            type: String,
            required: true
        },
        tokenVerified: {
            type: Boolean,
            required: true
        }
    },
    methods: {
        getExpenses() {
            const response = axios.get(`${this.$apiUrl}expenses/get`, {
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

