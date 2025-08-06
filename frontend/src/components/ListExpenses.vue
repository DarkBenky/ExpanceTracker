<template>
    <div class="expenses-container">
        <h2 class="title">Expenses Dashboard</h2>
        <div v-if="!isAuthenticated" class="auth-message">
            <p>Please log in to view your expenses.</p>
        </div>
        <div v-else class="dashboard">
            <!-- Chart Controls -->
            <div class="chart-controls">
                <button @click="toggleCharts" class="toggle-btn">
                    {{ showCharts ? 'Hide Charts' : 'Show Charts' }}
                </button>
            </div>

            <!-- Charts Section -->
            <div v-if="showCharts" class="charts-section">
                <div class="chart-container">
                    <h3>Monthly Expenses Trend</h3>
                    <canvas id="monthlyChart"></canvas>
                </div>
                <div class="chart-container">
                    <h3>Overall Expenses by Category</h3>
                    <canvas id="categoryChart"></canvas>
                </div>
            </div>

            <!-- Monthly Breakdown -->
            <div class="months-grid">
                <div v-for="month in this.months" :key="month.name" class="month-card">

                    <div class="month-header">
                        <h3>{{ month.name }}</h3>
                        <span class="month-total">${{ month.total.toFixed(2) }}</span>
                        <button v-if="addExpense != month.name" @click="selectForm(month.name)">Add Expense</button>
                        <button v-else @click="selectForm(month.name)">Close Form</button>
                        <div v-if="addExpense == month.name">
                            form
                        </div>
                    </div>

                    <div v-if="month.expenses.length === 0" class="no-expenses">
                        <p>No expenses for this month</p>
                    </div>

                    <div v-else class="month-content">
                        <!-- Individual Month Pie Chart -->
                        <div class="month-chart-container">
                            <canvas :id="`monthChart${month.name}`" class="month-chart"></canvas>
                        </div>

                        <!-- Category Summary -->
                        <div class="category-summary">
                            <div v-for="(amount, category) in month.byCategory" :key="category" class="category-item">
                                <span class="category-name">{{ category }}</span>
                                <span class="category-amount">${{ amount.toFixed(2) }}</span>
                            </div>
                        </div>

                        <!-- Expense List -->
                        <div class="expenses-list">
                            <div v-for="expense in month.expenses" :key="expense.id" class="expense-item">
                                <div class="expense-main">
                                    <span class="expense-description">{{ expense.description }}</span>
                                    <span class="expense-amount">${{ expense.amount }}</span>
                                </div>
                                <div class="expense-meta">
                                    <span class="expense-category">{{ expense.category }}</span>
                                    <span class="expense-date">{{ formatDate(expense.date) }}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import { Chart, registerables } from 'chart.js';

Chart.register(...registerables);

