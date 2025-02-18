import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/Login.vue';
import Profile from '../views/Profile.vue';
import Register from '../views/Register.vue';

const routes = [
  { path: '/login', component: Login, meta: { requiresGuest: true } }, 
  { path: '/register', component: Register, meta: { requiresGuest: true } }, 
  { path: '/profile', component: Profile, meta: { requiresAuth: true } }, // Hanya untuk pengguna yang sudah login
  { path: '/', redirect: '/login' },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});



// Navigation guard
router.beforeEach(async (to, from, next) => {
  // Cek status login
  try {

    const authStatus = localStorage.getItem('userInfo');


    if (to.meta.requiresAuth && !authStatus) {
      // Jika rute memerlukan autentikasi dan user belum login, arahkan ke login
      next('/login');
    } else if (to.meta.requiresGuest && authStatus) {
      // Jika rute hanya untuk tamu dan user sudah login, arahkan ke profile
      next('/profile');
    } else {
      // Izinkan akses ke rute
      next();
    }
  } catch (error) {
    // Jika terjadi error (misalnya, token tidak valid), arahkan ke login
    if (to.meta.requiresAuth) {
      next('/login');
    } else {
      next();
    }
  }
});

export default router;