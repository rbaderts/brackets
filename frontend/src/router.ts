/*
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import generatedRoutes from '~pages'

import mainLayout from 'src/layouts/mainLayout.vue'
import blankLayout from 'src/layouts/blankLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: mainLayout,
    children: generatedRoutes,
  },
  {
    path: '/',
    component: blankLayout,
    children: [
      {
        name: 'all',
        path: ':all(.*)*',
        component: () => import('src/pages/blank/404.vue'),
      },
      {
        name: 'login',
        path: 'login',
        component: () => import('src/pages/blank/login.vue'),
      },
      {
        name: 'register',
        path: 'register',
        component: () => import('src/pages/blank/register.vue'),
      },
    ],
  },
]
const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  console.log('router beforeEach', router.currentRoute.value.fullPath)
  next()
})
export default router
*/


import {createRouter, createWebHistory} from 'vue-router';

import Home from './views/Home.vue'
import TournamentView from './views/TournamentView.vue'
//import TournamentList from "./components/TournamentList.vue"

const router = createRouter({
     history: createWebHistory(),
    routes: [
  {
    path: '/',
    name: 'Home',
    component: Home,

    /*
    children: [
      {
        // UserPosts will be rendered inside User's <router-view>
        // when /user/:id/posts is matched
        path: 'tournament/:id',
        component: TournamentView
      }
    ],
      */
    meta: {
      isAuthenticated: false
    }
  },{
    path: '/tournament/:id',
    name: 'TournamentView',
    component: TournamentView,
  },

    ]
})
export default router


  /*
  {

    path: '/tournament/:id',
    name: 'TournamentView',
    components: {
      default: TouranmentView
    },
//    component: () => import('./views/TournamentView.vue'),
    meta: {
      isAuthenticated: true
    }
  },
  */
 /*
  {
    path: '/secured',
    name: 'Secured',
    meta: {
      isAuthenticated: true
    },
    component: () => import('./views/Secured.vue')
  },
  {
    path: '/unauthorized',
    name: 'Unauthorized',
    meta: {
      isAuthenticated: false
    },
    component: () => import('./views/Unauthorized.vue')
  }
  */