export default {
    name: 'ListExpenses',
    data() {
        return {
            months: {
                1: { name: 'January', expenses: [], total: 0, byCategory: {} },
                2: { name: 'February', expenses: [], total: 0, byCategory: {} },
                3: { name: 'March', expenses: [], total: 0, byCategory: {} },
                4: { name: 'April', expenses: [], total: 0, byCategory: {} },
                5: { name: 'May', expenses: [], total: 0, byCategory: {} },
                6: { name: 'June', expenses: [], total: 0, byCategory: {} },
                7: { name: 'July', expenses: [], total: 0, byCategory: {} },
                8: { name: 'August', expenses: [], total: 0, byCategory: {} },
                9: { name: 'September', expenses: [], total: 0, byCategory: {} },
                10: { name: 'October', expenses: [], total: 0, byCategory: {} },
                11: { name: 'November', expenses: [], total: 0, byCategory: {} },
                12: { name: 'December', expenses: [], total: 0, byCategory: {} },
            },
            groupID: -1,
            expenses: [],
            isAuthenticated: false,
            token: null,
            monthlyChart: null,
            categoryChart: null,
            showCharts: true,
            monthCharts: {},
            addExpense: null,
        }
    },
    async created() {
        await this.checkExistingToken();
        if (this.isAuthenticated) {
            this.getExpenses();
        }
    },
    methods: {
        selectForm(monthName) {
            if (this.addExpense === monthName) {
                this.addExpense = null; // Deselect if already selected
            } else {
                this.addExpense = monthName; // Select the month for adding expense
            }
        },
        // async addExpense() {
        //     try {
        //         const response
        //     }
        // },
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
        calculateTotalByMonth() {
            Object.keys(this.months).forEach(month => {
                this.months[month].total = this.months[month].expenses.reduce((sum, expense) => sum + expense.amount, 0);
            });
        },
        calculateExpensesByCategory() {
            Object.keys(this.months).forEach(month => {
                this.months[month].byCategory = {};
                this.months[month].expenses.forEach(expense => {
                    if (!this.months[month].byCategory[expense.category]) {
                        this.months[month].byCategory[expense.category] = 0;
                    }
                    this.months[month].byCategory[expense.category] += expense.amount;
                });
            });
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
        formatDate(dateString) {
            return new Date(dateString).toLocaleDateString();
        },
        toggleCharts() {
            this.showCharts = !this.showCharts;
            if (this.showCharts) {
                this.$nextTick(() => {
                    this.createMonthlyChart();
                    this.createCategoryChart();
                });
            } else {
                // Destroy charts when hiding to free memory
                if (this.monthlyChart) {
                    this.monthlyChart.destroy();
                    this.monthlyChart = null;
                }
                if (this.categoryChart) {
                    this.categoryChart.destroy();
                    this.categoryChart = null;
                }
            }
        },
        calculateRegressionLine(data) {
            const n = data.length;
            const x = data.map((_, i) => i + 1); // months 1-12
            const y = data;

            const sumX = x.reduce((a, b) => a + b, 0);
            const sumY = y.reduce((a, b) => a + b, 0);
            const sumXY = x.reduce((sum, xi, i) => sum + xi * y[i], 0);
            const sumXX = x.reduce((sum, xi) => sum + xi * xi, 0);

            const slope = (n * sumXY - sumX * sumY) / (n * sumXX - sumX * sumX);
            const intercept = (sumY - slope * sumX) / n;

            return x.map(xi => slope * xi + intercept);
        },
        createCharts() {
            this.$nextTick(() => {
                if (this.showCharts) {
                    this.createMonthlyChart();
                    this.createCategoryChart();
                }
                this.createMonthlyPieCharts();
            });
        },
        createMonthlyChart() {
            const ctx = document.getElementById('monthlyChart');
            if (!ctx) return;

            if (this.monthlyChart) {
                this.monthlyChart.destroy();
            }

            const monthlyData = Object.values(this.months).map(month => month.total);
            const monthLabels = Object.values(this.months).map(month => month.name);
            const regressionData = this.calculateRegressionLine(monthlyData);

            this.monthlyChart = new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: monthLabels,
                    datasets: [{
                        label: 'Monthly Expenses',
                        data: monthlyData,
                        backgroundColor: 'rgba(99, 102, 241, 0.8)',
                        borderColor: 'rgba(99, 102, 241, 1)',
                        borderWidth: 1
                    }, {
                        label: 'Trend Line',
                        type: 'line',
                        data: regressionData,
                        borderColor: '#ef4444',
                        backgroundColor: 'rgba(239, 68, 68, 0.1)',
                        borderWidth: 3,
                        fill: false,
                        tension: 0.1,
                        pointRadius: 4,
                        pointBackgroundColor: '#ef4444'
                    }]
                },
                options: {
                    responsive: true,
                    plugins: {
                        legend: {
                            labels: { color: '#e5e7eb' }
                        }
                    },
                    scales: {
                        y: {
                            beginAtZero: true,
                            ticks: { color: '#e5e7eb' },
                            grid: { color: '#374151' }
                        },
                        x: {
                            ticks: { color: '#e5e7eb' },
                            grid: { color: '#374151' }
                        }
                    }
                }
            });
        },
        createCategoryChart() {
            const ctx = document.getElementById('categoryChart');
            if (!ctx) return;

            if (this.categoryChart) {
                this.categoryChart.destroy();
            }

            const categoryTotals = {};
            Object.values(this.months).forEach(month => {
                Object.entries(month.byCategory).forEach(([category, amount]) => {
                    categoryTotals[category] = (categoryTotals[category] || 0) + amount;
                });
            });

            const colors = ['#ef4444', '#f97316', '#eab308', '#22c55e', '#06b6d4', '#8b5cf6', '#ec4899'];

            this.categoryChart = new Chart(ctx, {
                type: 'doughnut',
                data: {
                    labels: Object.keys(categoryTotals),
                    datasets: [{
                        data: Object.values(categoryTotals),
                        backgroundColor: colors,
                        borderColor: '#1f2937',
                        borderWidth: 2
                    }]
                },
                options: {
                    responsive: true,
                    plugins: {
                        legend: {
                            labels: { color: '#e5e7eb' }
                        }
                    }
                }
            });
        },
        createMonthlyPieCharts() {
            // Destroy existing month charts
            Object.values(this.monthCharts).forEach(chart => chart.destroy());
            this.monthCharts = {};

            Object.values(this.months).forEach(month => {
                if (month.expenses.length > 0) {
                    this.$nextTick(() => {
                        const ctx = document.getElementById(`monthChart${month.name}`);
                        if (ctx) {
                            const colors = ['#ef4444', '#f97316', '#eab308', '#22c55e', '#06b6d4', '#8b5cf6', '#ec4899'];

                            this.monthCharts[month.name] = new Chart(ctx, {
                                type: 'pie',
                                data: {
                                    labels: Object.keys(month.byCategory),
                                    datasets: [{
                                        data: Object.values(month.byCategory),
                                        backgroundColor: colors,
                                        borderColor: '#1f2937',
                                        borderWidth: 2
                                    }]
                                },
                                options: {
                                    responsive: true,
                                    maintainAspectRatio: false,
                                    plugins: {
                                        legend: {
                                            position: 'bottom',
                                            labels: {
                                                color: '#e5e7eb',
                                                font: { size: 10 }
                                            }
                                        }
                                    }
                                }
                            });
                        }
                    });
                }
            });
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
                this.sortExpensesByMonth();
                this.calculateTotalByMonth();
                this.calculateExpensesByCategory();
                this.$nextTick(() => {
                    this.createCharts();
                });
                console.log("Expenses fetched and sorted by month:", this.months);
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

<style scoped>
.expenses-container {
    background: #111827;
    color: #e5e7eb;
    min-height: 100vh;
    padding: 2rem;
    font-family: 'Inter', sans-serif;
}

.title {
    color: #f9fafb;
    font-size: 2.5rem;
    font-weight: 700;
    margin-bottom: 2rem;
    text-align: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
}

.auth-message {
    text-align: center;
    padding: 4rem;
    background: #1f2937;
    border-radius: 12px;
    border: 1px solid #374151;
}

.dashboard {
    max-width: 1400px;
    margin: 0 auto;
}

.chart-controls {
    text-align: center;
    margin-bottom: 2rem;
}

.toggle-btn {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 0.75rem 2rem;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: transform 0.2s, box-shadow 0.2s;
}

.toggle-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.charts-section {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    margin-bottom: 3rem;
}

.chart-container {
    background: #1f2937;
    border-radius: 12px;
    padding: 1.5rem;
    border: 1px solid #374151;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.chart-container h3 {
    color: #f9fafb;
    margin-bottom: 1rem;
    font-size: 1.25rem;
    font-weight: 600;
}

.months-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    gap: 1.5rem;
}

