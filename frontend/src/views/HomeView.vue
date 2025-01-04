<script setup lang="ts">
import Toast from 'primevue/toast';

import LoginUser from '@/components/LoginUser.vue';
import UserInfo from '@/components/UserInfo.vue'
import store from '../store'
</script>

<template>
  <main>
    <Toast />
    <h1 class="green">Welcome to Simple Bank!</h1>
    <UserInfo v-if="store.state.user" :user="store.state.user" />
    <LoginUser v-else />
  </main>
</template>
