<script lang="ts">
import { defineComponent, PropType, reactive } from "vue";
import {axiosApiInstance} from '../main.ts';

export default defineComponent({
  name: "TournamentInfo",
  inject: ["router"],
  props: {
    tournament: {
      type: Object as PropType<Tournament>
      /*
      default: function () {
        return {
          tournamentId: 0,
          name: "Unknown Tournament",
        };
      },
    },*/
    },
    currentTab: String,
  },
  watch: {
    currentTab(val) {
      console.log("val = " + val);
      if (val == "Info") {
        console.log("Info activated");
      }
    },
  },

  data: function () {
    return {
      isLoading: true,
      error: null,
      tournament: {id:0} as Tournament
    };
  },

  updated() {
    //this.updateList()
  },

  mounted: function () {
    //this.fetchTournament().then(function() {
    //})
  },

  methods: {
    /*
             fetchTournament() {

                let cmp = this;
                let url = "/tournaments/" + cmp.tournamentId
                return axiosApiInstance.get(url)
                    .then(function (response) {
                        console.log("response = " + JSON.stringify(response))
                        cmp.tournament = response.data
                    })
                    .catch(function (error) {
                        cmp.error = error.toString()
                        console.log(error);
                    })
                    .then(function () {
                    });
            } 
            */
    /*
            updateList() {
                const table = document.querySelector("#tournament_list");
                for (const currentRow of table.rows) {
                     currentRow.onclick = this.createClickHandler(currentRow);
                }
            }
            */
  },
});
</script>

<template>
  <div class="container">
    <div class="q-pa-xs">
      <div class="q-gutter-y-md column" style="max-width: 300px">
        <q-input :dense="true" v-model="tournament.name" label="Name"></q-input>
        <q-input
          :dense="true"
          v-model="tournament.status"
          label="Status"
        ></q-input>
        <q-field :dense="true" label="Tournament Type" stack-label>
          <template v-slot:control>
            <div class="self-center full-width no-outline" tabindex="0">
              {{ tournament.typ == 1 ? "Single Elim" : "Double Elim" }}
            </div>
          </template>
        </q-field>
      </div>
    </div>
  </div>
</template>


<style lang="scss" scoped></style>
