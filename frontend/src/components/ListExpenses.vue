<template>
    <div class="expenses-container">
        <div class="groups-section">
            <div class="groups-header">
                <h3>Groups</h3>
                <button @click="createGroup = !createGroup" class="btn btn-primary">
                    {{ createGroup ? 'Cancel' : 'Create Group' }}
                </button>
            </div>
            <div v-if="createGroup == true" class="create-group-form">
                <form @submit.prevent="createNewGroup" class="group-form">
                    <div class="input-group">
                        <input type="text" v-model="newGroupName" placeholder="Enter group name" required />
                    </div>
                    <div class="form-buttons">
                        <button @click="addNewGroup(newGroupName)" type="submit" class="btn btn-primary">Create</button>
                        <button type="button" @click="createGroup = false" class="btn btn-secondary">Cancel</button>
                    </div>
                </form>
            </div>
            <div v-for="group in groups" :key="group.id" class="group-item">
                <span>{{ group.name }}</span>
                <button @click="selectGroup(group.id)" :class="{ active: groupID === group.id }">
                    {{ groupID === group.id ? 'Selected' : 'Select' }}
                </button>
                <div v-if="groupID === group.id">
                    <div class="group-members">
                        <p class="members-title">Members:</p>
                        <ul v-if="users.length > 0" class="members-list">
                            <li v-for="user in users" :key="user.id" class="member-item">
                                <span class="member-avatar">{{ user.username.charAt(0).toUpperCase() }}</span>
                                <span class="member-name">{{ user.username }}</span>
                            </li>
                        </ul>
                        <p v-else class="no-members">No members found</p>
                        <button @click="addMember = !addMember" class="btn btn-secondary add-member-btn">
                            {{ addMember ? 'Cancel' : 'Add Member' }}
                        </button>
                        <div v-if="addMember" class="add-member-form">
                            <div class="input-group">
                                <select v-model="selectedUser" class="member-select">
                                    <option value="" disabled>Select User</option>
                                    <option v-for="user in this.usersNotInGroup()" :key="user.id" :value="user.id">{{
                                        user.username }}
                                    </option>
                                </select>
                            </div>
                            <div class="form-buttons">
                                <button @click="addUserToGroup" class="btn btn-primary">Add User</button>
                                <button @click="addMember = false" class="btn btn-secondary">Cancel</button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="!isAuthenticated" class="auth-message">
            <p>Please log in to view your expenses.</p>
        </div>
        <div v-else class="dashboard">
            <!-- Controls -->
            <div class="controls">
                <div class="chart-controls">
                    <button @click="toggleCharts" class="toggle-btn">
                        {{ showCharts ? 'Hide Charts' : 'Show Charts' }}
                    </button>
                </div>
                <div class="refresh-controls">
                    <button @click="manualRefresh" class="refresh-btn" :disabled="isLoading">
                        {{ isLoading ? 'Loading...' : 'Refresh Now' }}
                    </button>
                    <button @click="toggleAutoRefresh" class="auto-refresh-btn" :class="{ active: autoRefresh }">
                        Auto-refresh: {{ autoRefresh ? 'ON' : 'OFF' }}
                    </button>
                </div>
            </div>

            <!-- Charts Section -->
            <div v-if="showCharts" class="charts-section" :key="'charts-' + groupRenderKey">
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
                <div v-for="month in months" :key="month.name + '-' + groupRenderKey" class="month-card">

                    <div class="month-header">
                        <h3>{{ month.name }} </h3>
                        <span class="month-total">${{ month.total.toFixed(2) }}</span>
                        <button v-if="addExpense != month.name" @click="selectForm(month.name)"
                            class="btn btn-primary add-expense-btn">Add Expense</button>
                        <button v-else @click="selectForm(month.name)" class="btn btn-secondary add-expense-btn">Close
                            Form</button>
                    </div>

                    <div v-if="addExpense == month.name" class="add-expense-form">
                        <form @submit.prevent="addExpenseToMonth(month.name)" class="expense-form">
                            <div class="input-group">
                                <label for="description">Description</label>
                                <input type="text" v-model="newExpense.description" required />
                            </div>
                            <div class="input-group">
                                <label for="amount">Amount</label>
                                <input type="number" step="0.01" v-model.number="newExpense.amount" required />
                            </div>
                            <div class="input-group">
                                <label for="category">Category</label>
                                <select v-model="newExpense.category" required>
                                    <option value="" disabled>Select Category</option>
                                    <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
                                </select>
                            </div>
                            <div class="form-buttons">
                                <button type="submit" class="btn btn-primary">Add Expense</button>
                                <button type="button" @click="addExpense = null"
                                    class="btn btn-secondary">Cancel</button>
                            </div>
                        </form>
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
                                <div class="expense-actions">
                                    <button @click="deleteExpense(expense.id)"
                                        class="btn btn-danger expense-action-btn">Delete</button>
                                    <button v-if="editExpenseId !== expense.id" @click="editExpenseId = expense.id"
                                        class="btn btn-secondary expense-action-btn">Edit</button>
                                    <button v-else @click="editExpenseId = null"
                                        class="btn btn-secondary expense-action-btn">Cancel</button>
                                </div>
                                <form v-if="this.editExpenseId == expense.id" class="edit-expense-form">
                                    <div class="input-group">
                                        <label for="editDescription">Description</label>
                                        <input type="text" v-model="expense.description" required />
                                    </div>
                                    <div class="input-group">
                                        <label for="editAmount">Amount</label>
                                        <input type="number" step="0.01" v-model.number="expense.amount" required />
                                    </div>
                                    <div class="input-group">
                                        <label for="editCategory">Category</label>
                                        <select v-model="expense.category" required>
                                            <option value="" disabled>Select Category</option>
                                            <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
                                        </select>
                                    </div>
                                    <div class="input-group">
                                        <label for="editDate">Date</label>
                                        <div class="date-selectors">
                                            <select v-if="editDate.month != '' && editDate.year != ''"
                                                v-model="editDate.day" @change="updateExpenseDate(expense)" required>
                                                <option value="" disabled>Day</option>
                                                <option v-for="day in getDaysInEditMonth(editDate.month, editDate.year)"
                                                    :key="day" :value="day">
                                                    {{ day }}
                                                </option>
                                            </select>
                                            <select v-model="editDate.month" @change="updateExpenseDate(expense)"
                                                required>
                                                <option value="" disabled>Month</option>
                                                <option v-for="month in getMonthOptions()" :key="month.value"
                                                    :value="month.value">
                                                    {{ month.label }}
                                                </option>
                                            </select>
                                            <select v-model="editDate.year" @change="updateExpenseDate(expense)"
                                                required>
                                                <option value="" disabled>Year</option>
                                                <option v-for="year in getYearOptions()" :key="year" :value="year">
                                                    {{ year }}
                                                </option>
                                            </select>
                                        </div>
                                    </div>
                                    <div class="form-buttons">
                                        <button type="submit" @click.prevent="updateExpense(expense)"
                                            class="btn btn-primary">Update Expense</button>
                                        <button type="button" @click="this.editExpenseId = null"
                                            class="btn btn-secondary">Cancel</button>
                                    </div>
                                </form>
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
            allUsers: [],
            users: {},
            addMember: false,
            createGroup: false,
            editExpenseId: null,
            groups: [],
            groupID: -1,
            expenses: [],
            isAuthenticated: false,
            token: null,
            monthlyChart: null,
            categoryChart: null,
            showCharts: true,
            monthCharts: {},
            addExpense: null,
            refreshInterval: null,
            autoRefresh: true,
            refreshIntervalTime: 30000, // 30 seconds
            newExpense: {
                description: '',
                amount: 0,
                category: '',
                date: ''
            },
            editDate: {
                day: '',
                month: '',
                year: ''
            },
            categories: [
                'Jedlo - Food & Dining',
                'Cestovanie - Transportation', 
                'Nakupy - Shopping',
                'Zabava - Entertainment',
                'Účty a služby - Bills & Utilities',
                'Zdravotníctvo - Healthcare',
                'Cestovanie - Travel',
                'Vzdelávanie - Education',
                'Bývanie - Housing',
                'Oblečenie - Clothing',
                'Šport a fitness - Sports & Fitness',
                'Darčeky - Gifts',
                'Elektronika - Electronics',
                'Automobil - Car & Transport',
                'Investície - Investments',
                'Poistenie - Insurance',
                'Krása - Beauty & Personal Care',
                'Domáce potreby - Household Items',
                'Opravy - Repairs & Maintenance',
                'Ostatné - Other'
            ],
            isLoading: false,
            errorMessage: '',
            groupRenderKey: 0, // force canvas re-mount on group change
        }
    },
    async created() {
        await this.checkExistingToken();
        if (this.isAuthenticated) {
            this.getExpenses();
            this.GetGroups();
            this.getAllUsers();
            this.startPeriodicRefresh();
        }
    },
    beforeUnmount() {
        this.stopPeriodicRefresh();
        this.destroyAllCharts();
    },
    watch: {
        groupID(newGroupId) {
            if (newGroupId !== -1 && this.isAuthenticated) {
                this.destroyAllCharts();
                this.groupRenderKey++; // force fresh canvases
                this.getExpenses();
            }
        },
        isAuthenticated(newAuth) {
            if (newAuth) {
                this.startPeriodicRefresh();
            } else {
                this.stopPeriodicRefresh();
            }
        }
    },
    methods: {
        addNewGroup(groupName) {
            if (!groupName) {
                console.log('Please enter a group name');
                return;
            }
            axios.post(`${this.$apiUrl}groups`, {
                token: this.token,
                name: groupName
            }).then(response => {
                if (response.status === 201) {
                    this.createGroup = false;
                    this.newGroupName = '';
                    this.GetGroups();
                    console.log('Group created successfully:', response.data);
                }
            }).catch(error => {
                console.error('Error creating group:', error);
                console.log('Failed to create group. Please try again.');
            });
        },
        addUserToGroup() {
            if (!this.selectedUser) {
                console.log('Please select a user to add');
                return;
            }

            axios.post(`${this.$apiUrl}groups/members/add`, {
                token: this.token,
                group_id: this.groupID,
                user_id: this.selectedUser
            }).then(response => {
                if (response.status === 200) {
                    this.getUsersFromGroup(this.groupID);
                    this.selectedUser = null; // Reset selection
                }
            }).catch(error => {
                console.error('Error adding user to group:', error);
                console.log('Failed to add user to group. Please try again.');
            });
        },
        usersNotInGroup() {
            return this.allUsers.filter(user => !this.users.some(u => u.id === user.id));
        },
        getAllUsers() {
            axios.post(`${this.$apiUrl}users/get`, {
                token: this.token
            }).then(response => {
                if (response.status === 200) {
                    // Try accessing the data directly first
                    this.allUsers = response.data || [];
                    console.log('All users fetched:', this.allUsers);
                    console.log('Full response data:', response.data);
                }
            }).catch(error => {
                console.error('Error fetching all users:', error);
            });
        },
        selectGroup(groupId) {
            this.groupID = groupId;
            this.getUsersFromGroup(groupId);
        },
        getUsersFromGroup(groupId) {
            // Fetch users from the selected group
            axios.post(`${this.$apiUrl}groups/members/get`, {
                token: this.token,
                group_id: groupId
            }).then(response => {
                if (response.status === 200) {
                    this.users = response.data.users || [];
                    console.log('Group users fetched:', this.users);
                }
            }).catch(error => {
                console.error('Error fetching group users:', error);
            });
        },
        updateExpense(expense) {
            if (!expense.description || !expense.amount || !expense.category) {
                console.log('Please fill in all fields');
                return;
            }

            // Ensure date is properly formatted
            let formattedDate = expense.date;
            if (this.editDate.day && this.editDate.month && this.editDate.year) {
                // If date was edited using the selectors, format it properly
                const day = String(this.editDate.day).padStart(2, '0');
                const month = String(this.editDate.month).padStart(2, '0');
                formattedDate = `${this.editDate.year}-${month}-${day}`;
            }

            // Update the expense in the database
            axios.post(`${this.$apiUrl}expenses/update`, {
                token: this.token,
                group_id: this.groupID,
                expense_id: expense.id,
                description: expense.description,
                amount: parseFloat(expense.amount), // Ensure amount is a number
                category: expense.category,
                date: formattedDate // Ensure date is in YYYY-MM-DD format
            }).then(response => {
                if (response.status === 200) {
                    this.editExpenseId = null; // Close edit form
                    // Reset edit date selectors
                    this.editDate = {
                        day: '',
                        month: '',
                        year: ''
                    };
                    console.log('Expense updated successfully!');
                }
            }).catch(error => {
                console.error('Error updating expense:', error);
                if (error.response && error.response.data && error.response.data.error) {
                    console.log('Error: ' + error.response.data.error);
                } else {
                    console.log('Failed to update expense. Please try again.');
                }
            });
        },

        // Add this method to properly initialize edit form with current expense data
        editExpense(expense) {
            this.editExpenseId = expense.id;

            // Initialize date selectors with current expense date
            if (expense.date) {
                const date = new Date(expense.date);
                this.editDate = {
                    day: date.getDate(),
                    month: date.getMonth() + 1, // getMonth() returns 0-11
                    year: date.getFullYear()
                };
            } else {
                this.editDate = {
                    day: '',
                    month: '',
                    year: ''
                };
            }
        },

        getDaysInMonth(monthName) {
            const currentYear = new Date().getFullYear();
            const monthNumber = Object.keys(this.months).find(key =>
                this.months[key].name === monthName
            );

            // Get the number of days in this month
            const daysInMonth = new Date(currentYear, parseInt(monthNumber), 0).getDate();

            const days = [];
            for (let day = 1; day <= daysInMonth; day++) {
                const dateValue = `${currentYear}-${monthNumber.toString().padStart(2, '0')}-${day.toString().padStart(2, '0')}`;
                const dateLabel = `${monthName} ${day}, ${currentYear}`;
                days.push({
                    value: dateValue,
                    label: dateLabel
                });
            }

            return days;
        },
        deleteExpense(expenseId) {
            console.log('Deleting expense with ID:', expenseId);
            axios.delete(`${this.$apiUrl}expenses`, {
                headers: { 'Content-Type': 'application/json' },
                data: {
                    token: this.token,
                    group_id: this.groupID,
                    expense_id: expenseId
                }
            }).then(response => {
                if (response.status === 200) {
                    this.getExpenses();
                }
            }).catch(error => {
                console.error('Error deleting expense:', error?.response?.data || error);
            });
        },
        destroyAllCharts() {
            if (this.monthlyChart && typeof this.monthlyChart.destroy === 'function') {
                this.monthlyChart.destroy();
                this.monthlyChart = null;
            }
            if (this.categoryChart && typeof this.categoryChart.destroy === 'function') {
                this.categoryChart.destroy();
                this.categoryChart = null;
            }
            Object.values(this.monthCharts).forEach(c => {
                if (c && typeof c.destroy === 'function') c.destroy();
            });
            this.monthCharts = {};
        },
        startPeriodicRefresh() {
            if (this.refreshInterval) {
                clearInterval(this.refreshInterval);
            }

            this.refreshInterval = setInterval(() => {
                if (this.autoRefresh && this.isAuthenticated && this.groupID !== -1) {
                    console.log('Auto-refreshing expenses data...');
                    this.getExpenses();
                }
            }, this.refreshIntervalTime);
        },
        stopPeriodicRefresh() {
            if (this.refreshInterval) {
                clearInterval(this.refreshInterval);
                this.refreshInterval = null;
            }
        },
        toggleAutoRefresh() {
            this.autoRefresh = !this.autoRefresh;
            if (this.autoRefresh) {
                this.startPeriodicRefresh();
            } else {
                this.stopPeriodicRefresh();
            }
        },
        manualRefresh() {
            if (this.isAuthenticated && this.groupID !== -1) {
                this.getExpenses();
            }
        },
        selectForm(monthName) {
            if (this.addExpense === monthName) {
                this.addExpense = null; // Deselect if already selected
            } else {
                this.addExpense = monthName; // Select the month for adding expense
            }
        },
        async addExpenseToMonth(monthName) {
            if (!this.newExpense.description || !this.newExpense.amount || !this.newExpense.category) {
                console.log('Please fill in all fields');
                return;
            }

            // Get the month number from month name
            const monthNumber = Object.keys(this.months).find(key =>
                this.months[key].name === monthName
            );

            // Create date for the selected month (using first day of the month)
            const currentYear = new Date().getFullYear();
            const expenseDate = new Date(currentYear, parseInt(monthNumber), 1);

            try {
                const response = await axios.post(`${this.$apiUrl}expenses`, {
                    token: this.token,
                    group_id: this.groupID,
                    description: this.newExpense.description,
                    amount: this.newExpense.amount,
                    category: this.newExpense.category,
                    date: expenseDate.toISOString().split('T')[0] // Format as YYYY-MM-DD
                });

                if (response.status === 200) {
                    // Reset form
                    this.newExpense = {
                        description: '',
                        amount: 0,
                        category: '',
                        date: ''
                    };

                    // Close form
                    this.addExpense = null;

                    // Refresh expenses to show the new one
                    await this.getExpenses();

                    console.log('Expense added successfully!');
                }
            } catch (error) {
                console.error('Error adding expense:', error);
                if (error.response && error.response.data && error.response.data.error) {
                    console.log('Error: ' + error.response.data.error);
                } else {
                    console.log('Failed to add expense. Please try again.');
                }
            }
        },
        async GetGroups() {
            try {
                const response = await axios.post(`${this.$apiUrl}groups/get`, {
                    token: this.token,
                });

                // The API returns the groups array directly, not nested under 'groups'
                if (response.data && Array.isArray(response.data)) {
                    this.groups = response.data;
                    console.log("Groups fetched successfully:", this.groups);
                } else {
                    console.error("No groups found in response:", response.data);
                }
            } catch (error) {
                console.error("Error fetching groups:", error);
            }
        },
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
                this.groupRenderKey++; // remount canvases when showing again
                this.$nextTick(() => {
                    this.createMonthlyChart();
                    this.createCategoryChart();
                    this.createMonthlyPieCharts();
                });
            } else {
                this.destroyAllCharts();
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
            if (!this.showCharts) return;
            const canvas = document.getElementById('monthlyChart');
            if (!canvas) return;
            if (this.monthlyChart) this.monthlyChart.destroy();

            const monthlyData = Object.values(this.months).map(month => month.total);
            const monthLabels = Object.values(this.months).map(month => month.name);
            const regressionData = this.calculateRegressionLine(monthlyData);

            this.monthlyChart = new Chart(canvas, {
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
            if (!this.showCharts) return;
            const canvas = document.getElementById('categoryChart');
            if (!canvas) return;
            if (this.categoryChart) this.categoryChart.destroy();

            const categoryTotals = {};
            Object.values(this.months).forEach(month => {
                Object.entries(month.byCategory).forEach(([category, amount]) => {
                    categoryTotals[category] = (categoryTotals[category] || 0) + amount;
                });
            });

            const colors = ['#ef4444', '#f97316', '#eab308', '#22c55e', '#06b6d4', '#8b5cf6', '#ec4899'];

            this.categoryChart = new Chart(canvas, {
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
            }).then(response => {
                this.expenses = response.data;
                this.sortExpensesByMonth();
                this.calculateTotalByMonth();
                this.calculateExpensesByCategory();
                this.getAllUsers();
                this.$nextTick(() => {
                    // canvases just remounted due to groupRenderKey increment
                    if (this.showCharts) {
                        this.createMonthlyChart();
                        this.createCategoryChart();
                    }
                    this.createMonthlyPieCharts();
                });
                this.isLoading = false;
            }).catch(error => {
                console.error("Error fetching expenses:", error);
                this.errorMessage = "Failed to load expenses.";
                this.isLoading = false;
            });
        },
        getMonthOptions() {
            return Object.keys(this.months).map(key => ({
                value: key,
                label: this.months[key].name
            }));
        },
        getYearOptions() {
            const currentYear = new Date().getFullYear();
            return Array.from({ length: 10 }, (_, i) => currentYear - i);
        },
        getDaysInEditMonth(month, year) {
            if (!month || !year) return [];
            const daysInMonth = new Date(year, month, 0).getDate();
            return Array.from({ length: daysInMonth }, (_, i) => i + 1);
        },
        updateExpenseDate(expense) {
            const { day, month, year } = this.editDate;
            if (day && month && year) {
                const date = new Date(year, month - 1, day);
                expense.date = date.toISOString().split('T')[0]; // Update expense date
            }
        }
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

.controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
}

.chart-controls {
    text-align: center;
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

.refresh-controls {
    display: flex;
    gap: 1rem;
}

.refresh-btn,
.auto-refresh-btn {
    background: #374151;
    color: #e5e7eb;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s, transform 0.2s;
}

.refresh-btn:disabled {
    background: #6b7280;
    cursor: not-allowed;
}

.auto-refresh-btn.active {
    background: #22c55e;
    color: white;
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

.groups-section {
    margin-bottom: 2rem;
}

.group-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    background: #1f2937;
    border-radius: 8px;
    margin-bottom: 0.75rem;
    border: 1px solid transparent;
    transition: border-color 0.2s;
}

.group-item:hover {
    border-color: #374151;
}

.group-item span {
    color: #d1d5db;
    font-weight: 500;
}

.group-item button {
    background: #374151;
    color: #e5e7eb;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s, transform 0.2s;
}

.group-item button.active {
    background: #6366f1;
    color: white;
}

.group-item button:hover {
    transform: translateY(-2px);
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

.input-group {
    margin-bottom: 1rem;
}

.input-group label {
    display: block;
    color: #d1d5db;
    font-weight: 500;
    margin-bottom: 0.5rem;
}

.input-group input,
.input-group select {
    width: 100%;
    background: #1f2937;
    border: 1px solid #4b5563;
    border-radius: 6px;
    padding: 0.75rem;
    color: #e5e7eb;
    font-size: 1rem;
}

.input-group input:focus,
.input-group select:focus {
    outline: none;
    border-color: #6366f1;
    box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.form-buttons {
    display: flex;
    gap: 1rem;
    margin-top: 1.5rem;
}

.btn {
    padding: 0.75rem 1.5rem;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    border: none;
}

.btn-primary {
    background: #6366f1;
    color: white;
}

.btn-primary:hover {
    background: #5b57d9;
    transform: translateY(-1px);
}

.btn-secondary {
    background: #6b7280;
    color: white;
}

.btn-secondary:hover {
    background: #5a616d;
}

.month-header button {
    background: #4f46e5;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s, transform 0.2s;
}

.month-header button:hover {
    background: #4338ca;
    transform: translateY(-1px);
}

.date-selectors {
    display: flex;
    gap: 0.5rem;
}

.date-selectors select {
    flex: 1;
}

/* Groups Section Styling */
.groups-section {
    background: #1f2937;
    border-radius: 12px;
    padding: 1.5rem;
    margin-bottom: 2rem;
    border: 1px solid #374151;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.groups-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
}

.groups-header h3 {
    color: #f9fafb;
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
}

.create-group-form {
    background: #374151;
    border-radius: 8px;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
    border: 1px solid #4b5563;
}

.group-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

/* Group Members Styling */
.group-members {
    margin-top: 1rem;
    padding: 1rem;
    background: #374151;
    border-radius: 8px;
    border: 1px solid #4b5563;
}

.members-title {
    color: #d1d5db;
    margin: 0 0 1rem 0;
    font-weight: 600;
}

.members-list {
    list-style: none;
    padding: 0;
    margin: 0 0 1rem 0;
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
}

.member-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: #4b5563;
    padding: 0.5rem 0.75rem;
    border-radius: 20px;
    color: #e5e7eb;
}

.member-avatar {
    background: #6366f1;
    color: white;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 600;
}

.member-name {
    font-size: 0.875rem;
    font-weight: 500;
}

.no-members {
    color: #9ca3af;
    margin: 0 0 1rem 0;
    font-style: italic;
}

.add-member-btn {
    margin-bottom: 1rem;
}

.add-member-form {
    background: #4b5563;
    border-radius: 6px;
    padding: 1rem;
    border: 1px solid #6b7280;
}

.member-select {
    width: 100%;
    background: #374151;
    border: 1px solid #6b7280;
    border-radius: 6px;
    padding: 0.75rem;
    color: #e5e7eb;
    font-size: 1rem;
}

.member-select:focus {
    outline: none;
    border-color: #6366f1;
    box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

/* Month Header Styling */
.add-expense-btn {
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    margin-left: auto;
}

.add-expense-form {
    width: 100%;
    background: #374151;
    border-radius: 8px;
    padding: 1.5rem;
    margin-top: 1rem;
    border: 1px solid #4b5563;
}

.expense-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

/* Expense Actions Styling */
.expense-actions {
    display: flex;
    gap: 0.5rem;
    margin-top: 0.75rem;
    padding-top: 0.75rem;
    border-top: 1px solid #4b5563;
}

.expense-action-btn {
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    border-radius: 6px;
    font-weight: 500;
}

.btn-danger {
    background: #dc2626;
    color: white;
}

.btn-danger:hover {
    background: #b91c1c;
    transform: translateY(-1px);
}

/* Edit Expense Form Styling */
.edit-expense-form {
    background: #374151;
    border-radius: 8px;
    padding: 1.5rem;
    margin-top: 1rem;
    border: 1px solid #4b5563;
}

/* Fix CSS errors */
.title {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
}

.expenses-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
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
