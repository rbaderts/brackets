
<script lang="ts">
    import Bracket from '../components/Bracket.vue'
    import TournamentInfo from '../components/TournamentInfo.vue'
    import Participants from '../components/Participants.vue'
    import {axiosApiInstance} from '../main.ts';
    import { ref } from "vue";
    import { defineComponent } from 'vue';

    export default defineComponent ({
      name: "TournamentView",
      provides:  {},

      
      components: {Bracket, TournamentInfo, Participants},

      data: function() {
        return { 
          tournament: {id:0} as Tournament,
          tabs: ['Info', 'Participants', 'Bracket'],
          currentTab: 'Info',
          tab: ref('Info'),
          error: null,
          loading: true
        }

      },

      created() {
        let cmp = this
        console.log("mount entr")
        this.fetchTournament().then( function() {
              console.log("tournament: " +JSON.stringify(cmp.tournament))
        })
        console.log("mount leave")
      
      },
      watch: {
         // call again the method if the route changes
         '$route': 'fetchTournament'
      },

      computed: {
          tournamentName() {
            if (this.tournament != null) {
              return this.tournament.name
            }
          },
          playerCount() {
            if (this.tournament != null) {
              return Object.keys(this.tournament.participants).length
            }
          },
          currentTabComponent() {
            if (this.tab == "Info") {
              return 'TournamentInfo'
            }
            if (this.tab == "Participants") {
              return 'Participants'
            }
            if (this.tab == "Bracket") {
              return 'Bracket'
            }
            return 'None'
          },

      },

        methods:  {

          getTournamentId() {
            let parts = this.$route.path.split('/')
            let last = parts[parts.length-1]
            return last
          },

          fetchTournament() {

            let cmp = this;
            let id = cmp.getTournamentId();
            let url = "/tournaments/" + id
            console.log("id = " + id)
            return axiosApiInstance.get(url)
                .then(function (response) {
                    console.log("response = " + JSON.stringify(response))
                    if (cmp.getTournamentId() !== id) return
                    cmp.tournament = response.data
                    cmp.loading = false
                    cmp.$root.title = cmp.tournamentName
                    cmp.$root.players = cmp.playerCount
                })
                .catch(function (error) {
                    cmp.error = error.toString()
                    console.log(error);
                })
                .then(function () {
                });
          },

          onBeforeTransition(newVal, oldVal) {
            console.log("newValue = " + newVal);
            console.log("oldValue = " + oldVal);
          },

          openTab(evt, tabName) {
              var i, tabcontent, tablinks;
              tabcontent = document.getElementsByClassName("tabcontent");
              for (i = 0; i < tabcontent.length; i++) {
                tabcontent[i].style.display = "none";
              }
              tablinks = document.getElementsByClassName("tablinks");
              for (i = 0; i < tablinks.length; i++) {
                tablinks[i].className = tablinks[i].className.replace(" active", "");
              }
              let elem = document.getElementById(tabName)
              if (elem != null) {
                 elem.style.display = "block";
              }
              this.currentTab = tabName;
              evt.currentTarget.className += " active";

          }

        }
    })

</script>


<template>
  <div v-if="loading">
      <h1> Loading... </h1>
  </div>
  <div v-else >
    <q-page>
      <q-tabs
        v-model="tab"
        class="text-dark bg-secondary shadow-2"
        align="justify"
        indicator-color="dark"
      >
        <q-tab name="Info" icon="edit_attributes" label="Details"></q-tab>
        <q-tab name="Participants" icon="people" label="Participants"></q-tab>
        <q-tab name="Bracket" icon="account_tree" label="Bracket"></q-tab>
      </q-tabs>

      <q-separator></q-separator>
        
    <!--<q-page-sticky> -->
        <q-tab-panels keep-alive v-model="tab" @before-transition="onBeforeTransition">

          <q-tab-panel name="Info">
            <TournamentInfo :tournament="tournament" v-bind:currentTab="tab"/>
          </q-tab-panel>

          <q-tab-panel name="Participants">
            <Participants :tournamentId="tournament.id" v-bind:currentTab="tab"/>
          </q-tab-panel>

          <q-tab-panel name="Bracket">
            <Bracket :tournamentId="tournament.id" v-bind:currentTab="tab"/>
          </q-tab-panel>
        </q-tab-panels>
<!--      </q-page-sticky> -->

    </q-page>
    </div>
  </template>


<style scoped>

a {
  color: #42b983;
}
body {font-family: Arial;}

/* Style the tab */
/*
.tab {
  overflow: hidden;
  border: 1px solid #ccc;
  background-color: #f1f1f1;
}

/* Style the buttons inside the tab */
/*
.q-tab-panel {
  min-height: 300px;
  padding-left: 88x;
  padding-right:8px;
  padding-top:12px;
}

.tab button {
  background-color: inherit;
  float: left;
  border: none;
  outline: none;
  cursor: pointer;
  padding: 14px 16px;
  transition: 0.3s;
  font-size: 20px;
}

.tab button:hover {
  background-color: #ddd;
}

.tab button.active {
  background-color: #ccc;
}

.tabcontent {
  display: none;
  padding: 2px 2px;
  border: 1px solid #ccc;
  border-top: none;
}
*/
</style>