.month-card {
    background: #1f2937;
    border-radius: 12px;
    border: 1px solid #374151;
    overflow: hidden;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s, box-shadow 0.2s;
}

.month-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 25px -3px rgba(0, 0, 0, 0.2);
}

.month-header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 1rem 1.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.month-header h3 {
    color: white;
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
}

.month-total {
    color: white;
    font-size: 1.5rem;
    font-weight: 700;
}

.month-content {
    padding: 1.5rem;
}

.no-expenses {
    padding: 2rem 1.5rem;
    text-align: center;
    color: #9ca3af;
    font-style: italic;
}

.month-chart-container {
    height: 200px;
    margin-bottom: 1.5rem;
    display: flex;
    justify-content: center;
    align-items: center;
}

.category-summary {
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #374151;
}

.category-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0;
}

.category-name {
    color: #d1d5db;
    font-weight: 500;
}

.category-amount {
    color: #10b981;
    font-weight: 600;
}

.expenses-list {
    space-y: 0.75rem;
}

.expense-item {
    background: #374151;
    border-radius: 8px;
    padding: 1rem;
    margin-bottom: 0.75rem;
    border-left: 4px solid #6366f1;
}

.expense-main {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
}

.expense-description {
    color: #f9fafb;
    font-weight: 600;
}

.expense-amount {
    color: #10b981;
    font-weight: 700;
    font-size: 1.1rem;
}

.expense-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.875rem;
}

.expense-category {
    background: #4f46e5;
    color: white;
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 500;
}

.expense-date {
    color: #9ca3af;
}

@media (max-width: 768px) {
    .charts-section {
        grid-template-columns: 1fr;
    }

    .months-grid {
        grid-template-columns: 1fr;
    }

    .expenses-container {
        padding: 1rem;
    }
}
</style>
