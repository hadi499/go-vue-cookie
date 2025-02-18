<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';

const username = ref('');
const router = useRouter();

onMounted(async () => {
  try {
    const response = await axios.get('http://localhost:8080/protected/profile', { withCredentials: true });
    if (response.status === 200 && response.data.username) {
      username.value = response.data.username;
    } else {
      router.push('/login');
    }
  } catch (error) {
    router.push('/login');
  }
});


</script>

<template>
  <div class="w-80 mt-6 mx-auto">
    <h1>Profile</h1>
    <p>Welcome, <span class="text-lg font-semibold">{{ username }}</span></p>
    
  </div>
</template>