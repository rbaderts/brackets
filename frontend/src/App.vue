<!--
<template>
  <router-view />
</template>
-->
<script lang="ts">
import { defineComponent } from "vue";
import Home from "./views/Home.vue";
import TournamentView from "./views/TournamentView.vue";
export default defineComponent({
  name: "App",
  inject: ["keycloak"],
  data() {
     return {
         title: "",
         players: 0,
     }
  },
  components: {
    Home,
    TournamentView,
  },
  created() {
    // Should work referencing custom property on 'this' in typescript
    console.log("App created", this.$keycloak);
  },
});
</script>

<template>
  <q-layout view="hHh lpR fFf">
    <q-header elevated class="bg-primary text-white" height-hint="40">
      <q-toolbar>
        <q-toolbar-title>
          <q-avatar>
            <img src="/BracketsLogo.png" />
          </q-avatar>
          Brackets
        </q-toolbar-title>
      </q-toolbar>
      <div class="q-py-none q-my-0 q-px-sm row">
        <div class="col-4">
            {{title}}
        </div>
        <div class="col-4">
            Players: {{players}}
         </div>
      </div>

<!--
      <q-tabs align="left">
        <q-route-tab to="/page1" label="Page One" />
        <q-route-tab to="/page2" label="Page Two" />
        <q-route-tab to="/page3" label="Page Three" />
      </q-tabs>
      -->
    </q-header>

   <q-page-container style="padding-top: 10px; padding-bottom: 15px">
      <router-view />
   </q-page-container>

  </q-layout>
</template>

<style>
/*
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: left;
  color: #2c3e50;
  */
#app {
  margin-top: 80px;
}
</style>
