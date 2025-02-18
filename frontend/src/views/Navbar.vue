<template>
  <nav class="bg-blue-600 p-4 text-white shadow-md">
    <div class="container mx-auto flex justify-between items-center">
      <router-link to="/" class="text-lg font-semibold">MyApp</router-link>
      
      <div class="space-x-4">
        <template v-if="!isAuthenticated">
          <router-link to="/login" class="hover:underline">Login</router-link>
          <router-link to="/register" class="hover:underline">Register</router-link>
        </template>
        <template v-else>
          <span class="mr-2">Hello, {{ userInfo.username }}</span>
          <router-link to="/profile" class="hover:underline">Profile</router-link>
          <button
            @click="logout"
            class="bg-red-500 px-3 py-1 rounded hover:bg-red-600 transition duration-300"
          >
            Logout
          </button>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const router = useRouter();
const isAuthenticated = ref(false);
const userInfo = ref({ username: '' });

// Watch route changes to update authentication state
watch(
  () => router.currentRoute.value,
  () => {
    const storedUser = localStorage.getItem('userInfo');
    if (storedUser) {
      userInfo.value = JSON.parse(storedUser);
      isAuthenticated.value = true;
    } else {
      isAuthenticated.value = false;
      userInfo.value = { username: '' };
    }
  },
  { immediate: true } // Trigger on component mount
);

const logout = async () => {
  try {
    const response = await axios.get('http://localhost:8080/logout', { withCredentials: true });

    if (response.status === 200) {
      // Clear user data and redirect to login
      localStorage.removeItem('userInfo');
      isAuthenticated.value = false;
      userInfo.value = { username: '' };
      router.push('/login');
    } else {
      throw new Error('Unexpected response from server');
    }
  } catch (error) {
    console.error('Logout error:', error);
   
  }
};
</script>