<script setup lang="ts">
import InputGroup from 'primevue/inputgroup'
import InputGroupAddon from 'primevue/inputgroupaddon'
import InputText from 'primevue/inputtext'
import FloatLabel from 'primevue/floatlabel';
import Button from 'primevue/button';
import { ref } from 'vue';
import axios from 'axios'

const username = ref<string>('')
const password = ref<string>('')
const handleLogin = async () => {
  const response = await axios.post('http://localhost:8080/v1/login_user',
      {
        username: username.value,
        password: password.value
      },)
  }
</script>

<template>  
  <div class="flex flex-column row-gap-5">
    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-user"></i>
      </InputGroupAddon>
      <FloatLabel>
        <InputText id="username" v-model="username" />
        <label for="username">Username</label>
      </FloatLabel>
    </InputGroup>
    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-lock"></i>
      </InputGroupAddon>
      <FloatLabel>
        <InputText id="password" type="password" v-model="password" />
        <label for="password">Password</label>
      </FloatLabel>
    </InputGroup>
    <Button label="Login" @click="handleLogin"/>
  </div>
</template>



